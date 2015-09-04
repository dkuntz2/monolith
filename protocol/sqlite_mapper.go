package protocol

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteMapper struct {
	db         *sql.DB
	statements struct {
		createUser    *sql.Stmt
		getUser       *sql.Stmt
		getUsers      *sql.Stmt
		createMessage *sql.Stmt
		getMessage    *sql.Stmt
		getMessages   *sql.Stmt
	}
}

func NewSqliteMapper(dbName string) (*SqliteMapper, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        create table if not exists user (
            id integer primary key autoincrement,
            name string not null
        )
    `)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        create table if not exists message (
            id integer primary key autoincrement,
            author_id integer not null,
            payload text not null,
            foreign key(author_id) references user(id)
        )
    `)
	if err != nil {
		return nil, err
	}

	mapper := &SqliteMapper{}
	mapper.db = db

	mapper.statements.createUser, err = db.Prepare("insert into user (name) values (?)")
	if err != nil {
		return nil, err
	}

	mapper.statements.getUser, err = db.Prepare("select id, name from user where id = ?")
	if err != nil {
		return nil, err
	}

	mapper.statements.getUsers, err = db.Prepare("select id, name from user")
	if err != nil {
		return nil, err
	}

	mapper.statements.createMessage, err = db.Prepare("insert into message (author_id, payload) values (?, ?)")
	if err != nil {
		return nil, err
	}

	mapper.statements.getMessage, err = db.Prepare("select id, author_id, payload from message where id = ?")
	if err != nil {
		return nil, err
	}

	mapper.statements.getMessages, err = db.Prepare("select id, author_id, payload from message")
	if err != nil {
		return nil, err
	}

	return mapper, nil
}

func (mapper *SqliteMapper) SaveUser(user *User) (int64, error) {
	result, err := mapper.statements.createUser.Exec(user.Name)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (mapper *SqliteMapper) GetUser(id int64) (*User, error) {
	row := mapper.statements.getUser.QueryRow(id)

	user := User{}
	err := row.Scan(&user.Id, &user.Name)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (mapper *SqliteMapper) GetUsers() ([]*User, error) {
	rows, err := mapper.statements.getUsers.Query()
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0, 5)
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Name)
		users = append(users, &user)
	}

	return users, nil
}

func (mapper *SqliteMapper) SaveMessage(message *Message) (int64, error) {
	result, err := mapper.statements.createMessage.Exec(message.AuthorId, message.Payload)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}

func (mapper *SqliteMapper) GetMessage(id int64) (*Message, error) {
	row := mapper.statements.getMessage.QueryRow(id)

	message := Message{}
	err := row.Scan(&message.Id, &message.AuthorId, &message.Payload)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (mapper *SqliteMapper) GetMessages() ([]*Message, error) {
	rows, err := mapper.statements.getMessages.Query()
	if err != nil {
		return nil, err
	}

	messages := make([]*Message, 0, 5)
	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.Id, &message.AuthorId, &message.Payload)
		messages = append(messages, &message)
	}

	return messages, nil
}
