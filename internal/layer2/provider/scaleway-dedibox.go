package provider

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	//ScalewayDediboxAPIURL scaleway-dedibox api url
	ScalewayDediboxAPIURL = "https://api.online.net/api/v1"
)

// ScalewayDediboxClient client
type ScalewayDediboxClient struct {
	genericClient
}

// ScalewayDediboxIPInfo struct
type ScalewayDediboxIPInfo struct {
	Address     string `json:"address"`
	Destination string `json:"destination"`
}

// SetIP change ip destination
func (c *ScalewayDediboxClient) SetIP(ip net.IP, destination string) error {

	// Prepare request
	request := gorequest.New()

	// Check ip is not already on destination
	var ipInfoRes ScalewayDediboxIPInfo

	resp, bodyByte, errs := c.setRequestOptions(request.Get(ScalewayDediboxAPIURL + "/server/ip/" + ip.String())).
		EndStruct(&ipInfoRes)

	if errs != nil {
		return fmt.Errorf("Failed to get ip %s information: %s", ip.String(), c.joinErrors(errs))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to get ip %s information: %s", ip.String(), errors.New(string(bodyByte)))
	}

	if ipInfoRes.Destination == destination {
		c.logger.Log("op", "setIPDestinationOnSoyoustart", "info", fmt.Sprintf("Ip %s already set on destination %s", ip.String(), destination))
		return nil
	}

	// Set ip for destination
	request = gorequest.New()

	ipEditParams := &ScalewayDediboxIPInfo{
		Address:     ip.String(),
		Destination: destination,
	}

	resp, body, errs := c.setRequestOptions(request.Post(ScalewayDediboxAPIURL + "/server/ip/edit")).
		Send(ipEditParams).
		End()

	if errs != nil {
		return fmt.Errorf("Failed to set ip %s on destination %s: %s", ip.String(), destination, c.joinErrors(errs))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to set ip %s on destination %s: %s", ip.String(), destination, errors.New(string(body)))
	}

	return nil
}

func (c *ScalewayDediboxClient) setRequestOptions(request *gorequest.SuperAgent) *gorequest.SuperAgent {
	return request.Timeout(180*time.Second).
		Retry(3, 1*time.Second, http.StatusInternalServerError, http.StatusNotFound).
		Set("Authorization", fmt.Sprintf("Bearer %s", c.auth.Token)).
		Set("Accept", "application/json")
}

func (c *ScalewayDediboxClient) joinErrors(errs []error) error {

	s := ""

	for _, err := range errs {
		s += err.Error()
	}

	return errors.New(s)
}
