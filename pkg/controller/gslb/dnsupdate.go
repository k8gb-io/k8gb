package gslb

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	ibclient "github.com/AbsaOSS/infoblox-go-client"
	ohmyglbv1beta1 "github.com/AbsaOSS/ohmyglb/pkg/apis/ohmyglb/v1beta1"
	externaldns "github.com/kubernetes-incubator/external-dns/endpoint"
	"github.com/miekg/dns"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ReconcileGslb) getGslbIngressIPs(gslb *ohmyglbv1beta1.Gslb) ([]string, error) {
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

func getExternalTargets(host string) ([]string, error) {

	extGslbClustersVar := os.Getenv("EXT_GSLB_CLUSTERS")

	if extGslbClustersVar == "" {
		log.Info("No other Gslb enabled clusters are defined in the configuration...Working standalone")
		return nil, nil
	}

	extGslbClusters := strings.Split(extGslbClustersVar, ",")

	var targets []string

	for _, cluster := range extGslbClusters {
		log.Info(fmt.Sprintf("Adding external Gslb targets from %s cluster...", cluster))
		g := new(dns.Msg)
		host = fmt.Sprintf("localtargets.%s.", host) //Convert to true FQDN with dot at the endOtherwise dns lib freaks out
		g.SetQuestion(host, dns.TypeA)
		ns := fmt.Sprintf("%s:53", cluster)
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

func (r *ReconcileGslb) gslbDNSEndpoint(gslb *ohmyglbv1beta1.Gslb) (*externaldns.DNSEndpoint, error) {
	var gslbHosts []*externaldns.Endpoint

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
				RecordTTL:  30,
				RecordType: "A",
				Targets:    localTargets,
			}
			gslbHosts = append(gslbHosts, dnsRecord)
		}

		// Check if host is alive on external Gslb
		externalTargets, err := getExternalTargets(host)
		if err != nil {
			return nil, err
		}
		if len(externalTargets) > 0 {
			switch gslb.Spec.Strategy {
			case "roundRobin":
				finalTargets = append(finalTargets, externalTargets...)
			}
		}

		if len(finalTargets) > 0 {
			dnsRecord := &externaldns.Endpoint{
				DNSName:    host,
				RecordTTL:  30,
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
			Annotations: map[string]string{"ohmyglb.absa.oss/dnstype": "local"},
		},
		Spec: dnsEndpointSpec,
	}

	err = controllerutil.SetControllerReference(gslb, dnsEndpoint, r.scheme)
	if err != nil {
		return nil, err
	}
	return dnsEndpoint, err
}

func nsServerName(gslb *ohmyglbv1beta1.Gslb) string {
	clusterGeoTag := os.Getenv("CLUSTER_GEO_TAG")
	edgeDNSZone := os.Getenv("EDGE_DNS_ZONE")
	if len(clusterGeoTag) == 0 {
		clusterGeoTag = "default"
	}
	return fmt.Sprintf("%s-ns-%s.%s", gslb.Name, clusterGeoTag, edgeDNSZone)
}

func (r *ReconcileGslb) gslbEdgeDNSEndpoint(gslb *ohmyglbv1beta1.Gslb) (*externaldns.DNSEndpoint, error) {
	edgeDNSEndpointSpec := externaldns.DNSEndpointSpec{}
	localTargets, err := r.getGslbIngressIPs(gslb)
	if err != nil {
		return nil, err
	}
	// Type A record to be registered resolve NS records in edge dns responses(infoblox, route53,...)
	edgeDNSRecord := &externaldns.Endpoint{
		DNSName:    nsServerName(gslb),
		RecordTTL:  30,
		RecordType: "A",
		Targets:    localTargets,
	}
	edgeDNSEndpointSpec = externaldns.DNSEndpointSpec{
		Endpoints: []*externaldns.Endpoint{
			edgeDNSRecord,
		},
	}

	edgeDNSEndpoint := &externaldns.DNSEndpoint{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-ns", gslb.Name),
			Namespace:   gslb.Namespace,
			Annotations: map[string]string{"ohmyglb.absa.oss/dnstype": "edgedns"},
		},
		Spec: edgeDNSEndpointSpec,
	}
	err = controllerutil.SetControllerReference(gslb, edgeDNSEndpoint, r.scheme)
	if err != nil {
		return nil, err
	}
	return edgeDNSEndpoint, nil
}

func infobloxConnection() (*ibclient.ObjectManager, error) {
	hostConfig := ibclient.HostConfig{
		Host:     os.Getenv("INFOBLOX_GRID_HOST"),
		Version:  os.Getenv("INFOBLOX_WAPI_VERSION"),
		Port:     os.Getenv("INFOBLIX_WAPI_PORT"),
		Username: os.Getenv("EXTERNAL_DNS_INFOBLOX_WAPI_USERNAME"),
		Password: os.Getenv("EXTERNAL_DNS_INFOBLOX_WAPI_PASSWORD"),
	}
	transportConfig := ibclient.NewTransportConfig("false", 20, 10)
	requestBuilder := &ibclient.WapiRequestBuilder{}
	requestor := &ibclient.WapiHttpRequestor{}
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

	objMgr := ibclient.NewObjectManager(conn, "ohmyclient", "")
	return objMgr, nil
}

func (r *ReconcileGslb) configureZoneDelegation(gslb *ohmyglbv1beta1.Gslb) (*reconcile.Result, error) {
	if len(os.Getenv("INFOBLOX_GRID_HOST")) > 0 {

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
			nameServer := ibclient.NameServer{Address: address, Name: nsServerName(gslb)}
			delegateTo = append(delegateTo, nameServer)
		}
		gslbZoneName := os.Getenv("DNS_ZONE")
		findZone, err := objMgr.GetZoneDelegated(gslbZoneName)
		if err != nil {
			return &reconcile.Result{}, err
		}

		if len(findZone.Ref) > 0 {
			log.Info(fmt.Sprintf("Updating delegated zone(%s)...", gslbZoneName))
			_, err = objMgr.UpdateZoneDelegated(findZone.Ref, delegateTo)
		} else {
			log.Info(fmt.Sprintf("Creating delegated zone(%s)...", gslbZoneName))
			_, err = objMgr.CreateZoneDelegated(gslbZoneName, delegateTo)
		}
		if err != nil {
			return &reconcile.Result{}, err
		}
		return &reconcile.Result{}, err
	}
	return nil, nil
}

func (r *ReconcileGslb) ensureDNSEndpoint(request reconcile.Request,
	gslb *ohmyglbv1beta1.Gslb,
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

func prettyPrint(s interface{}) string {
	prettyStruct, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		fmt.Println("can't convert struct to json")
	}
	return string(prettyStruct)
}
