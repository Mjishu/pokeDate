package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Notification struct {
	Id           uuid.UUID
	Actor        uuid.UUID
	Notifier     uuid.UUID
	Entity_text  string
	Entity_type  int
	Status       string //accepted, denied, unseen
	Date_created time.Time
	Date_seen    time.Time
}

func CreateNotification(pool *pgxpool.Pool, notification Notification) error {
	sql := `
		INSERT INTO notifications (actor, notifier, entity_text, entity_type, status) VALLUES ($1,$2,$3,$4,$5)
	`

	_, err := pool.Exec(context.TODO(), sql, notification.Actor, notification.Notifier, notification.Entity_text, notification.Entity_type, notification.Status)
	if err != nil {
		return err
	}
	return nil
}

func GetNotification(pool *pgxpool.Pool, idToGet uuid.UUID) ([]Notification, error) {
	sql := `SELECT id,actor,notifier,entity_text,entity_type,status,date_created,date_seen FROM notifications WHERE notifier = $1`

	rows, err := pool.Query(context.TODO(), sql, idToGet)
	if err != nil {
		return []Notification{}, err
	}

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		var date_seen pgtype.Date

		err := rows.Scan(
			&notification.Id,
			&notification.Actor,
			&notification.Notifier,
			&notification.Entity_text,
			&notification.Entity_type,
			&notification.Status,
			&notification.Date_created,
			&date_seen,
		)
		if err != nil {
			return []Notification{}, err
		}
		if date_seen.Status == pgtype.Present {
			notification.Date_seen = date_seen.Time
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
