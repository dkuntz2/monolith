package protocol

import (
	"time"
)

type ProtocolMessage struct {
	Id   int64     `json:"id"`
	Type string    `json:"type"`
	Text string    `json:"text"`
	Date time.Time `json:"time"`
}

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Id       int64  `json:"id"`
	AuthorId int64  `json:"author_id"`
	Payload  string `json:"payload"`
}

type DataMapper interface {
	SaveMessage(*Message) (int64, error)
	GetMessage(int64) (*Message, error)
	GetMessages() ([]*Message, error)

	SaveUser(*User) (int64, error)
	GetUser(int64) (*User, error)
	GetUsers() ([]*User, error)
}
