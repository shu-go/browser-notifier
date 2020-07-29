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

	"github.com/shu-go/browser-notifier"
	"github.com/shu-go/gli"
)

type globalCmd struct {
	Port    uint16 `cli:"port, p=PORT"  default:"7878"  help:"server port"`
	Storage string `cli:"storage, s=FILE" help:"file name to store contents"`
}

func main() {
	global := globalCmd{}
	app := gli.NewWith(&global)
	app.Name = "Browser Notifier"
	app.Usage = "A notifier using Web Notifications."
	app.Version = "0.1.0"
	_, _, err := app.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "app error: %v\n", err)
		os.Exit(1)
	}

	notifications := []notifier.Notification{}
	if global.Storage != "" {
		if f, err := os.Open(global.Storage); err == nil {
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

		if global.Storage != "" {
			ioutil.WriteFile(global.Storage, b, 0x600)
		}
	})

	url := fmt.Sprintf("http://localhost:%d/", global.Port)
	fmt.Fprintf(os.Stderr, "Access URL %s\n", url)
	browser.OpenURL(url)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", global.Port), nil))
}
