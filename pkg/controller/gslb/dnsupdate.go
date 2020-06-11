package gslb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	ibclient "github.com/AbsaOSS/infoblox-go-client"
	k8gbv1beta1 "github.com/AbsaOSS/k8gb/pkg/apis/k8gb/v1beta1"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	"github.com/miekg/dns"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileGslb) getGslbIngressIPs(gslb *k8gbv1beta1.Gslb) ([]string, error) {
	nn := types.NamespacedName{
		Name:      gslb.Name,
		Namespace: gslb.Namespace,
	}

	gslbIngress := &v1beta1.Ingress{}

	err := r.client.Get(context.TODO(), nn, gslbIngress)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info(fmt.Sprintf("Can't find gslb Ingress: %s", gslb.Name))
		}
		return nil, err
	}

	var gslbIngressIPs []string

	for _, ip := range gslbIngress.Status.LoadBalancer.Ingress {
		gslbIngressIPs = append(gslbIngressIPs, ip.IP)
	}

	return gslbIngressIPs, nil
}

func getExternalClusterFQDNs(gslb *k8gbv1beta1.Gslb) []string {
	extGslbClustersGeoTagsVar := os.Getenv("EXT_GSLB_CLUSTERS_GEO_TAGS")

	if extGslbClustersGeoTagsVar == "" {
		log.Info("No other Gslb enabled clusters are defined in the configuration...Working standalone")
		return nil
	}

	extGslbClustersGeoTags := strings.Split(extGslbClustersGeoTagsVar, ",")

	var extGslbClusters []string
	for _, geoTag := range extGslbClustersGeoTags {
		cluster := nsServerName(gslb, geoTag)
		extGslbClusters = append(extGslbClusters, cluster)
	}
	return extGslbClusters
}

func getExternalClusterHeartbeatFQDNs(gslb *k8gbv1beta1.Gslb) []string {
	extGslbClustersGeoTagsVar := os.Getenv("EXT_GSLB_CLUSTERS_GEO_TAGS")

	if extGslbClustersGeoTagsVar == "" {
		log.Info("No other Gslb enabled clusters are defined in the configuration...Working standalone")
		return nil
	}

	extGslbClustersGeoTags := strings.Split(extGslbClustersGeoTagsVar, ",")

	var extGslbClusters []string
	for _, geoTag := range extGslbClustersGeoTags {
		cluster := heartbeatFQDN(gslb, geoTag)
		extGslbClusters = append(extGslbClusters, cluster)
	}
	return extGslbClusters
}

func getExternalTargets(gslb *k8gbv1beta1.Gslb, host string) ([]string, error) {

	extGslbClusters := getExternalClusterFQDNs(gslb)

	var targets []string

	for _, cluster := range extGslbClusters {
		log.Info(fmt.Sprintf("Adding external Gslb targets from %s cluster...", cluster))
		g := new(dns.Msg)
		host = fmt.Sprintf("localtargets.%s.", host) //Convert to true FQDN with dot at the end. Otherwise dns lib freaks out
		g.SetQuestion(host, dns.TypeA)

		localTestDNSinject := os.Getenv("OVERRIDE_WITH_FAKE_EXT_DNS")

		var ns string

		if localTestDNSinject == "true" {
			ns = "127.0.0.1:7753"
		} else {
			ns = fmt.Sprintf("%s:53", cluster)
		}

		a, err := dns.Exchange(g, ns)
		if err != nil {
			log.Info(fmt.Sprintf("Error contacting external Gslb cluster(%s) : (%v)", cluster, err))
			return nil, nil
		}
		var clusterTargets []string

		for _, A := range a.Answer {
			IP := strings.Split(A.String(), "\t")[4]
			clusterTargets = append(clusterTargets, IP)
		}
		if len(clusterTargets) > 0 {
			targets = append(targets, clusterTargets...)
			log.Info(fmt.Sprintf("Added external %s Gslb targets from %s cluster", clusterTargets, cluster))
		}
	}

	return targets, nil
}

