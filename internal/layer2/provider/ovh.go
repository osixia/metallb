package provider

import (
	"fmt"
	"net"

	"github.com/ovh/go-ovh/ovh"
)

// OvhClient client
type OvhClient struct {
	genericClient
	endpoint string
}

// OvhIPInfoRoutedTo struct
type OvhIPInfoRoutedTo struct {
	ServiceName string `json:"serviceName"`
}

// OvhIPInfo struct
type OvhIPInfo struct {
	RoutedTo OvhIPInfoRoutedTo `json:"routedTo"`
}

// OvhIPMove struct
type OvhIPMove struct {
	To string `json:"to"`
}

// SetIP change ip destination
func (c *OvhClient) SetIP(ip net.IP, destination string) error {

	// Create client
	client, err := ovh.NewClient(
		c.endpoint,
		c.auth.ApplicationKey,
		c.auth.ApplicationSecret,
		c.auth.ConsumerKey,
	)
	if err != nil {
		return fmt.Errorf("Failed to set ip %s on destination %s: %s", ip.String(), destination, err.Error())
	}

	// Check ip is not already on destination
	var ipInfoRes OvhIPInfo

	if err := client.Get(fmt.Sprintf("/ip/%s", ip.String()), &ipInfoRes); err != nil {
		return fmt.Errorf("Failed to get ip %s information: %s", ip.String(), err.Error())
	}

	if ipInfoRes.RoutedTo.ServiceName == destination {
		c.logger.Log("op", "setIPDestinationOnSoyoustart", "info", fmt.Sprintf("Ip %s already set on destination %s", ip.String(), destination))
		return nil
	}

	c.logger.Log("op", "setIPDestinationOnSoyoustart", "ipInfoRes.routedTo.serviceName", ipInfoRes.RoutedTo.ServiceName, "destination", destination)

	ipMoveParams := &OvhIPMove{
		To: destination,
	}

	if err := client.Post(fmt.Sprintf("/ip/%s/move", ip.String()), ipMoveParams, nil); err != nil {
		return fmt.Errorf("Failed to set ip %s on destination %s: %s", ip.String(), destination, err.Error())
	}

	return nil
}
