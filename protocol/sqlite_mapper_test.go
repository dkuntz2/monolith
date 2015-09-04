package protocol

import (
	"os"
	"testing"
)

func TestSqliteGetMessages(t *testing.T) {
	mapper, err := NewSqliteMapper("testing.db")
	if err != nil {
		t.Error(err)
	}

	messages := []*Message{
		&Message{AuthorId: 1, Payload: "Hello!"},
		&Message{AuthorId: 2, Payload: "Howdy!"},
		&Message{AuthorId: 1, Payload: "What's up?"},
		&Message{AuthorId: 2, Payload: "testing..."},
		&Message{AuthorId: 3, Payload: "did this by chance insert?"},
		&Message{AuthorId: 4, Payload: "well, actually, we know insertions work"},
		&Message{AuthorId: 1, Payload: "but retrieval is having problems"},
	}

	for _, message := range messages {
		mapper.SaveMessage(message)
	}

	dbMessages, err := mapper.GetMessages()
	if err != nil {
		t.Error(err)
	}

	for _, oMsg := range messages {
		// find in inserted messages
		found := false
		for _, dMsg := range dbMessages {
			if oMsg.AuthorId == dMsg.AuthorId && oMsg.Payload == dMsg.Payload {
				found = true
				break
			}
		}

		if !found {
			t.Error("Message(", oMsg.Payload, ") not found in DB response")
		}
	}

	os.Remove("testing.db")
}
