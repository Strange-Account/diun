package notif

import (
	"github.com/Strange-Account/diun/internal/notif/script"
	"github.com/Strange-Account/diun/internal/model"
	"github.com/Strange-Account/diun/internal/notif/mail"
	"github.com/Strange-Account/diun/internal/notif/notifier"
	"github.com/Strange-Account/diun/internal/notif/slack"
	"github.com/Strange-Account/diun/internal/notif/telegram"
	"github.com/Strange-Account/diun/internal/notif/webhook"
	"github.com/rs/zerolog/log"
)

// Client represents an active webhook notification object
type Client struct {
	cfg       model.Notif
	app       model.App
	notifiers []notifier.Notifier
}

// New creates a new notification instance
func New(config model.Notif, app model.App) (*Client, error) {
	var c = &Client{
		cfg:       config,
		app:       app,
		notifiers: []notifier.Notifier{},
	}

	// Add notifiers
	if config.Mail.Enable {
		c.notifiers = append(c.notifiers, mail.New(config.Mail, app))
	}
	if config.Slack.Enable {
		c.notifiers = append(c.notifiers, slack.New(config.Slack, app))
	}
	if config.Telegram.Enable {
		c.notifiers = append(c.notifiers, telegram.New(config.Telegram, app))
	}
	if config.Webhook.Enable {
		c.notifiers = append(c.notifiers, webhook.New(config.Webhook, app))
	}
	if config.Script.Enable {
		c.notifiers = append(c.notifiers, script.New(config.Script, app))
	}

	log.Debug().Msgf("%d notifier(s) created", len(c.notifiers))
	return c, nil
}

// Send creates and sends notifications to notifiers
func (c *Client) Send(entry model.NotifEntry) {
	for _, n := range c.notifiers {
		log.Debug().Str("image", entry.Image.String()).Msgf("Sending %s notification...", n.Name())
		if err := n.Send(entry); err != nil {
			log.Error().Err(err).Str("image", entry.Image.String()).Msgf("%s notification failed", n.Name())
		}
	}
}
