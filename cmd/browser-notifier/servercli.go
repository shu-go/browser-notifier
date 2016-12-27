package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/pkg/browser"
	"github.com/urfave/cli"

	"bitbucket.org/shu/browser-notifier"
)

func main() {
	app := cli.NewApp()
	app.Name = "Browser Notifier"
	app.Usage = "A notifier using Web Notifications."
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "port, p", Value: "7878", Usage: "server `PORT`"},
		cli.StringFlag{Name: "storage, s", Value: "", Usage: "the `FILE` for storageence"},
	}
	app.Action = func(c *cli.Context) error {
		port := c.Int("port")
		storage := c.String("storage")

		notifications := []notifier.Notification{}
		if storage != "" {
			if f, err := os.Open(storage); err == nil {
				defer f.Close()
				if b, err := ioutil.ReadAll(f); err == nil {
					if err = json.Unmarshal(b, &notifications); err != nil {
						notifications = nil
					}
				}
			}
		}

		if _, err := os.Stat("./cmd/assets"); err == nil {
			http.Handle("/", http.FileServer(http.Dir("./cmd/assets")))
		} else {
			http.Handle("/", http.FileServer(assetFS()))
		}

		server := NewServer()
		go server.Run()

		http.Handle("/push", websocket.Handler(func(ws *websocket.Conn) {
			client := NewClient(ws)
			server.AppendClient(client)
			client.Run() //synchronous
		}))

		nmutex := sync.Mutex{}
		http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				n := notifier.Notification{}
				data, err := ioutil.ReadAll(r.Body)

				if err == nil {
					log.Printf("Received: %v", string(data))
					err = json.Unmarshal(data, &n)
					if err == nil {
						n.RawTS = time.Now()
						n.Timestamp = n.RawTS.Format(time.RFC3339)

						nmutex.Lock()
						notifications = append(notifications, n)
						nmutex.Unlock()

						server.Send(n)
						//log.Println("main: returned from server.Send()")
					}
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "error occurred: %v\n", err)
				}
			}

			b, err := json.Marshal(notifications)
			if err != nil {
				//nop
			}
			w.Write(b)

			if storage != "" {
				ioutil.WriteFile(storage, b, 0x600)
			}
		})

		url := fmt.Sprintf("http://localhost:%d/", port)
		fmt.Fprintf(os.Stderr, "Access URL %s\n", url)
		browser.OpenURL(url)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

		return nil
	}

	app.Run(os.Args)
}
