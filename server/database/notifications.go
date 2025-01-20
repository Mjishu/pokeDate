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
	Notifier     uuid.UUID `json:"Notifier_id"`
	Animal_id    uuid.UUID `json:"Animal_id"`
	Entity_text  string
	Entity_type  int
	Status       string //accepted, denied, unseen
	Date_created time.Time
	Date_seen    time.Time
}

func CreateNotification(pool *pgxpool.Pool, notification Notification) error {
	//* get or add from animal_groups, add the notification_id thats created here to it with the animal id provided
	sql := `
		INSERT INTO notifications (actor, notifier, entity_text, entity_type, status) VALUES ($1,$2,$3,$4,$5) RETURNING id
	`

	var notification_id uuid.UUID
	err := pool.QueryRow(context.TODO(), sql, notification.Actor, notification.Notifier, notification.Entity_text, notification.Entity_type, notification.Status).Scan(&notification_id)
	if err != nil {
		return err
	}

	animalGroupsSql := `INSERT INTO animal_groups (animal_id, notification_id) VALUES ($1, $2)`

	_, err = pool.Exec(context.TODO(), animalGroupsSql, notification.Animal_id, notification_id)
	if err != nil {
		return err

	}
	return nil
}

func GetNotification(pool *pgxpool.Pool, idToGet uuid.UUID) ([]Notification, error) {
	//* get or add from animal_groups
	sql := `SELECT n.id,n.actor,n.notifier,n.entity_text,n.entity_type,n.status,n.date_created,n.date_seen, ag.animal_id FROM notifications n LEFT JOIN animal_groups ag ON n.id = ag.notification_id WHERE n.notifier = $1`

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
			&notification.Animal_id,
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
