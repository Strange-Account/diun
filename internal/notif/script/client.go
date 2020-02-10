package script

import (
	"fmt"
	"log"
	"os/exec"

	//	"fmt"
	//	"net/http"

	"github.com/crazy-max/diun/internal/model"
	"github.com/crazy-max/diun/internal/notif/notifier"
)

// Client represents an active slack notification object
type Client struct {
	*notifier.Notifier
	cfg model.NotifScript
	app model.App
}

// New creates a new slack notification instance
func New(config model.NotifScript, app model.App) notifier.Notifier {
	return notifier.Notifier{
		Handler: &Client{
			cfg: config,
			app: app,
		},
	}
}

// Name returns notifier's name
func (c *Client) Name() string {
	return "script"
}

// Send creates and sends a webhook notification with an entry
func (c *Client) Send(entry model.NotifEntry) error {

	// body, err := json.Marshal(struct {
	// 	Version      string        `json:"diun_version"`
	// 	Status       string        `json:"status"`
	// 	Provider     string        `json:"provider"`
	// 	Image        string        `json:"image"`
	// 	MIMEType     string        `json:"mime_type"`
	// 	Digest       digest.Digest `json:"digest"`
	// 	Created      *time.Time    `json:"created"`
	// 	Architecture string        `json:"architecture"`
	// 	Os           string        `json:"os"`
	// }{
	// 	Version:      c.app.Version,
	// 	Status:       string(entry.Status),
	// 	Provider:     entry.Provider,
	// 	Image:        entry.Image.String(),
	// 	MIMEType:     entry.Manifest.MIMEType,
	// 	Digest:       entry.Manifest.Digest,
	// 	Created:      entry.Manifest.Created,
	// 	Architecture: entry.Manifest.Architecture,
	// 	Os:           entry.Manifest.Os,
	// })
	// if err != nil {
	// 	return err
	// }

	out, err := exec.Command(c.cfg.Endpoint, entry.Image.String()).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("output is %s\n", out)

	// req, err := http.NewRequest(c.cfg.Method, c.cfg.Endpoint, bytes.NewBuffer([]byte(body)))
	// if err != nil {
	// 	return err
	// }

	// if len(c.cfg.Headers) > 0 {
	// 	for key, value := range c.cfg.Headers {
	// 		req.Header.Add(key, value)
	// 	}
	// }

	// req.Header.Set("User-Agent", fmt.Sprintf("%s %s", c.app.Name, c.app.Version))

	// _, err = hc.Do(req)
	return err
}
