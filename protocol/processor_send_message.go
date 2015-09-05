package protocol

import (
	"strconv"
)

func SendMessageProcessor(hub *Hub, request *ProtocolMessage) (*ProtocolMessage, error) {
	id, err := hub.Mapper.SaveMessage(&Message{
		AuthorId: request.Id,
		Payload:  request.Text,
	})
	if err != nil {
		return nil, err
	}

	hub.GlobalBroadcast(&ProtocolMessage{
		Type: "message_broadcast",
		Id:   request.Id,
		Text: request.Text,
	})

	return &ProtocolMessage{
		Type: "new_message",
		Text: strconv.Itoa(int(id)),
	}, nil
}