func (r *ReconcileGslb) gslbDNSEndpoint(gslb *k8gbv1beta1.Gslb) (*externaldns.DNSEndpoint, error) {
	var gslbHosts []*externaldns.Endpoint
	var ttl = externaldns.TTL(gslb.Spec.Strategy.DNSTtlSeconds)

	serviceHealth, err := r.getServiceHealthStatus(gslb)
	if err != nil {
		return nil, err
	}

	localTargets, err := r.getGslbIngressIPs(gslb)
	if err != nil {
		return nil, err
	}

	for host, health := range serviceHealth {
		var finalTargets []string

		if health == "Healthy" {
			finalTargets = append(finalTargets, localTargets...)
			localTargetsHost := fmt.Sprintf("localtargets.%s", host)
			dnsRecord := &externaldns.Endpoint{
				DNSName:    localTargetsHost,
				RecordTTL:  ttl,
				RecordType: "A",
				Targets:    localTargets,
			}
			gslbHosts = append(gslbHosts, dnsRecord)
		}

		// Check if host is alive on external Gslb
		externalTargets, err := getExternalTargets(gslb, host)
		if err != nil {
			return nil, err
		}
		if len(externalTargets) > 0 {
			switch gslb.Spec.Strategy.Type {
			case "roundRobin":
				finalTargets = append(finalTargets, externalTargets...)
			case "failover":
				clusterGeoTag := os.Getenv("CLUSTER_GEO_TAG")
				// If cluster is Primary
				if gslb.Spec.Strategy.PrimaryGeoTag == clusterGeoTag {
					// If cluster is Primary and Healthy return only own targets
					// If cluster is Primary and Unhealthy return Secondary external targets
					if health != "Healthy" {
						finalTargets = externalTargets
					}
				} else {
					// If cluster is Secondary and Primary external cluster is Healthy
					// then return Primary external targets.
					// Return own targets by default.
					if len(externalTargets) > 0 {
						finalTargets = externalTargets
					}
				}
			}
		}

		if len(finalTargets) > 0 {
			dnsRecord := &externaldns.Endpoint{
				DNSName:    host,
				RecordTTL:  ttl,
				RecordType: "A",
				Targets:    finalTargets,
			}
			gslbHosts = append(gslbHosts, dnsRecord)
		}
	}
	dnsEndpointSpec := externaldns.DNSEndpointSpec{
		Endpoints: gslbHosts,
	}

	dnsEndpoint := &externaldns.DNSEndpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gslb.Name,
			Namespace:   gslb.Namespace,
			Annotations: map[string]string{"k8gb.absa.oss/dnstype": "local"},
		},
		Spec: dnsEndpointSpec,
	}

	err = controllerutil.SetControllerReference(gslb, dnsEndpoint, r.scheme)
	if err != nil {
		return nil, err
	}
	return dnsEndpoint, err
}

func nsServerName(gslb *k8gbv1beta1.Gslb, clusterGeoTag string) string {
	edgeDNSZone := os.Getenv("EDGE_DNS_ZONE")
	if len(clusterGeoTag) == 0 {
		clusterGeoTag = "default"
	}
	return fmt.Sprintf("%s-ns-%s.%s", gslb.Name, clusterGeoTag, edgeDNSZone)
}

func heartbeatFQDN(gslb *k8gbv1beta1.Gslb, clusterGeoTag string) string {
	edgeDNSZone := os.Getenv("EDGE_DNS_ZONE")
	if len(clusterGeoTag) == 0 {
		clusterGeoTag = "default"
	}
	return fmt.Sprintf("%s-heartbeat-%s.%s", gslb.Name, clusterGeoTag, edgeDNSZone)
}

type fakeInfobloxConnector struct {
	//createObjectObj interface{}

	getObjectObj interface{}
	getObjectRef string

	//deleteObjectRef string

	//updateObjectObj interface{}
	//updateObjectRef string

	resultObject interface{}

	fakeRefReturn string
}

func (c *fakeInfobloxConnector) CreateObject(obj ibclient.IBObject) (string, error) {

	return c.fakeRefReturn, nil
}

func (c *fakeInfobloxConnector) GetObject(obj ibclient.IBObject, ref string, res interface{}) (err error) {
	return nil
}

func (c *fakeInfobloxConnector) DeleteObject(ref string) (string, error) {

	return c.fakeRefReturn, nil
}

func (c *fakeInfobloxConnector) UpdateObject(obj ibclient.IBObject, ref string) (string, error) {

	return c.fakeRefReturn, nil
}

