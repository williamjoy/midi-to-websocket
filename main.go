package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rakyll/portmidi"
)

var (
	deviceIDFlag      = flag.Int("deviceid", 3, "MIDI Device ID")
	httpListenAddress = flag.String("webSocketAddress", "0.0.0.0:8080", "WebSocket Listen Address")
)

func main() {
	flag.Parse()

	// midi
	deviceID := portmidi.DeviceID(*deviceIDFlag)
	deviceInfo := portmidi.Info(deviceID)
	if deviceInfo == nil {
		for i := 0; i < portmidi.CountDevices(); i++ {
			deviceInfo = portmidi.Info(portmidi.DeviceID(i))
			log.Printf("DeviceID <%d>, info is <%+v>\n", i, deviceInfo)
		}
		log.Fatal("Device not exists for id=", deviceID)
		os.Exit(-2)
	}

	log.Printf("DeviceID is <%d>, info is <%+v>\n", deviceID, deviceInfo)
	in, err := portmidi.NewInputStream(deviceID, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	midiEvents := in.Listen()

	clients := map[int64]*websocket.Conn{}

	go func() {
		for {
			select {
			case events := <-midiEvents:
				jsonPayload, err := json.Marshal(events)
				if err != nil {
					log.Fatal(err)
					continue
				}
				log.Printf("write midi: %+s\n", jsonPayload)
				for _, client := range clients {
					err := client.WriteMessage(websocket.TextMessage, jsonPayload)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}()

	// server
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	http.HandleFunc("/midi", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}

		uuid := time.Now().UnixNano()
		clients[uuid] = conn

		log.Printf("client connected: %d\n", uuid)

		defer func() {
			conn.Close()
			delete(clients, uuid)
			log.Printf("client disconnected: %d\n", uuid)
		}()

		for {
			// err is returned when the client is disconnected
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	})

	log.Fatal(http.ListenAndServe(*httpListenAddress, nil))
}
