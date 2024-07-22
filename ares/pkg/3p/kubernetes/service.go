package kubernetes

import (
	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"
	"k8s.io/api/core/v1"
	errorv1 "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var ServiceSet = wire.NewSet(NewService, wire.Bind(new(ares.DeploymentService), new(Service)))

type Service struct {
	Client *kubernetes.Clientset

	cfg *Config
	log logger.Logger
}

func (s *Service) upsertSecret(orgName string, secret *v1.Secret) error {
	c := s.Client.CoreV1().Secrets(orgName)

	res, err := c.Update(secret)
	switch {
	case err == nil:
		s.log.Debugf("updated secret %q.\n", res.GetObjectMeta().GetName())
	case !errorv1.IsNotFound(err):
		return errors.Wrap(err, "could not update secret: ")
	default:
		res, err = c.Create(secret)
		if err != nil {
			return errors.Wrap(err, "could not create secret: ")
		}
		s.log.Debugf("secret created")
	}

	return nil
}

func NewService(cfg *Config, log logger.Logger) (*Service, error) {
	apicfg, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
		// If ConfigPath is not specified, return the inlined config
		if cfg.ConfigPath == "" {
			return cfg.Config, nil
		}

		// If ConfigPath is specified, use it to load the config
		kubecfg, err := clientcmd.LoadFromFile(cfg.ConfigPath)

		if err != nil {
			return nil, errors.Wrapf(err, "Failed to parse kube kubecfg")
		}

		return kubecfg, nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create kube api config")
	}

	clientset, err := kubernetes.NewForConfig(apicfg)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create kube client set")
	}

	return &Service{
		Client: clientset,
		cfg:    cfg,
		log:    log,
	}, nil
}
