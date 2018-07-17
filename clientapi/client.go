package clientapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/shu_go/browser-notifier"
)

type Client struct {
	Server string
	Port   int

	App string
}

func New(server string, port int, app string) *Client {
	return &Client{Server: server, Port: port, App: app}
}

func (c *Client) Notify(n notifier.Notification) error {
	if n.App == "" {
		n.App = c.App
	}

	n.App = strings.Replace(n.App, `\n`, `\\n`, -1)
	n.Title = strings.Replace(n.Title, `\n`, `\\n`, -1)
	n.Text = strings.Replace(n.Text, `\n`, `\\n`, -1)

	b, err := json.Marshal(n)
	if err != nil {
		return err
	}

	data := bytes.NewBuffer(b)

	resp, err := http.Post(fmt.Sprintf("http://%s:%d/notifications", c.Server, c.Port), "application/json", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