func infobloxConnection() (*ibclient.ObjectManager, error) {
	hostConfig := ibclient.HostConfig{
		Host:     os.Getenv("INFOBLOX_GRID_HOST"),
		Version:  os.Getenv("INFOBLOX_WAPI_VERSION"),
		Port:     os.Getenv("INFOBLOX_WAPI_PORT"),
		Username: os.Getenv("EXTERNAL_DNS_INFOBLOX_WAPI_USERNAME"),
		Password: os.Getenv("EXTERNAL_DNS_INFOBLOX_WAPI_PASSWORD"),
	}
	transportConfig := ibclient.NewTransportConfig("false", 20, 10)
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}

	fakeInfoblox := os.Getenv("FAKE_INFOBLOX")

	var objMgr *ibclient.ObjectManager

	if len(fakeInfoblox) > 0 {
		fqdn := "fakezone.example.com"
		fakeRefReturn := "zone_delegated/ZG5zLnpvbmUkLl9kZWZhdWx0LnphLmNvLmFic2EuY2Fhcy5vaG15Z2xiLmdzbGJpYmNsaWVudA:fakezone.example.com/default"
		ohmyFakeConnector := &fakeInfobloxConnector{
			getObjectObj: ibclient.NewZoneDelegated(ibclient.ZoneDelegated{Fqdn: fqdn}),
			getObjectRef: "",
			resultObject: []ibclient.ZoneDelegated{*ibclient.NewZoneDelegated(ibclient.ZoneDelegated{Fqdn: fqdn, Ref: fakeRefReturn})},
		}
		objMgr = ibclient.NewObjectManager(ohmyFakeConnector, "ohmyclient", "")
	} else {
		conn, err := ibclient.NewConnector(hostConfig, transportConfig, requestBuilder, requestor)
		if err != nil {
			return nil, err
		}
		defer func() {
			err = conn.Logout()
			if err != nil {
				log.Error(err, "Failed to close connection to infoblox")
			}
		}()
		objMgr = ibclient.NewObjectManager(conn, "ohmyclient", "")
	}
	return objMgr, nil
}

func checkAliveFromTXT(dnsserver string, fqdn string, splitBrainThreshold time.Duration) error {
	localTestDNSinject := os.Getenv("OVERRIDE_WITH_FAKE_EXT_DNS")

	var ns string

	if localTestDNSinject == "true" {
		ns = "127.0.0.1:7753"
	} else {
		ns = fmt.Sprintf("%s:53", dnsserver)
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeTXT)
	txt, err := dns.Exchange(m, ns)
	if err != nil {
		log.Info(fmt.Sprintf("Error contacting EdgeDNS server (%s) for TXT split brain record: (%s)", ns, err))
		return err
	}
	var timestamp string
	if len(txt.Answer) > 0 {
		if t, ok := txt.Answer[0].(*dns.TXT); ok {
			log.Info(fmt.Sprintf("Split brain TXT raw record: %s", t.String()))
			timestamp = strings.Split(t.String(), "\t")[4]
			timestamp = strings.Trim(timestamp, "\"") // Otherwise time.Parse() will miserably fail
		}
	}

	if len(timestamp) > 0 {
		log.Info(fmt.Sprintf("Split brain TXT raw time stamp: %s", timestamp))
		timeFromTXT, err := time.Parse("2006-01-02T15:04:05", timestamp)
		if err != nil {
			return err
		}

		log.Info(fmt.Sprintf("Split brain TXT parsed time stamp: %s", timeFromTXT))
		now := time.Now().UTC()

		diff := now.Sub(timeFromTXT)
		log.Info(fmt.Sprintf("Split brain TXT time diff: %s", diff))

		if diff > splitBrainThreshold {
			return errors.NewGone(fmt.Sprintf("Split brain TXT record expired the time threshold: (%s)", splitBrainThreshold))
		}

		return nil
	}
	return errors.NewGone(fmt.Sprintf("Can't find split brain TXT record at EdgeDNS server(%s) and record %s ", ns, fqdn))

}

func filterOutDelegateTo(delegateTo []ibclient.NameServer, fqdn string) []ibclient.NameServer {
	for i := 0; i < len(delegateTo); i++ {
		if delegateTo[i].Name == fqdn {
			delegateTo = append(delegateTo[:i], delegateTo[i+1:]...)
			i--
		}
	}
	return delegateTo
}

