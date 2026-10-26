# Copyright 2026 The k8gb Contributors.
# Licensed under the Apache License, Version 2.0.

# Keep the upstream CPU engine; ai-dynamo[vllm] would install CUDA dependencies.
ARG TARGETARCH
FROM vllm/vllm-openai-cpu:v0.30.0-x86_64@sha256:0fd700a207f4f4f0cc9c87fb425b5fa2aafa420c9affca904e2fd44920cbf5a7 AS cpu-amd64
FROM vllm/vllm-openai-cpu:v0.30.0-arm64@sha256:09e656588c2dc572240cb62a22f75681dc640a99c61537bbd55e9eb1dc3d4fbc AS cpu-arm64
FROM cpu-${TARGETARCH}

RUN python3 -m pip install --no-cache-dir \
      ai-dynamo==1.6.0 uvloop==0.22.1 blake3==1.0.9

ENV HF_HOME=/tmp/huggingface \
    VLLM_NO_USAGE_STATS=1 \
    DO_NOT_TRACK=1

ENTRYPOINT []
CMD ["python3", "-m", "dynamo.vllm", "--help"]
