package dns

import (
	"fmt"
	"time"

	externaldns "sigs.k8s.io/external-dns/endpoint"

	"github.com/AbsaOSS/k8gb/controllers/providers/assistant"

	k8gbv1beta1 "github.com/AbsaOSS/k8gb/api/v1beta1"
	"github.com/AbsaOSS/k8gb/controllers/depresolver"
	ibclient "github.com/infobloxopen/infoblox-go-client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type InfobloxProvider struct {
	assistant assistant.IAssistant
	config    depresolver.Config
}

func NewInfobloxDNS(config depresolver.Config, assistant assistant.IAssistant) *InfobloxProvider {
	return &InfobloxProvider{
		assistant: assistant,
		config:    config,
	}
}

func (p *InfobloxProvider) CreateZoneDelegationForExternalDNS(gslb *k8gbv1beta1.Gslb) (*reconcile.Result, error) {
	objMgr, err := p.infobloxConnection()
	if err != nil {
		return &reconcile.Result{}, err
	}
	addresses, err := p.assistant.GslbIngressExposedIPs(gslb)
	if err != nil {
		return &reconcile.Result{}, err
	}
	var delegateTo []ibclient.NameServer

	for _, address := range addresses {
		nameServer := ibclient.NameServer{Address: address, Name: nsServerName(p.config)}
		delegateTo = append(delegateTo, nameServer)
	}

	findZone, err := objMgr.GetZoneDelegated(p.config.DNSZone)
	if err != nil {
		return &reconcile.Result{}, err
	}

	if findZone != nil {
		err = p.checkZoneDelegated(findZone)
		if err != nil {
			return &reconcile.Result{}, err
		}
		if len(findZone.Ref) > 0 {

			// Drop own records for straight away update
			existingDelegateTo := p.filterOutDelegateTo(findZone.DelegateTo, nsServerName(p.config))
			existingDelegateTo = append(existingDelegateTo, delegateTo...)

			// Drop external records if they are stale
			extClusters := getExternalClusterHeartbeatFQDNs(gslb, p.config)
			for _, extCluster := range extClusters {
				err = p.assistant.InspectTXTThreshold(
					extCluster,
					p.config.Override.FakeDNSEnabled,
					time.Second*time.Duration(gslb.Spec.Strategy.SplitBrainThresholdSeconds))
				if err != nil {
					p.assistant.Error(err, "got the error from TXT based checkAlive")
					p.assistant.Info("External cluster (%s) doesn't look alive, filtering it out from "+
						"delegated zone configuration...", extCluster)
					existingDelegateTo = p.filterOutDelegateTo(existingDelegateTo, extCluster)
				}
			}
			p.assistant.Info("Updating delegated zone(%s) with the server list(%v)", p.config.DNSZone, existingDelegateTo)

			_, err = objMgr.UpdateZoneDelegated(findZone.Ref, existingDelegateTo)
			if err != nil {
				return &reconcile.Result{}, err
			}
		}
	} else {
		p.assistant.Info("Creating delegated zone(%s)...", p.config.DNSZone)
		_, err = objMgr.CreateZoneDelegated(p.config.DNSZone, delegateTo)
		if err != nil {
			return &reconcile.Result{}, err
		}
	}

	edgeTimestamp := fmt.Sprint(time.Now().UTC().Format("2006-01-02T15:04:05"))
	heartbeatTXTName := fmt.Sprintf("%s-heartbeat-%s.%s", gslb.Name, p.config.ClusterGeoTag, p.config.EdgeDNSZone)
	heartbeatTXTRecord, err := objMgr.GetTXTRecord(heartbeatTXTName)
	if err != nil {
		return &reconcile.Result{}, err
	}
	if heartbeatTXTRecord == nil {
		p.assistant.Info("Creating split brain TXT record(%s)...", heartbeatTXTName)
		_, err := objMgr.CreateTXTRecord(heartbeatTXTName, edgeTimestamp, gslb.Spec.Strategy.DNSTtlSeconds, "default")
		if err != nil {
			return &reconcile.Result{}, err
		}
	} else {
		p.assistant.Info("Updating split brain TXT record(%s)...", heartbeatTXTName)
		_, err := objMgr.UpdateTXTRecord(heartbeatTXTName, edgeTimestamp)
		if err != nil {
			return &reconcile.Result{}, err
		}
	}
	return nil, nil
}

func (p *InfobloxProvider) Finalize(gslb *k8gbv1beta1.Gslb) error {
	objMgr, err := p.infobloxConnection()
	if err != nil {
		return err
	}
	findZone, err := objMgr.GetZoneDelegated(p.config.DNSZone)
	if err != nil {
		return err
	}

	if findZone != nil {
		err = p.checkZoneDelegated(findZone)
		if err != nil {
			return err
		}
		if len(findZone.Ref) > 0 {
			p.assistant.Info("Deleting delegated zone(%s)...", p.config.DNSZone)
			_, err := objMgr.DeleteZoneDelegated(findZone.Ref)
			if err != nil {
				return err
			}
		}
	}

	heartbeatTXTName := fmt.Sprintf("%s-heartbeat-%s.%s", gslb.Name, p.config.ClusterGeoTag, p.config.EdgeDNSZone)
	findTXT, err := objMgr.GetTXTRecord(heartbeatTXTName)
	if err != nil {
		return err
	}

	if findTXT != nil {
		if len(findTXT.Ref) > 0 {
			p.assistant.Info("Deleting split brain TXT record(%s)...", heartbeatTXTName)
			_, err := objMgr.DeleteTXTRecord(findTXT.Ref)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *InfobloxProvider) GetExternalTargets(host string) (targets []string) {
	return p.assistant.GetExternalTargets(host, p.config.Override.FakeDNSEnabled, nsServerNameExt(p.config))
}

func (p *InfobloxProvider) GslbIngressExposedIPs(gslb *k8gbv1beta1.Gslb) ([]string, error) {
	return p.assistant.GslbIngressExposedIPs(gslb)
}

func (p *InfobloxProvider) SaveDNSEndpoint(gslb *k8gbv1beta1.Gslb, i *externaldns.DNSEndpoint) (*reconcile.Result, error) {
	return p.assistant.SaveDNSEndpoint(gslb.Namespace, i)
}

func (p *InfobloxProvider) String() string {
	return "Infoblox"
}
