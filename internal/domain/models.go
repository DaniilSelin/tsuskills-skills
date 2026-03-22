package domain

import (
	"time"

	"github.com/google/uuid"
)

type Skill struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

type Resume struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	Description string
	AboutMe     string
	Skills      []Skill
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Organization struct {
	ID         uuid.UUID
	DirectorID uuid.UUID
	Name       string
	AboutUs    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ApplicationStatus описывает состояние отклика
type ApplicationStatus string

const (
	AppStatusPending  ApplicationStatus = "pending"
	AppStatusAccepted ApplicationStatus = "accepted"
	AppStatusRejected ApplicationStatus = "rejected"
)

type Application struct {
	ID        uuid.UUID
	ResumeID  uuid.UUID
	VacancyID uuid.UUID
	Status    ApplicationStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	// подгружаемые при чтении
	ResumeName  string // join-поле
	VacancyTitle string // join-поле
}
