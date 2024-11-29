package entities

import "time"

type InboxMessage struct {
	Id        int64     `gorm:"column:id;primaryKey"`
	MessageId string    `gorm:"column:message_id"`
	CreatedOn time.Time `gorm:"column:created_on"`
}