func (r *ReconcileGslb) configureZoneDelegation(gslb *k8gbv1beta1.Gslb) (*reconcile.Result, error) {
	clusterGeoTag := os.Getenv("CLUSTER_GEO_TAG")
	infobloxGridHost := os.Getenv("INFOBLOX_GRID_HOST")
	if len(infobloxGridHost) > 0 {

		objMgr, err := infobloxConnection()
		if err != nil {
			return &reconcile.Result{}, err
		}
		addresses, err := r.getGslbIngressIPs(gslb)
		if err != nil {
			return &reconcile.Result{}, err
		}
		delegateTo := []ibclient.NameServer{}

		for _, address := range addresses {
			nameServer := ibclient.NameServer{Address: address, Name: nsServerName(gslb, clusterGeoTag)}
			delegateTo = append(delegateTo, nameServer)
		}

		gslbZoneName := os.Getenv("DNS_ZONE")
		edgeDNSZone := os.Getenv("EDGE_DNS_ZONE")
		edgeDNSServer := os.Getenv("EDGE_DNS_SERVER")
		findZone, err := objMgr.GetZoneDelegated(gslbZoneName)
		if err != nil {
			return &reconcile.Result{}, err
		}

		if findZone != nil {
			err = checkZoneDelegated(findZone, gslbZoneName)
			if err != nil {
				return &reconcile.Result{}, err
			}
			if len(findZone.Ref) > 0 {

				// Drop own records for straight away update
				existingDelegateTo := filterOutDelegateTo(findZone.DelegateTo, nsServerName(gslb, clusterGeoTag))
				existingDelegateTo = append(existingDelegateTo, delegateTo...)

				// Drop external records if they are stale
				extClusters := getExternalClusterHeartbeatFQDNs(gslb)
				for _, extCluster := range extClusters {
					err = checkAliveFromTXT(edgeDNSServer, extCluster, time.Second * time.Duration(gslb.Spec.Strategy.SplitBrainThresholdSeconds))
					if err != nil {
						log.Error(err, "got the error from TXT based checkAlive")
						log.Info(fmt.Sprintf("External cluster (%s) doesn't look alive, filtering it out from delegated zone configuration...", extCluster))
						existingDelegateTo = filterOutDelegateTo(existingDelegateTo, extCluster)
					}
				}
				log.Info(fmt.Sprintf("Updating delegated zone(%s) with the server list(%v)", gslbZoneName, existingDelegateTo))

				_, err = objMgr.UpdateZoneDelegated(findZone.Ref, existingDelegateTo)
				if err != nil {
					return &reconcile.Result{}, err
				}
			}
		} else {
			log.Info(fmt.Sprintf("Creating delegated zone(%s)...", gslbZoneName))
			_, err = objMgr.CreateZoneDelegated(gslbZoneName, delegateTo)
			if err != nil {
				return &reconcile.Result{}, err
			}
		}

		edgeTimestamp := fmt.Sprint(time.Now().UTC().Format("2006-01-02T15:04:05"))
		heartbeatTXTName := fmt.Sprintf("%s-heartbeat-%s.%s", gslb.Name, clusterGeoTag, edgeDNSZone)
		heartbeatTXTRecord, err := objMgr.GetTXTRecord(heartbeatTXTName)
		if err != nil {
			return &reconcile.Result{}, err
		}
		if heartbeatTXTRecord == nil {
			log.Info(fmt.Sprintf("Creating split brain TXT record(%s)...", heartbeatTXTName))
			_, err := objMgr.CreateTXTRecord(heartbeatTXTName, edgeTimestamp, gslb.Spec.Strategy.DNSTtlSeconds, "default")
			if err != nil {
				return &reconcile.Result{}, err
			}
		} else {
			log.Info(fmt.Sprintf("Updating split brain TXT record(%s)...", heartbeatTXTName))
			_, err := objMgr.UpdateTXTRecord(heartbeatTXTName, edgeTimestamp)
			if err != nil {
				return &reconcile.Result{}, err
			}
		}
	}
	return nil, nil
}

func (r *ReconcileGslb) ensureDNSEndpoint(request reconcile.Request,
	gslb *k8gbv1beta1.Gslb,
	i *externaldns.DNSEndpoint,
) (*reconcile.Result, error) {
	found := &externaldns.DNSEndpoint{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      i.Name,
		Namespace: gslb.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the DNSEndpoint
		log.Info(fmt.Sprintf("Creating a new DNSEndpoint:\n %s", prettyPrint(i)))
		err = r.client.Create(context.TODO(), i)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new DNSEndpoint", "DNSEndpoint.Namespace", i.Namespace, "DNSEndpoint.Name", i.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get DNSEndpoint")
		return &reconcile.Result{}, err
	}

	// Update existing object with new spec
	found.Spec = i.Spec
	err = r.client.Update(context.TODO(), found)

	if err != nil {
		// Update failed
		log.Error(err, "Failed to update DNSEndpoint", "DNSEndpoint.Namespace", found.Namespace, "DNSEndpoint.Name", found.Name)
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func checkZoneDelegated(findZone *ibclient.ZoneDelegated, gslbZoneName string) error {
	if findZone.Fqdn != gslbZoneName {
		err := fmt.Errorf("delegated zone returned from infoblox(%s) does not match requested gslb zone(%s)", findZone.Fqdn, gslbZoneName)
		return err
	}
	return nil
}

func prettyPrint(s interface{}) string {
	prettyStruct, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		fmt.Println("can't convert struct to json")
	}
	return string(prettyStruct)
}
