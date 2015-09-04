package protocol

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"time"
)

type Hub struct {
	connections      map[*websocket.Conn]bool
	globalBroadcasts chan *ProtocolMessage
	mapper           DataMapper
}

func NewHub() *Hub {
	mapper, err := NewSqliteMapper("test.db")
	if err != nil {
		log.Fatal(err)
	}
	return &Hub{
		connections:      make(map[*websocket.Conn]bool),
		globalBroadcasts: make(chan *ProtocolMessage),
		mapper:           mapper,
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case message := <-hub.globalBroadcasts:
			for connection := range hub.connections {
				go func(connection *websocket.Conn, message *ProtocolMessage) {
					err := connection.WriteJSON(*message)
					if err != nil {
						log.Println(err)
					}
				}(connection, message)
			}
		}
	}
}

func (hub *Hub) Attach(conn *websocket.Conn) {
	hub.connections[conn] = true
	go hub.listen(conn)
}

func (hub *Hub) listen(conn *websocket.Conn) {
	for {
		protoMessage := &ProtocolMessage{}
		err := conn.ReadJSON(protoMessage)
		if err != nil {
			conn.Close()
			log.Println(err)
			return
		}

		response := ProtocolMessage{}
		switch protoMessage.Type {
		case "hello":
			id, err := hub.mapper.SaveUser(&User{Name: protoMessage.Text})
			if err != nil {
				response.Type = "error"
				response.Text = err.Error()
			} else {
				response.Type = "new_user"
				response.Text = strconv.Itoa(int(id))
			}
		case "get_messages":
			messages, err := hub.mapper.GetMessages()
			if err != nil {
				response.Type = "error"
				response.Text = err.Error()
			} else {
				messagesString, err := json.Marshal(messages)
				if err != nil {
					response.Type = "error"
					response.Text = err.Error()
				} else {
					response.Type = "messages"
					response.Text = string(messagesString)
				}
			}
		case "send_message":
			id, err := hub.mapper.SaveMessage(&Message{AuthorId: protoMessage.Id, Payload: protoMessage.Text})
			if err != nil {
				response.Type = "error"
				response.Text = err.Error()
			} else {
				response.Type = "new_message"
				response.Text = strconv.Itoa(int(id))

				hub.globalBroadcasts <- &ProtocolMessage{
					Type: "message_broadcast",
					Id:   protoMessage.Id,
					Text: protoMessage.Text,
				}
			}
		}

		response.Date = time.Now()

		err = conn.WriteJSON(response)
		if err != nil {
			conn.Close()
			log.Println(err)
			return
		}
	}
	conn.Close()
}
