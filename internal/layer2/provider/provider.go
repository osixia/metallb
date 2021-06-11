package provider

import (
	"fmt"
	"net"

	"github.com/go-kit/kit/log"
)

// Provider holds the ip providers we support.
type Provider string

// MetalLB supported ip providers.
const (
	SclalewayDedibox Provider = "scaleway-dedibox"
	Soyoustart                = "soyoustart"
	Ovh                       = "ovh"
	Kimsufi                   = "kimsufi"
)

// Auth holds ip provider auth configuration.
type Auth struct {
	Token             string `yaml:"token"`
	ApplicationKey    string `yaml:"application-key"`
	ApplicationSecret string `yaml:"application-secret"`
	ConsumerKey       string `yaml:"consumer-key"`
}

// Client can change ip addresse  destination on a provider an IP address.
type Client interface {
	SetIP(net.IP, string) error
}

type genericClient struct {
	logger log.Logger
	auth   *Auth
}

// GetClient returns provider client
func GetClient(logger log.Logger, provider Provider, auth *Auth) (Client, error) {

	var handler Client

	genericClient := genericClient{
		logger: logger,
		auth:   auth,
	}

	switch provider {
	case Ovh:
		handler = &OvhClient{
			genericClient: genericClient,
			endpoint:      "ovh-eu",
		}
	case Soyoustart:
		handler = &OvhClient{
			genericClient: genericClient,
			endpoint:      "soyoustart-eu",
		}
	case Kimsufi:
		handler = &OvhClient{
			genericClient: genericClient,
			endpoint:      "kimsufi-eu",
		}
	case SclalewayDedibox:
		handler = &ScalewayDediboxClient{
			genericClient: genericClient,
		}
	}

	if handler == nil {
		return nil, fmt.Errorf("No client for provider : %s", provider)
	}

	return handler, nil
}
