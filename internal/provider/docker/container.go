package docker

import (
	"fmt"
	"reflect"

	"github.com/Strange-Account/diun/internal/model"
	"github.com/Strange-Account/diun/internal/provider"
	"github.com/Strange-Account/diun/pkg/docker"
	"github.com/docker/docker/api/types/filters"
	"github.com/rs/zerolog/log"
)

func (c *Client) listContainerImage(id string, elt model.PrdDocker) []model.Image {
	sublog := log.With().
		Str("provider", fmt.Sprintf("docker-%s", id)).
		Logger()

	cli, err := docker.New(docker.Options{
		Endpoint:    elt.Endpoint,
		APIVersion:  elt.APIVersion,
		TLSCertPath: elt.TLSCertsPath,
		TLSVerify:   elt.TLSVerify,
	})
	if err != nil {
		sublog.Error().Err(err).Msg("Cannot create Docker client")
		return []model.Image{}
	}

	ctnFilter := filters.NewArgs()
	ctnFilter.Add("status", "running")
	if elt.WatchStopped {
		ctnFilter.Add("status", "created")
		ctnFilter.Add("status", "exited")
	}

	ctns, err := cli.ContainerList(ctnFilter)
	if err != nil {
		sublog.Error().Err(err).Msg("Cannot list Docker containers")
		return []model.Image{}
	}

	var list []model.Image
	for _, ctn := range ctns {
		local, err := cli.IsLocalImage(ctn.Image)
		if err != nil {
			sublog.Error().Err(err).Msgf("Cannot inspect image from container %s", ctn.ID)
			continue
		} else if local {
			sublog.Debug().Msgf("Skip locally built image for container %s", ctn.ID)
			continue
		}
		image, err := provider.ValidateContainerImage(ctn.Image, ctn.Labels, elt.WatchByDefault)
		if err != nil {
			sublog.Error().Err(err).Msgf("Cannot get image from container %s", ctn.ID)
			continue
		} else if reflect.DeepEqual(image, model.Image{}) {
			sublog.Debug().Msgf("Watch disabled for container %s", ctn.ID)
			continue
		}
		list = append(list, image)
	}

	return list
}
