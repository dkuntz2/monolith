package protocol

import (
	"errors"
	"github.com/gorilla/websocket"
	"log"
)

type Hub struct {
	connections      map[*websocket.Conn]bool
	globalBroadcasts chan *ProtocolMessage
	Mapper           DataMapper
	processors       map[string]HubProcessor
}

type HubProcessor func(*Hub, *ProtocolMessage) (*ProtocolMessage, error)

func NewHub() *Hub {
	mapper, err := NewSqliteMapper("test.db")
	if err != nil {
		log.Fatal(err)
	}
	hub := &Hub{
		connections:      make(map[*websocket.Conn]bool),
		globalBroadcasts: make(chan *ProtocolMessage),
		Mapper:           mapper,
		processors:       make(map[string]HubProcessor),
	}

	hub.RegisterProcessor("hello", HelloProcessor)
	hub.RegisterProcessor("get_messages", GetMessagesProcessor)
	hub.RegisterProcessor("send_message", SendMessageProcessor)

	return hub
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

func (hub *Hub) RegisterProcessor(message_name string, fn HubProcessor) error {
	_, exists := hub.processors[message_name]
	if exists {
		return errors.New("Processor already exists for message_name")
	}

	hub.processors[message_name] = fn
	return nil
}

func (hub *Hub) GlobalBroadcast(message *ProtocolMessage) {
	hub.globalBroadcasts <- message
}

func (hub *Hub) Attach(conn *websocket.Conn) {
	hub.connections[conn] = true
	go hub.listen(conn)
}

func (hub *Hub) listen(conn *websocket.Conn) {
	for {
		request := &ProtocolMessage{}
		err := conn.ReadJSON(request)
		if err != nil {
			conn.Close()
			log.Println(err)
			return
		}

		fn, exists := hub.processors[request.Type]
		var response *ProtocolMessage
		if exists {
			response, err = fn(hub, request)
			if err != nil {
				response = &ProtocolMessage{
					Type: "error",
					Text: err.Error(),
				}
			}
		} else {
			response = &ProtocolMessage{
				Type: "error",
				Text: "unknown method",
			}
		}

		err = conn.WriteJSON(response)
		if err != nil {
			conn.Close()
			log.Println(err)
			return
		}
	}
	conn.Close()
}
