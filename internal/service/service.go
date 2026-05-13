package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"tsuskills-skills/internal/domain"
	"tsuskills-skills/internal/infra/kafka"
	"tsuskills-skills/internal/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IRepository interface {
	// Skills
	SearchSkills(ctx context.Context, query string) ([]domain.Skill, error)
	CreateSkill(ctx context.Context, name string) (*domain.Skill, error)
	DeleteSkill(ctx context.Context, id int) error
	// Organizations
	CreateOrganization(ctx context.Context, o *domain.Organization) (uuid.UUID, error)
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error)
	GetOrganizationByDirector(ctx context.Context, directorID uuid.UUID) (*domain.Organization, error)
	UpdateOrganization(ctx context.Context, o *domain.Organization) error
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	// Resumes
	CreateResume(ctx context.Context, res *domain.Resume) (uuid.UUID, error)
	GetResumeByID(ctx context.Context, id uuid.UUID) (*domain.Resume, error)
	ListResumesByUser(ctx context.Context, userID uuid.UUID) ([]domain.Resume, error)
	UpdateResume(ctx context.Context, res *domain.Resume) error
	DeleteResume(ctx context.Context, id uuid.UUID) error
	DeleteResumesByUser(ctx context.Context, userID uuid.UUID) error
	// Applications
	CreateApplication(ctx context.Context, app *domain.Application) (uuid.UUID, error)
	GetApplicationsByVacancy(ctx context.Context, vacancyID uuid.UUID) ([]domain.Application, error)
	GetApplicationsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Application, error)
	GetApplicationByID(ctx context.Context, id uuid.UUID) (*domain.Application, error)
}

type Service struct {
	repo IRepository
	log  logger.Logger
}

func New(repo IRepository, log logger.Logger) *Service {
	return &Service{repo: repo, log: log}
}

func (s *Service) errCode(err error) domain.ErrorCode {
	if errors.Is(err, domain.ErrNotFound) {
		return domain.CodeNotFound
	}
	if errors.Is(err, domain.ErrConflict) {
		return domain.CodeConflict
	}
	return domain.CodeInternal
}

// ──── Skills ─────────────────────────────

