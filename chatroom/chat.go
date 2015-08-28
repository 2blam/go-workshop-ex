package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/net/websocket"
)

type IncomingMessage struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type OutgoingMessage struct {
	Event string    `json:"event"`
	Data  EventData `json:"data,omitempty"`
}

type EventData struct {
	Username  string `json:"username,omitempty"`
	UserCount int    `json:"numUsers,omitempty"`
	Message   string `json:"message,omitempty"`
}

var chatroom = Chatroom{}

func ChatServer(ws *websocket.Conn) {
	chatroom.Add(ws)

	var username string
	for {
		var msgin IncomingMessage
		if err := websocket.JSON.Receive(ws, &msgin); err != nil {
			log.Printf("Error receving message: %v", err)
			break
		}

		log.Printf("Received: %#v", msgin)

		//checking
		log.Printf("event: %#v", msgin.Event)
		log.Printf("data: %#v", msgin.Data)
		if msgin.Event == "add user" {
			//emit a login event with count of participants
			loginMsg := OutgoingMessage{
				Event: "login",
				Data: EventData{
					UserCount: chatroom.Count(),
				},
			}
			if err := websocket.JSON.Send(ws, loginMsg); err != nil {
				log.Printf("Failed to send message: %v", err)
			}

			//update username
			username = msgin.Data

			//boardcast a user joined event
			joinMsg := OutgoingMessage{
				Event: "user joined",
				Data: EventData{
					Username:  username,
					UserCount: chatroom.Count(),
				},
			}
			if err := chatroom.Broadcast(joinMsg, ws); err != nil {
				log.Printf("Failed to boardcast join message: %v", err)
			}

		} else if msgin.Event == "new message" {
			//boardcast a new message
			newMsg := OutgoingMessage{
				Event: "new message",
				Data: EventData{
					Username: username,
					Message:  msgin.Data,
				},
			}
			if err := chatroom.Broadcast(newMsg, ws); err != nil {
				log.Printf("Failed to boardcast new message: %v", err)
			}
		}

	}

	if username != "" {
		leaveMsg := OutgoingMessage{
			Event: "user left",
			Data: EventData{
				Username:  username,
				UserCount: chatroom.Count() - 1,
			},
		}
		if err := chatroom.Broadcast(leaveMsg, ws); err != nil {
			log.Printf("Failed to leave room: %v", err)
		}
	}

	chatroom.Remove(ws)
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))

	chatServer := websocket.Server{
		Handler: ChatServer,
	}
	http.Handle("/chat", chatServer)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	host := ":" + port
	log.Printf("Listening on %s", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// Chatroom maintain a list of websocket connections
type Chatroom struct {
	websockets []*websocket.Conn
	sync.Mutex
}

// Count returns the number of participanting conn in a Chatroom
func (rm *Chatroom) Count() int {
	return len(rm.websockets)
}

// Add a Conn to list of participants
func (rm *Chatroom) Add(ws *websocket.Conn) {
	rm.Lock()
	defer rm.Unlock()

	rm.websockets = append(rm.websockets, ws)
}

// Remove a Conn from list of participants
func (rm *Chatroom) Remove(ws *websocket.Conn) {
	rm.Lock()
	defer rm.Unlock()

	for i, conn := range rm.websockets {
		if conn == ws {
			rm.websockets = append(rm.websockets[:i], rm.websockets[i+1:]...)
			break
		}
	}
}

// Broadcast a message to all participants in Chatroom except sender
func (rm *Chatroom) Broadcast(i interface{}, sender *websocket.Conn) error {
	rm.Lock()
	defer rm.Unlock()

	log.Printf("Broadcasting: %v", i)

	var lasterr error
	for _, ws := range rm.websockets {
		if ws != sender {
			if err := websocket.JSON.Send(ws, i); err != nil {
				log.Printf("send: %v", err)
				lasterr = err
			}
		}
	}

	// returning only the last error encountered
	// in practice we should return a custom error that has all errors collected
	return lasterr
}
