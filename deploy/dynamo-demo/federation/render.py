# SPDX-FileCopyrightText: Copyright (c) 2026 The k8gb Contributors.
# SPDX-License-Identifier: Apache-2.0

"""Render an opt-in regional DGD, Relay and HTTP federation entry."""

import argparse
import copy
from pathlib import Path

import yaml
from router import Config

CPU_IMAGE = "k8gb-dynamo-cpu:1.7.0.dev20261026"
ROUTER_IMAGE = "k8gb-dynamo-global-router:prototype"
REVISION = "c1899de289a04d12100db370d81485cdf75e47ca"


def render(config, namespace, replicas=1, with_gslb=False):
    baseline = list(yaml.safe_load_all((Path(__file__).parent.parent / "regional.yaml").read_text()))
    dgd = copy.deepcopy(baseline[1])
    dgd["metadata"]["namespace"] = namespace
    for component in dgd["spec"]["components"]:
        component["globalDynamoNamespace"] = True
        component["runtimeVersionOverride"] = "1.7.0"
        if component["type"] == "worker":
            component["replicas"] = replicas
        pod = component["podTemplate"]["spec"]
        for container in pod["containers"] + pod["initContainers"]:
            container["image"] = CPU_IMAGE
        if component["type"] == "frontend":
            pod["volumes"].append({"name": "mdc-cache", "emptyDir": {}})
            pod["containers"][0]["volumeMounts"].append({"name": "mdc-cache", "mountPath": "/.cache"})

    def deployment(name, container, **pod_fields):
        labels = {"app": name}
        return {
            "apiVersion": "apps/v1",
            "kind": "Deployment",
            "metadata": {"name": name, "namespace": namespace},
            "spec": {
                "replicas": 1,
                "selector": {"matchLabels": labels},
                "template": {"metadata": {"labels": labels}, "spec": {"containers": [container], **pod_fields}},
            },
        }

    def service(name, port, node_port=None, selector=None):
        spec = {
            "selector": selector or {"app": name},
            "ports": [{"name": "http" if port != 50051 else "grpc", "port": port, "targetPort": port}],
        }
        if node_port:
            spec["type"] = "NodePort"
            spec["ports"][0]["nodePort"] = node_port
        return {"apiVersion": "v1", "kind": "Service", "metadata": {"name": name, "namespace": namespace}, "spec": spec}

    env = [
        {"name": key, "value": value}
        for key, value in {
            "DYN_DISCOVERY_BACKEND": "kubernetes",
            "DYN_REQUEST_PLANE": "tcp",
            "DYN_EVENT_PLANE": "zmq",
            "DYN_NAMESPACE": "federation-relay-" + config.local_region,
            "DYN_RUNTIME_NUM_WORKER_THREADS": "2",
            "DYN_RUNTIME_MAX_BLOCKING_THREADS": "8",
            "DYN_COMPUTE_THREADS": "2",
            "DYN_RELAY_POOL_HEARTBEAT_INTERVAL_MS": "1000",
            "DYN_RELAY_READINESS_HEARTBEAT_INTERVAL_MS": "1000",
        }.items()
    ]
    env.extend(
        {"name": key, "valueFrom": {"fieldRef": {"fieldPath": field}}}
        for key, field in {
            "POD_NAME": "metadata.name",
            "POD_NAMESPACE": "metadata.namespace",
            "POD_UID": "metadata.uid",
        }.items()
    )
    relay = deployment(
        "dynamo-relay",
        {
            "name": "main",
            "image": CPU_IMAGE,
            "command": ["python3", "-m", "dynamo.kv_dc_relay"],
            "args": [
                "--dc-id",
                config.local_region,
                "--watch-all",
                "--expected-unique-blocks",
                "4096",
                "--bind",
                "0.0.0.0:50051",
            ],
            "env": env,
            "ports": [{"name": "grpc", "containerPort": 50051}],
            "readinessProbe": {"tcpSocket": {"port": "grpc"}, "periodSeconds": 2},
            "resources": {"requests": {"cpu": "50m", "memory": "128Mi"}, "limits": {"cpu": "1", "memory": "512Mi"}},
        },
        serviceAccountName="dynamo-demo-k8s-service-discovery",
    )
    security = {"runAsNonRoot": True, "runAsUser": 1000, "runAsGroup": 1000, "fsGroup": 1000}
    model_mount = {"name": "model", "mountPath": "/opt/model"}
    router = deployment(
        "dynamo-global-router",
        {
            "name": "main",
            "image": ROUTER_IMAGE,
            "args": ["--config", "/config/router.yaml"],
            "ports": [{"name": "http", "containerPort": 8080}],
            "env": [{"name": "HF_HUB_OFFLINE", "value": "1"}],
            "volumeMounts": [
                {**model_mount, "readOnly": True},
                {"name": "config", "mountPath": "/config", "readOnly": True},
            ],
            "readinessProbe": {"httpGet": {"path": "/ready", "port": "http"}, "periodSeconds": 2},
            "livenessProbe": {"httpGet": {"path": "/live", "port": "http"}, "periodSeconds": 10},
            "resources": {"requests": {"cpu": "100m", "memory": "128Mi"}, "limits": {"cpu": "1", "memory": "512Mi"}},
        },
        securityContext=security,
        initContainers=[
            {
                "name": "tokenizer",
                "image": ROUTER_IMAGE,
                "command": ["hf", "download"],
                "args": [
                    config.model,
                    "config.json",
                    "tokenizer.json",
                    "tokenizer_config.json",
                    "--revision",
                    REVISION,
                    "--local-dir",
                    "/opt/model",
                ],
                "volumeMounts": [model_mount],
                "resources": {"requests": {"cpu": "50m", "memory": "128Mi"}, "limits": {"cpu": "1", "memory": "512Mi"}},
            }
        ],
        volumes=[{"name": "model", "emptyDir": {}}, {"name": "config", "configMap": {"name": "dynamo-global-router"}}],
    )
    config_data = {**config.__dict__, "regions": list(config.regions)}
    documents = [
        {"apiVersion": "v1", "kind": "Namespace", "metadata": {"name": namespace}},
        dgd,
        {
            "apiVersion": "v1",
            "kind": "ConfigMap",
            "metadata": {"name": "dynamo-global-router", "namespace": namespace},
            "data": {"router.yaml": yaml.safe_dump(config_data, sort_keys=False)},
        },
        relay,
        service("dynamo-relay", 50051, 30551 if with_gslb else None),
        router,
        service("dynamo-global-router", 8080),
    ]
    if with_gslb:
        documents.append(
            service("dynamo-regional-origin", 8000, 30800, {"nvidia.com/selector": "dynamo-demo-frontend"})
        )
        for item in copy.deepcopy(baseline[2:]):
            item["metadata"]["namespace"] = namespace
            if item["kind"] == "Ingress":
                backend = item["spec"]["rules"][0]["http"]["paths"][0]["backend"]
                backend["service"] = {"name": "dynamo-global-router", "port": {"name": "http"}}
            documents.append(item)
    return documents


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--config", required=True)
    parser.add_argument("--namespace", default="dynamo-demo")
    parser.add_argument("--worker-replicas", type=int, choices=(1, 2), default=1)
    parser.add_argument("--with-gslb", action="store_true")
    args = parser.parse_args()
    print(
        yaml.safe_dump_all(
            render(Config.load(args.config), args.namespace, args.worker_replicas, args.with_gslb), sort_keys=False
        )
    )
