package protocol

import (
	"encoding/json"
)

func GetMessagesProcessor(hub *Hub, request *ProtocolMessage) (*ProtocolMessage, error) {
	messages, err := hub.Mapper.GetMessages()
	if err != nil {
		return nil, err
	}

	messagesString, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	return &ProtocolMessage{
		Type: "messages",
		Text: string(messagesString),
	}, nil
}
