package database

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conversation struct {
	Id                uuid.UUID
	Conversation_name string
	// Members           []Conversation_member
	// Messages          []Messages
}

type Messages struct {
	Id              uuid.UUID
	From_id         uuid.UUID //References users
	Conversation_id uuid.UUID
	Message_text    string
	Sent_datetime   time.Time
}

type Conversation_member struct {
	Member_id       uuid.UUID
	Conversation_id uuid.UUID // References Conversation
	Joined_datetime time.Time
	Left_datetime   time.Time
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
		err := rows.Scan(&member.Member_id, &member.Conversation_id, &member.Joined_datetime, &member.Left_datetime)
		if err != nil {
			return []Conversation_member{}, err
		}
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
		conversationSlice = append(conversationSlice, conversation)
	}
	fmt.Printf("the items are %v\n", conversationSlice)
	return conversationSlice, nil
}

// * fix the sql query to properly store the resposes
func GetConversation(pool *pgxpool.Pool, conversationId uuid.UUID) (Messages, error) {
	sql := `
		SELECT c.conversation_name, m.id, m.from_id, m.message_text, m.sent_datetime FROM conversation c LEFT JOIN
		messages m ON c.id  = m.conversation_id
		WHERE c.id = $1
	` // select all the stuff from messages, instead of *

	var messages Messages
	pool.QueryRow(context.TODO(), sql, conversationId).Scan() // scan into messages all stuff from sql query

	return messages, nil
}
