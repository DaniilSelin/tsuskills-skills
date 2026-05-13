package kafka

import (
	"encoding/json"
	"time"
)

const (
	EntityVacancy = "vacancy"
	EntityUser    = "user"

	EventVacancyCreated = "vacancy.created"
	EventVacancyUpdated = "vacancy.updated"
	EventVacancyDeleted = "vacancy.deleted"
	EventUserCreated    = "user.created"
	EventUserUpdated    = "user.updated"
	EventUserDeleted    = "user.deleted"
)

type Event struct {
	Type      string          `json:"type"`
	Entity    string          `json:"entity"`
	EntityID  string          `json:"entity_id"`
	Source    string          `json:"source,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type ConsumerConfig struct {
	Brokers     []string      `mapstructure:"Brokers"`
	Topic       string        `mapstructure:"Topic"`
	GroupID     string        `mapstructure:"GroupID"`
	DialTimeout time.Duration `mapstructure:"DialTimeout"`
}

func NewEvent(eventType, entity, entityID string, payload interface{}) (Event, error) {
	var raw json.RawMessage
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return Event{}, err
		}
		raw = json.RawMessage(b)
	}

	return Event{
		Type:      eventType,
		Entity:    entity,
		EntityID:  entityID,
		Timestamp: time.Now().UTC(),
		Payload:   raw,
	}, nil
}
