package kubernetes

import (
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Config struct {
	RegistryCredentialsSecret string
	Config                    *clientcmdapi.Config
	ConfigPath                string
	RancherProjectID          string
	RancherClusterID          string
	TLS                       bool
}
