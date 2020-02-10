package static

import (
	"fmt"

	"github.com/Strange-Account/diun/internal/model"
	"github.com/Strange-Account/diun/internal/provider"
	"github.com/rs/zerolog/log"
)

// Client represents an active static provider object
type Client struct {
	*provider.Client
	elts []model.PrdStatic
}

// New creates new static provider instance
func New(elts []model.PrdStatic) *provider.Client {
	return &provider.Client{Handler: &Client{
		elts: elts,
	}}
}

// ListJob returns job list to process
func (c *Client) ListJob() []model.Job {
	if len(c.elts) == 0 {
		return []model.Job{}
	}

	log.Info().Msgf("Found %d static provider(s) to analyze...", len(c.elts))
	var list []model.Job
	for key, elt := range c.elts {
		list = append(list, model.Job{
			Provider: fmt.Sprintf("static-%d", key),
			Image:    model.Image(elt),
		})
	}

	return list
}
