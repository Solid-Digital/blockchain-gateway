package pipeline

import (
	"bitbucket.org/unchain/ares/pkg/3p/elastic"
	"bitbucket.org/unchain/ares/pkg/3p/harbor"
	"bitbucket.org/unchain/ares/pkg/3p/sql"
	"github.com/google/wire"

	"github.com/unchainio/interfaces/logger"

	"bitbucket.org/unchain/ares/pkg/ares"
)

var ServiceSet = wire.NewSet(NewService, wire.Bind(new(ares.PipelineService), new(Service)))

var _ ares.PipelineService = &Service{}

type Service struct {
	db      *sql.DB
	service *service
	cfg     *Config
	log     logger.Logger
}

type service struct {
	kube         ares.DeploymentService
	containers   ares.ContainerService
	imageBuilder ares.ImageBuilder
	elastic      *elastic.Client
	registry     *harbor.Client
	store        ares.FileStore
}

func NewService(
	db *sql.DB,
	store ares.FileStore,
	imageBuilder ares.ImageBuilder,
	containerService ares.ContainerService,
	kubeService ares.DeploymentService,
	harborClient *harbor.Client,
	elasticClient *elastic.Client,
	cfg *Config,
	log logger.Logger,
) (adapterService *Service) {
	return &Service{
		db: db,
		service: &service{
			kube:         kubeService,
			containers:   containerService,
			imageBuilder: imageBuilder,
			elastic:      elasticClient,
			registry:     harborClient,
			store:        store,
		},
		cfg: cfg,
		log: log,
	}
}
