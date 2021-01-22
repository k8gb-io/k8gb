package dns

import (
	ibclient "github.com/infobloxopen/infoblox-go-client"
)

type fakeInfobloxConnector struct {
	// createObjectObj interface{}

	getObjectObj interface{}
	getObjectRef string

	// deleteObjectRef string

	// updateObjectObj interface{}
	// updateObjectRef string

	resultObject interface{}

	fakeRefReturn string
}

func (c *fakeInfobloxConnector) CreateObject(ibclient.IBObject) (string, error) {
	return c.fakeRefReturn, nil
}

func (c *fakeInfobloxConnector) GetObject(ibclient.IBObject, string, interface{}) (err error) {
	return nil
}

func (c *fakeInfobloxConnector) DeleteObject(string) (string, error) {
	return c.fakeRefReturn, nil
}

func (c *fakeInfobloxConnector) UpdateObject(ibclient.IBObject, string) (string, error) {
	return c.fakeRefReturn, nil
}
