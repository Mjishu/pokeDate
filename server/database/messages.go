package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conversation struct {
	Id                uuid.UUID
	Conversation_name string
	Members           []Conversation_member
	Messages          []Messages
}

type Messages struct {
	Id              uuid.UUID
	From_id         uuid.UUID //References users
	Conversation_id uuid.UUID
	Message_text    string `json:"message"`
	Sent_datetime   time.Time
	From_user       User
}

type Conversation_member struct {
	Member_id       uuid.UUID
	Conversation_id uuid.UUID // References Conversation
	Joined_datetime time.Time
	Left_datetime   time.Time
	User            User
}

/*
	have conversation, conversation_member, messages

	select conversation_id from conversation_member where member_id = userId

	select messages where conversation_id = ^conversationId
*/

/*
SELECT c.id,c.conversation_name, cm.member_id, m.id, m.from_id, m.message_text, m.sent_datetime FROM conversation_member cm LEFT JOIN
			conversation c ON cm.conversation_id = c.id
			LEFT JOIN messages m ON c.id = m.conversation_id
			WHERE cm.member_id = $1
*/

func GetMessageUsers(pool *pgxpool.Pool, conversationId uuid.UUID) ([]Conversation_member, error) {
	sql := `
		SELECT member_id, conversation_id, joined_datetime, left_datetime FROM conversation_member WHERE conversation_id = $1
	`

	var members []Conversation_member

	rows, err := pool.Query(context.TODO(), sql, conversationId)
	if err != nil {
		return []Conversation_member{}, err
	}
	for rows.Next() {
		var member Conversation_member
		var left_date pgtype.Date

		err := rows.Scan(&member.Member_id, &member.Conversation_id, &member.Joined_datetime, &left_date)
		if err != nil {
			return []Conversation_member{}, err
		}
		if left_date.Status == pgtype.Present {
			member.Left_datetime = left_date.Time
		} else {
			member.Left_datetime = time.Time{} //! sets time to time.Nil , make sure to properly check for this later
		}

		user, err := GetUserById(pool, member.Member_id)
		if err != nil {
			return []Conversation_member{}, err
		}
		member.User = user

		members = append(members, member)
	}
	return members, nil
}

func GetMessages(pool *pgxpool.Pool, conversationId uuid.UUID) ([]Messages, error) {
	sql := `
		SELECT id,from_id,message_text,sent_datetime from messages where conversation_id = $1
	`

	var messageSlice []Messages
	rows, err := pool.Query(context.TODO(), sql, conversationId)
	if err != nil {
		return []Messages{}, err
	}

	for rows.Next() {
		var message Messages
		err := rows.Scan(&message.Id, &message.From_id, &message.Message_text, &message.Sent_datetime)
		if err != nil {
			return []Messages{}, err
		}
		message.Conversation_id = conversationId

		user, err := GetUserById(pool, message.From_id)
		if err != nil {
			return []Messages{}, err
		}
		message.From_user = user

		messageSlice = append(messageSlice, message)
	}
	return messageSlice, nil
}

// * how to iterate over different messages in a conversation
func GetConversations(pool *pgxpool.Pool, userId uuid.UUID) ([]Conversation, error) {
	fmt.Println("i was called")
	// make sure that this gets more than 1 conversation
	sql := `
		SELECT c.id, c.conversation_name  from conversation_member cm LEFT JOIN conversation c on cm.conversation_id = c.id WHERE cm.member_id = $1
	`

	var conversationSlice []Conversation
	rows, err := pool.Query(context.TODO(), sql, userId)
	if err != nil {
		return []Conversation{}, err
	}

	for rows.Next() {
		var conversation Conversation
		err := rows.Scan(&conversation.Id, &conversation.Conversation_name)
		if err != nil {
			return []Conversation{}, err
		}
		users, err := GetMessageUsers(pool, conversation.Id)
		if err != nil {
			fmt.Println("2")
			return []Conversation{}, err
		}

		messages, err := GetMessages(pool, conversation.Id)
		if err != nil {
			return []Conversation{}, err
		}

		conversation.Members = users
		conversation.Messages = messages

		conversationSlice = append(conversationSlice, conversation)
	}

	return conversationSlice, nil
}

// * fix the sql query to properly store the resposes
func GetConversation(pool *pgxpool.Pool, conversationId uuid.UUID) (Conversation, error) {
	fmt.Printf("conversation id is %v\n", conversationId)
	sql := `
		SELECT conversation_name, id FROM conversation WHERE id = $1
	` // select all the stuff from messages, instead of *
	// this should call GetMessages and getMessageUsers

	var conversation Conversation
	pool.QueryRow(context.TODO(), sql, conversationId).Scan(&conversation.Conversation_name, &conversation.Id) // scan into messages all stuff from sql query

	users, err := GetMessageUsers(pool, conversation.Id)
	if err != nil {
		fmt.Println("2")
		return Conversation{}, err
	}

	messages, err := GetMessages(pool, conversation.Id)
	if err != nil {
		return Conversation{}, err
	}

	conversation.Members = users
	conversation.Messages = messages

	fmt.Printf("message is %v\n", messages)

	return conversation, nil
}

func CreateConversation(pool *pgxpool.Pool, userId uuid.UUID, recipientId uuid.UUID) (uuid.UUID, error) {
	conversationSql := `INSERT INTO conversation (conversation_name) VALUES ('New Conversation') RETURNING id`
	memberSql := `INSERT INTO conversation_member (member_id, conversation_id) VALUES ($1, $2)`

	var conversationId uuid.UUID
	pool.QueryRow(context.TODO(), conversationSql).Scan(&conversationId)
	if (conversationId == uuid.UUID{}) {
		return uuid.UUID{}, errors.New("conversation id is nil")
	}

	_, err := pool.Exec(context.TODO(), memberSql, userId, conversationId)
	if err != nil {
		return uuid.UUID{}, nil
	}

	_, err = pool.Exec(context.TODO(), memberSql, recipientId, conversationId)
	if err != nil {
		return uuid.UUID{}, nil
	}

	return conversationId, nil

}

func CreateMessage(pool *pgxpool.Pool, userId uuid.UUID, messageData Messages) error {
	exists, err := UserMessageExists(pool, userId)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("user")
	}

	sql := `INSERT INTO messages (from_id, message_text, conversation_id) VALUES ($1, $2, $3) `

	_, err = pool.Exec(context.TODO(), sql, userId, messageData.Message_text, messageData.Conversation_id)
	if err != nil {
		return err
	}

	return nil
}

func UserMessageExists(pool *pgxpool.Pool, userId uuid.UUID) (bool, error) {
	sql := `SELECT c.id FROM conversation_member cm LEFT JOIN conversation c ON cm.conversation_id = c.id WHERE cm.member_id = $1`

	var id uuid.UUID
	err := pool.QueryRow(context.TODO(), sql, userId).Scan(&id)
	if err != nil {
		return false, err
	}

	return true, nil
}
