package docker

import (
	"context"
	"io"

	"encoding/base64"
	"encoding/json"

	"bitbucket.org/unchain/ares/pkg/ares"
	"github.com/docker/docker/api/types"
	dc "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/google/wire"
	"github.com/unchainio/interfaces/logger"
	"github.com/unchainio/pkg/errors"
)

var ContainerServiceSet = wire.NewSet(NewClient, NewContainerService, wire.Bind(new(ares.ContainerService), new(ContainerService)))

const defaultWorkingDir = "/opt/unchain/bin"

type ContainerService struct {
	client *Client

	log logger.Logger
	cfg *Config
}

// Wrapper around the docker client, must have no other direct dependencies
func NewContainerService(
	client *Client,
	log logger.Logger,
	cfg *Config,
) *ContainerService {
	return &ContainerService{
		client: client,
		log:    log,
		cfg:    cfg,
	}
}

// creates a new container, copies the artifacts into it and runs it
func (s *ContainerService) PrepareImage(imageRef string, baseImageRef string, cmd string, tarredArtifacts ...io.Reader) error {
	s.log.Debugf("Preparing docker image `%s`", imageRef)
	ctx := context.Background()

	authToken, err := getAuthToken(s.cfg.Auth)
	if err != nil {
		return err
	}

	err = s.fetchAdapterBaseImage(baseImageRef, authToken)
	if err != nil {
		return err
	}

	containerID, containerCfg, err := s.createContainer(ctx, baseImageRef, cmd)
	if err != nil {
		return err
	}

	err = s.copyToContainer(ctx, containerID, tarredArtifacts)
	if err != nil {
		return err
	}

	s.log.Debugf("Committing container")
	_, err = s.commitContainer(containerID, imageRef, containerCfg)
	if err != nil {
		return err
	}
	s.log.Debugf("Committing container...Done")

	s.log.Debugf("Pushing container to registry")
	err = s.pushImageToRegistry(imageRef, authToken)
	if err != nil {
		return err
	}
	s.log.Debugf("Pushing container to registry...Done")

	err = s.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	_, err = s.client.ImageRemove(ctx, imageRef, types.ImageRemoveOptions{})
	if err != nil {
		return err
	}

	return nil
}

//  create container
// setup logging driver for the container
func (s *ContainerService) createContainer(ctx context.Context, baseImageRef string, cmd string) (string, *dc.Config, error) {
	s.log.Debugf("Creating container\n")

	ports, err := getExposedPorts()

	if err != nil {
		return "", nil, err
	}

	var cfg = &dc.Config{
		Image:        baseImageRef,
		ExposedPorts: ports,
		Cmd:          []string{"/bin/sh", "-c", cmd},
		WorkingDir:   defaultWorkingDir,
		AttachStdout: true,
	}

	netConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{},
	}

	container, err := s.client.ContainerCreate(ctx, cfg, nil, netConfig, "")

	if err != nil {
		return "", nil, errors.Wrap(err, "Error creating container: ")
	}

	s.log.Debugf("Done creating image")

	return container.ID, cfg, nil
}

func (s *ContainerService) commitContainer(containerID string, imageName string, config *dc.Config) (string, error) {
	s.log.Debugf("Committing container to image")

	ctx := context.Background()

	// Pull registry stuff into a config
	options := types.ContainerCommitOptions{
		Config:    config,
		Author:    "unchain",
		Reference: imageName,
	}
	resp, err := s.client.ContainerCommit(ctx, containerID, options)
	if err != nil {
		return resp.ID, errors.Wrap(err, "Error committing containerID: ")
	}
	s.log.Debugf(resp.ID)

	return resp.ID, nil
}

func (s *ContainerService) pushImageToRegistry(ref string, authToken string) error {
	s.log.Debugf("Pushing image to registry")

	ctx := context.Background()
	_, err := s.client.ImagePush(ctx, ref, types.ImagePushOptions{
		RegistryAuth: authToken,
	})
	if err != nil {
		return errors.Wrap(err, "Error pushing image to registry: ")
	}

	return nil
}

func getExposedPorts() (nat.PortSet, error) {
	ports := make(nat.PortSet)

	httpsPort, err := nat.NewPort("tcp", "80")

	if err != nil {
		return nil, errors.Wrap(err, "Error creating container: ")
	}

	ports[httpsPort] = struct{}{}

	return ports, nil
}

// Copy artifacts into container
func (s *ContainerService) copyToContainer(ctx context.Context, containerID string, tarredArtifacts []io.Reader) error {
	s.log.Debugf("Copying artifacts to container")

	options := types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: false,
	}

	for _, artifact := range tarredArtifacts {
		err := s.client.CopyToContainer(ctx, containerID, defaultWorkingDir, artifact, options)

		if err != nil {
			return errors.Wrap(err, "failed to copy to container")
		}
	}
	s.log.Debugf("Copying artifacts to container...Done")

	return nil
}

func (s *ContainerService) fetchAdapterBaseImage(baseImageRef string, authToken string) error {
	var err error

	_, err = s.client.ImagePull(context.Background(), baseImageRef, types.ImagePullOptions{
		RegistryAuth: authToken,
	})

	if err != nil {
		return errors.Wrap(err, "Error creating Janus base image")
	}

	return nil
}

const emptyAuthToken = "emptyAuthToken"

func getAuthToken(authConfig *types.AuthConfig) (string, error) {
	if authConfig != nil {
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return "", err
		}
		return base64.URLEncoding.EncodeToString(encodedJSON), nil
	}

	// This is needed because the auth token must always be set to something even if it isn't used.
	return emptyAuthToken, nil
}
