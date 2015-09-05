package protocol

import (
	"strconv"
)

func HelloProcessor(hub *Hub, request *ProtocolMessage) (*ProtocolMessage, error) {
	id, err := hub.Mapper.SaveUser(&User{Name: request.Text})
	if err != nil {
		return nil, err
	}

	return &ProtocolMessage{
		Type: "new_user",
		Text: strconv.Itoa(int(id)),
	}, nil
}