func (s *Service) SearchSkills(ctx context.Context, query string) ([]domain.Skill, domain.ErrorCode) {
	list, err := s.repo.SearchSkills(ctx, query)
	if err != nil {
		s.log.Error(ctx, "SearchSkills", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return list, domain.CodeOK
}

func (s *Service) CreateSkill(ctx context.Context, name string) (*domain.Skill, domain.ErrorCode) {
	sk, err := s.repo.CreateSkill(ctx, name)
	if err != nil {
		s.log.Error(ctx, "CreateSkill", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return sk, domain.CodeOK
}

func (s *Service) DeleteSkill(ctx context.Context, id int) domain.ErrorCode {
	if err := s.repo.DeleteSkill(ctx, id); err != nil {
		return s.errCode(err)
	}
	return domain.CodeOK
}

// ──── Organizations ──────────────────────

func (s *Service) CreateOrganization(ctx context.Context, o *domain.Organization) (uuid.UUID, domain.ErrorCode) {
	now := time.Now()
	o.ID = uuid.New()
	o.CreatedAt = now
	o.UpdatedAt = now

	id, err := s.repo.CreateOrganization(ctx, o)
	if err != nil {
		s.log.Error(ctx, "CreateOrganization", zap.Error(err))
		return uuid.Nil, domain.CodeInternal
	}
	return id, domain.CodeOK
}

func (s *Service) GetOrganization(ctx context.Context, id uuid.UUID) (*domain.Organization, domain.ErrorCode) {
	o, err := s.repo.GetOrganizationByID(ctx, id)
	if err != nil {
		return nil, s.errCode(err)
	}
	return o, domain.CodeOK
}

func (s *Service) GetMyOrganization(ctx context.Context, directorID uuid.UUID) (*domain.Organization, domain.ErrorCode) {
	o, err := s.repo.GetOrganizationByDirector(ctx, directorID)
	if err != nil {
		return nil, s.errCode(err)
	}
	return o, domain.CodeOK
}

func (s *Service) UpdateOrganization(ctx context.Context, o *domain.Organization) domain.ErrorCode {
	o.UpdatedAt = time.Now()
	if err := s.repo.UpdateOrganization(ctx, o); err != nil {
		return s.errCode(err)
	}
	return domain.CodeOK
}

func (s *Service) DeleteOrganization(ctx context.Context, id uuid.UUID) domain.ErrorCode {
	if err := s.repo.DeleteOrganization(ctx, id); err != nil {
		return s.errCode(err)
	}
	return domain.CodeOK
}

// ──── Resumes ────────────────────────────

func (s *Service) CreateResume(ctx context.Context, res *domain.Resume) (uuid.UUID, domain.ErrorCode) {
	now := time.Now()
	res.ID = uuid.New()
	res.CreatedAt = now
	res.UpdatedAt = now

	id, err := s.repo.CreateResume(ctx, res)
	if err != nil {
		s.log.Error(ctx, "CreateResume", zap.Error(err))
		return uuid.Nil, domain.CodeInternal
	}
	return id, domain.CodeOK
}

func (s *Service) GetResume(ctx context.Context, id uuid.UUID) (*domain.Resume, domain.ErrorCode) {
	r, err := s.repo.GetResumeByID(ctx, id)
	if err != nil {
		return nil, s.errCode(err)
	}
	return r, domain.CodeOK
}

func (s *Service) ListMyResumes(ctx context.Context, userID uuid.UUID) ([]domain.Resume, domain.ErrorCode) {
	list, err := s.repo.ListResumesByUser(ctx, userID)
	if err != nil {
		s.log.Error(ctx, "ListMyResumes", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return list, domain.CodeOK
}

func (s *Service) UpdateResume(ctx context.Context, res *domain.Resume) domain.ErrorCode {
	res.UpdatedAt = time.Now()
	if err := s.repo.UpdateResume(ctx, res); err != nil {
		return s.errCode(err)
	}
	return domain.CodeOK
}

func (s *Service) DeleteResume(ctx context.Context, id uuid.UUID) domain.ErrorCode {
	if err := s.repo.DeleteResume(ctx, id); err != nil {
		return s.errCode(err)
	}
	return domain.CodeOK
}

// ──── Applications ───────────────────────

func (s *Service) CreateApplication(ctx context.Context, app *domain.Application) (uuid.UUID, domain.ErrorCode) {
	now := time.Now()
	app.ID = uuid.New()
	app.Status = domain.AppStatusPending
	app.CreatedAt = now
	app.UpdatedAt = now

	id, err := s.repo.CreateApplication(ctx, app)
	if err != nil {
		s.log.Error(ctx, "CreateApplication", zap.Error(err))
		return uuid.Nil, domain.CodeInternal
	}
	return id, domain.CodeOK
}

func (s *Service) GetApplicationsByVacancy(ctx context.Context, vacancyID uuid.UUID) ([]domain.Application, domain.ErrorCode) {
	list, err := s.repo.GetApplicationsByVacancy(ctx, vacancyID)
	if err != nil {
		s.log.Error(ctx, "GetApplicationsByVacancy", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return list, domain.CodeOK
}

func (s *Service) GetMyApplications(ctx context.Context, userID uuid.UUID) ([]domain.Application, domain.ErrorCode) {
	list, err := s.repo.GetApplicationsByUser(ctx, userID)
	if err != nil {
		s.log.Error(ctx, "GetMyApplications", zap.Error(err))
		return nil, domain.CodeInternal
	}
	return list, domain.CodeOK
}

func (s *Service) GetApplication(ctx context.Context, id uuid.UUID) (*domain.Application, domain.ErrorCode) {
	a, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, s.errCode(err)
	}
	return a, domain.CodeOK
}

func (s *Service) HandleEvent(ctx context.Context, event kafka.Event) error {
	switch event.Entity {
	case kafka.EntityVacancy:
		return s.handleVacancyEvent(ctx, event)
	case kafka.EntityUser:
		return s.handleUserEvent(ctx, event)
	default:
		return nil
	}
}

func (s *Service) handleVacancyEvent(ctx context.Context, event kafka.Event) error {
	switch event.Type {
	case kafka.EventVacancyDeleted:
		vacancyID, err := uuid.Parse(event.EntityID)
		if err != nil || vacancyID == uuid.Nil {
			var payload struct {
				ID string `json:"id"`
			}
			if err2 := json.Unmarshal(event.Payload, &payload); err2 != nil {
				s.log.Error(ctx, "handleVacancyEvent: invalid payload", zap.Error(err2))
				return nil
			}
			vacancyID, err = uuid.Parse(payload.ID)
			if err != nil {
				s.log.Error(ctx, "handleVacancyEvent: invalid vacancy id", zap.Error(err))
				return nil
			}
		}
		if err := s.repo.DeleteApplicationsByVacancy(ctx, vacancyID); err != nil {
			s.log.Error(ctx, "handleVacancyEvent: delete applications", zap.Error(err))
			return err
		}
		s.log.Info(ctx, "handleVacancyEvent: deleted related applications", zap.String("vacancy_id", vacancyID.String()))
	default:
	}
	return nil
}

func (s *Service) handleUserEvent(ctx context.Context, event kafka.Event) error {
	switch event.Type {
	case kafka.EventUserCreated:
		s.log.Info(ctx, "handleUserEvent: user created event", zap.String("user_id", event.EntityID))
	case kafka.EventUserUpdated:
		s.log.Info(ctx, "handleUserEvent: user updated event", zap.String("user_id", event.EntityID))
	case kafka.EventUserDeleted:
		userID, err := uuid.Parse(event.EntityID)
		if err != nil {
			s.log.Error(ctx, "handleUserEvent: invalid user id", zap.Error(err))
			return nil
		}
		if err := s.repo.DeleteResumesByUser(ctx, userID); err != nil {
			s.log.Error(ctx, "handleUserEvent: delete resumes by user", zap.Error(err))
			return err
		}
		s.log.Info(ctx, "handleUserEvent: deleted user resumes", zap.String("user_id", userID.String()))
	default:
		s.log.Debug(ctx, "handleUserEvent: unsupported user event type", zap.String("type", event.Type), zap.String("user_id", event.EntityID))
	}
	return nil
}
