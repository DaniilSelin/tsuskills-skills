package handler

import (
	"net/http"

	"tsuskills-skills/internal/delivery/dto"
	"tsuskills-skills/internal/delivery/mapper"
	"tsuskills-skills/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ════════════════════════════════════════════
// SKILLS
// ════════════════════════════════════════════

// SearchSkills GET /api/v1/skills?q=flutter
func (h *Handler) SearchSkills(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query().Get("q")

	list, code := h.svc.SearchSkills(ctx, query)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "SearchSkills")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, mapper.SkillsToDTO(list))
}

// CreateSkill POST /api/v1/skills
func (h *Handler) CreateSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateSkillRequest
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	sk, code := h.svc.CreateSkill(ctx, req.Name)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "CreateSkill")
		return
	}

	h.writeJSON(ctx, w, http.StatusCreated, mapper.SkillToDTO(*sk))
}

// DeleteSkill DELETE /api/v1/skills/{id}
func (h *Handler) DeleteSkill(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.intParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid skill ID")
		return
	}

	code := h.svc.DeleteSkill(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "DeleteSkill")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "skill deleted"})
}

// ════════════════════════════════════════════
// ORGANIZATIONS
// ════════════════════════════════════════════

// CreateOrganization POST /api/v1/organizations
func (h *Handler) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateOrganizationRequest
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	org, err := mapper.OrgFromCreate(req)
	if err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid director_id")
		return
	}

	id, code := h.svc.CreateOrganization(ctx, org)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "CreateOrganization")
		return
	}

	h.log.Info(ctx, "CreateOrganization: success", zap.String("id", id.String()))
	h.writeJSON(ctx, w, http.StatusCreated, map[string]string{"id": id.String(), "message": "organization created"})
}

// GetOrganization GET /api/v1/organizations/{id}
func (h *Handler) GetOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.uuidParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid organization ID")
		return
	}

	org, code := h.svc.GetOrganization(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "GetOrganization")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, mapper.OrgToDTO(org))
}

// GetMyOrganization GET /api/v1/organizations/my?director_id=...
func (h *Handler) GetMyOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	directorID, err := uuid.Parse(r.URL.Query().Get("director_id"))
	if err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid or missing director_id")
		return
	}

	org, code := h.svc.GetMyOrganization(ctx, directorID)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "GetMyOrganization")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, mapper.OrgToDTO(org))
}

// UpdateOrganization PUT /api/v1/organizations/{id}
func (h *Handler) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.uuidParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid organization ID")
		return
	}

	var req dto.UpdateOrganizationRequest
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	// получаем текущую для partial update
	existing, code := h.svc.GetOrganization(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "UpdateOrganization")
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.AboutUs != "" {
		existing.AboutUs = req.AboutUs
	}

	code = h.svc.UpdateOrganization(ctx, existing)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "UpdateOrganization")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "organization updated"})
}

// DeleteOrganization DELETE /api/v1/organizations/{id}
func (h *Handler) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.uuidParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid organization ID")
		return
	}

	code := h.svc.DeleteOrganization(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "DeleteOrganization")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "organization deleted"})
}

// ════════════════════════════════════════════
// RESUMES
// ════════════════════════════════════════════

// CreateResume POST /api/v1/resumes
func (h *Handler) CreateResume(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateResumeRequest
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	res, err := mapper.ResumeFromCreate(req)
	if err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid user_id")
		return
	}

	id, code := h.svc.CreateResume(ctx, res)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "CreateResume")
		return
	}

	h.log.Info(ctx, "CreateResume: success", zap.String("id", id.String()))
	h.writeJSON(ctx, w, http.StatusCreated, map[string]string{"id": id.String(), "message": "resume created"})
}

// GetResume GET /api/v1/resumes/{id}
func (h *Handler) GetResume(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.uuidParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid resume ID")
		return
	}

	res, code := h.svc.GetResume(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "GetResume")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, mapper.ResumeToDTO(res))
}

