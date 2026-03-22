package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"tsuskills-skills/internal/delivery/validator"
	"tsuskills-skills/internal/domain"
	"tsuskills-skills/internal/logger"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type IService interface {
	SearchSkills(ctx context.Context, query string) ([]domain.Skill, domain.ErrorCode)
	CreateSkill(ctx context.Context, name string) (*domain.Skill, domain.ErrorCode)
	DeleteSkill(ctx context.Context, id int) domain.ErrorCode

	CreateOrganization(ctx context.Context, o *domain.Organization) (uuid.UUID, domain.ErrorCode)
	GetOrganization(ctx context.Context, id uuid.UUID) (*domain.Organization, domain.ErrorCode)
	GetMyOrganization(ctx context.Context, directorID uuid.UUID) (*domain.Organization, domain.ErrorCode)
	UpdateOrganization(ctx context.Context, o *domain.Organization) domain.ErrorCode
	DeleteOrganization(ctx context.Context, id uuid.UUID) domain.ErrorCode

	CreateResume(ctx context.Context, res *domain.Resume) (uuid.UUID, domain.ErrorCode)
	GetResume(ctx context.Context, id uuid.UUID) (*domain.Resume, domain.ErrorCode)
	ListMyResumes(ctx context.Context, userID uuid.UUID) ([]domain.Resume, domain.ErrorCode)
	UpdateResume(ctx context.Context, res *domain.Resume) domain.ErrorCode
	DeleteResume(ctx context.Context, id uuid.UUID) domain.ErrorCode

	CreateApplication(ctx context.Context, app *domain.Application) (uuid.UUID, domain.ErrorCode)
	GetApplicationsByVacancy(ctx context.Context, vacancyID uuid.UUID) ([]domain.Application, domain.ErrorCode)
	GetMyApplications(ctx context.Context, userID uuid.UUID) ([]domain.Application, domain.ErrorCode)
	GetApplication(ctx context.Context, id uuid.UUID) (*domain.Application, domain.ErrorCode)
}

type Handler struct {
	svc IService
	log logger.Logger
}

func NewHandler(svc IService, l logger.Logger) *Handler {
	return &Handler{svc: svc, log: l}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message"`
}

func (h *Handler) writeJSON(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(ctx context.Context, w http.ResponseWriter, status int, code domain.ErrorCode, msg string) {
	h.writeJSON(ctx, w, status, ErrorResponse{Error: http.StatusText(status), Code: string(code), Message: msg})
}

func (h *Handler) handleSvcError(ctx context.Context, w http.ResponseWriter, code domain.ErrorCode, op string) {
	switch code {
	case domain.CodeNotFound:
		h.writeError(ctx, w, http.StatusNotFound, code, "Not found")
	case domain.CodeConflict:
		h.writeError(ctx, w, http.StatusConflict, code, "Already exists")
	default:
		h.log.Error(ctx, op, zap.String("code", string(code)))
		h.writeError(ctx, w, http.StatusInternalServerError, domain.CodeInternal, "Internal error")
	}
}

func (h *Handler) decodeAndValidate(ctx context.Context, w http.ResponseWriter, r *http.Request, dst interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid JSON")
		return false
	}
	if err := validator.ValidateStruct(dst); err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, err.Error())
		return false
	}
	return true
}

func (h *Handler) uuidParam(r *http.Request, name string) (uuid.UUID, bool) {
	raw := mux.Vars(r)[name]
	id, err := uuid.Parse(raw)
	return id, err == nil
}

func (h *Handler) intParam(r *http.Request, name string) (int, bool) {
	raw := mux.Vars(r)[name]
	v, err := strconv.Atoi(raw)
	return v, err == nil
}
