package swarm

import (
	"fmt"

	"github.com/Strange-Account/diun/internal/model"
	"github.com/Strange-Account/diun/internal/provider"
	"github.com/rs/zerolog/log"
)

// Client represents an active swarm provider object
type Client struct {
	*provider.Client
	elts map[string]model.PrdSwarm
}

// New creates new swarm provider instance
func New(elts map[string]model.PrdSwarm) *provider.Client {
	return &provider.Client{Handler: &Client{
		elts: elts,
	}}
}

// ListJob returns job list to process
func (c *Client) ListJob() []model.Job {
	if len(c.elts) == 0 {
		return []model.Job{}
	}

	log.Info().Msgf("Found %d swarm provider(s) to analyze...", len(c.elts))
	var list []model.Job
	for id, elt := range c.elts {
		for _, img := range c.listServiceImage(id, elt) {
			list = append(list, model.Job{
				Provider: fmt.Sprintf("swarm-%s", id),
				Image:    img,
			})
		}
	}

	return list
}