// ListMyResumes GET /api/v1/resumes?user_id=...
func (h *Handler) ListMyResumes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := uuid.Parse(r.URL.Query().Get("user_id"))
	if err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid or missing user_id")
		return
	}

	list, code := h.svc.ListMyResumes(ctx, userID)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "ListMyResumes")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, mapper.ResumesToDTO(list))
}

// UpdateResume PUT /api/v1/resumes/{id}
func (h *Handler) UpdateResume(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.uuidParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid resume ID")
		return
	}

	var req dto.UpdateResumeRequest
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	existing, code := h.svc.GetResume(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "UpdateResume")
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.AboutMe != "" {
		existing.AboutMe = req.AboutMe
	}
	if req.SkillNames != nil {
		skills := make([]domain.Skill, 0, len(req.SkillNames))
		for _, name := range req.SkillNames {
			skills = append(skills, domain.Skill{Name: name})
		}
		existing.Skills = skills
	}

	code = h.svc.UpdateResume(ctx, existing)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "UpdateResume")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "resume updated"})
}

// DeleteResume DELETE /api/v1/resumes/{id}
func (h *Handler) DeleteResume(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.uuidParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid resume ID")
		return
	}

	code := h.svc.DeleteResume(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "DeleteResume")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, map[string]string{"message": "resume deleted"})
}

// ════════════════════════════════════════════
// APPLICATIONS
// ════════════════════════════════════════════

// CreateApplication POST /api/v1/applications
func (h *Handler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateApplicationRequest
	if !h.decodeAndValidate(ctx, w, r, &req) {
		return
	}

	resumeID, err := uuid.Parse(req.ResumeID)
	if err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid resume_id")
		return
	}
	vacancyID, err := uuid.Parse(req.VacancyID)
	if err != nil {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid vacancy_id")
		return
	}

	app := &domain.Application{ResumeID: resumeID, VacancyID: vacancyID}

	id, code := h.svc.CreateApplication(ctx, app)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "CreateApplication")
		return
	}

	h.log.Info(ctx, "CreateApplication: success", zap.String("id", id.String()))
	h.writeJSON(ctx, w, http.StatusCreated, map[string]string{"id": id.String(), "message": "application created"})
}

// GetApplicationsByVacancy GET /api/v1/applications?vacancy_id=...
func (h *Handler) ListApplications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// по vacancy_id — для работодателя
	if vid := r.URL.Query().Get("vacancy_id"); vid != "" {
		vacancyID, err := uuid.Parse(vid)
		if err != nil {
			h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid vacancy_id")
			return
		}
		list, code := h.svc.GetApplicationsByVacancy(ctx, vacancyID)
		if code != domain.CodeOK {
			h.handleSvcError(ctx, w, code, "ListApplications")
			return
		}
		h.writeJSON(ctx, w, http.StatusOK, mapper.AppsToDTO(list))
		return
	}

	// по user_id — для соискателя (мои отклики)
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		userID, err := uuid.Parse(uid)
		if err != nil {
			h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid user_id")
			return
		}
		list, code := h.svc.GetMyApplications(ctx, userID)
		if code != domain.CodeOK {
			h.handleSvcError(ctx, w, code, "ListApplications")
			return
		}
		h.writeJSON(ctx, w, http.StatusOK, mapper.AppsToDTO(list))
		return
	}

	h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Provide vacancy_id or user_id")
}

// GetApplication GET /api/v1/applications/{id}
func (h *Handler) GetApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, ok := h.uuidParam(r, "id")
	if !ok {
		h.writeError(ctx, w, http.StatusBadRequest, domain.CodeInvalidRequestBody, "Invalid application ID")
		return
	}

	app, code := h.svc.GetApplication(ctx, id)
	if code != domain.CodeOK {
		h.handleSvcError(ctx, w, code, "GetApplication")
		return
	}

	h.writeJSON(ctx, w, http.StatusOK, mapper.AppToDTO(app))
}
