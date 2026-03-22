package dto

import "time"

// ──── Skills ─────────────────
type SkillResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSkillRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// ──── Organizations ──────────
type CreateOrganizationRequest struct {
	DirectorID string `json:"director_id" validate:"required,uuid4"`
	Name       string `json:"name" validate:"required,min=1,max=255"`
	AboutUs    string `json:"about_us"`
}

type UpdateOrganizationRequest struct {
	Name    string `json:"name" validate:"omitempty,min=1,max=255"`
	AboutUs string `json:"about_us"`
}

type OrganizationResponse struct {
	ID         string    `json:"id"`
	DirectorID string    `json:"director_id"`
	Name       string    `json:"name"`
	AboutUs    string    `json:"about_us"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ──── Resumes ────────────────
type CreateResumeRequest struct {
	UserID      string   `json:"user_id" validate:"required,uuid4"`
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Description string   `json:"description"`
	AboutMe     string   `json:"about_me"`
	SkillNames  []string `json:"skill_names"`
}

type UpdateResumeRequest struct {
	Name        string   `json:"name" validate:"omitempty,min=1,max=255"`
	Description string   `json:"description"`
	AboutMe     string   `json:"about_me"`
	SkillNames  []string `json:"skill_names"`
}

type ResumeResponse struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	AboutMe     string          `json:"about_me"`
	Skills      []SkillResponse `json:"skills"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ──── Applications ───────────
type CreateApplicationRequest struct {
	ResumeID  string `json:"resume_id" validate:"required,uuid4"`
	VacancyID string `json:"vacancy_id" validate:"required,uuid4"`
}

type ApplicationResponse struct {
	ID          string    `json:"id"`
	ResumeID    string    `json:"resume_id"`
	VacancyID   string    `json:"vacancy_id"`
	Status      string    `json:"status"`
	ResumeName  string    `json:"resume_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
