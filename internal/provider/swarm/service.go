package swarm

import (
	"fmt"
	"reflect"

	"github.com/Strange-Account/diun/internal/model"
	"github.com/Strange-Account/diun/internal/provider"
	"github.com/Strange-Account/diun/pkg/docker"
	"github.com/docker/docker/api/types/filters"
	"github.com/rs/zerolog/log"
)

func (c *Client) listServiceImage(id string, elt model.PrdSwarm) []model.Image {
	sublog := log.With().
		Str("provider", fmt.Sprintf("swarm-%s", id)).
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

	svcs, err := cli.ServiceList(filters.NewArgs())
	if err != nil {
		sublog.Error().Err(err).Msg("Cannot list Swarm services")
		return []model.Image{}
	}

	var list []model.Image
	for _, svc := range svcs {
		local, err := cli.IsLocalImage(svc.Spec.TaskTemplate.ContainerSpec.Image)
		if err != nil {
			sublog.Error().Err(err).Msgf("Cannot inspect image from service %s", svc.ID)
			continue
		} else if local {
			sublog.Debug().Msgf("Skip locally built image for service %s", svc.ID)
			continue
		}
		image, err := provider.ValidateContainerImage(svc.Spec.TaskTemplate.ContainerSpec.Image, svc.Spec.Labels, elt.WatchByDefault)
		if err != nil {
			sublog.Error().Err(err).Msgf("Cannot get image from service %s", svc.ID)
			continue
		} else if reflect.DeepEqual(image, model.Image{}) {
			sublog.Debug().Msgf("Watch disabled for service %s", svc.ID)
			continue
		}
		list = append(list, image)
	}

	return list
}
