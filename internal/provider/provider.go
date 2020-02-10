package provider

import (
	"github.com/Strange-Account/diun/internal/model"
)

// Handler is a provider interface
type Handler interface {
	ListJob() []model.Job
	Close() error
}

// Client represents an active provider object
type Client struct {
	Handler
}
