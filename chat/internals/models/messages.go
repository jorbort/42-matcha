package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	DB *pgxpool.Pool
}

type Conversation struct {
	ID        int       `json:"id"`
	User1ID   string    `json:"user1_id"`
	User2ID   string    `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID             int       `json:"id"`
	ConversationID int       `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	ReceiverID     string    `json:"receiver_id"`
	Message        string    `json:"message"`
	CreatedAt      time.Time `json:"created_at"`
	IsRead         bool      `json:"is_read"`
}

func (m *Models) GetOrCreateConversation(user1ID, user2ID string) (int, error) {
	query := `
	SELECT id FROM conversations
	WHERE (user1_id = $1 AND user2_id = $2)
	OR  (user1_id = $2 AND user2_id = $1)`
	var conversationID int
	err := m.DB.QueryRow(context.Background(), query, user1ID, user2ID).Scan(&conversationID)
	if err == nil {
		return conversationID, nil
	}
	query = `INSERT INTO conversations (user1_id, user2_id)
	VALUES ($1, $2)
	RETURNING id`
	err = m.DB.QueryRow(context.Background(), query, user1ID, user2ID).Scan(&conversationID)
	if err != nil {
		return 0, err
	}
	return conversationID, nil
}

func (m *Models) SaveMessage(msg *Message) error {
	query := `INSERT INTO messages (conversation_id
	sender_id, receiver_id, message)
	VALUES($1, $2, $3, $4)
	RETURNING id , created_at`
	return m.DB.QueryRow(context.Background(), query,
		msg.ConversationID, msg.SenderID, msg.ReceiverID, msg.Message).Scan(&msg.ID, &msg.CreatedAt)
}

func (m *Models) GetConversationMsgs(conversationID int) ([]Message, error) {
	query := `
	SELECT id, conversation_id,sender_id, receiver_id, message, created_at
	FROM messages
	WHERE conversation_id = $1
	ORDER BY created_at ASC`
	rows, err := m.DB.Query(context.Background(), query, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var msgs []Message
	for rows.Next() {
		var msg Message
		err := rows.Scan(
			&msg.ID,
			&msg.ConversationID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Message,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

func (m *Models) GetUserConversations(userID string) ([]Conversation, error) {
	query := `
        SELECT id, user1_id, user2_id, created_at
        FROM conversations
        WHERE user1_id = $1 OR user2_id = $1`

	rows, err := m.DB.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conv Conversation
		err := rows.Scan(
			&conv.ID,
			&conv.User1ID,
			&conv.User2ID,
			&conv.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}
