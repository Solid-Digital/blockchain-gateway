package ares

import (
	"io"
)

// BuildManifest is the object defining what to build.
type BuildManifest struct {
	Components []*Component `json:"components,omitempty"`
	Tag        string       `json:"tag,omitempty"`
	BaseImage  string       `json:"baseImage,omitempty"`
	Cmd        []string     `json:"cmd,omitempty"`
}

// Component defines a single component filename and location.
type Component struct {
	FileName string
	FileId   string
}

// Wrapper around a container provider such as docker or another one, must have no direct dependencies other than the provider it wraps
type ContainerService interface {
	PrepareImage(imageName string, baseImageRef string, cmd string, tarredArtifacts ...io.Reader) (err error)
}

// Wrapper around a container provider such as docker or another one, must have no direct dependencies other than the provider it wraps
type ImageBuilder interface {
	BuildImage(manifest *BuildManifest) error
}
