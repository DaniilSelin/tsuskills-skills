package mapper

import (
	"tsuskills-skills/internal/delivery/dto"
	"tsuskills-skills/internal/domain"

	"github.com/google/uuid"
)

func SkillToDTO(s domain.Skill) dto.SkillResponse {
	return dto.SkillResponse{ID: s.ID, Name: s.Name, CreatedAt: s.CreatedAt}
}

func SkillsToDTO(list []domain.Skill) []dto.SkillResponse {
	out := make([]dto.SkillResponse, 0, len(list))
	for _, s := range list {
		out = append(out, SkillToDTO(s))
	}
	return out
}

func OrgToDTO(o *domain.Organization) dto.OrganizationResponse {
	return dto.OrganizationResponse{
		ID: o.ID.String(), DirectorID: o.DirectorID.String(),
		Name: o.Name, AboutUs: o.AboutUs,
		CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt,
	}
}

func OrgFromCreate(r dto.CreateOrganizationRequest) (*domain.Organization, error) {
	dirID, err := uuid.Parse(r.DirectorID)
	if err != nil {
		return nil, err
	}
	return &domain.Organization{DirectorID: dirID, Name: r.Name, AboutUs: r.AboutUs}, nil
}

func ResumeToDTO(r *domain.Resume) dto.ResumeResponse {
	skills := make([]dto.SkillResponse, 0, len(r.Skills))
	for _, s := range r.Skills {
		skills = append(skills, SkillToDTO(s))
	}
	return dto.ResumeResponse{
		ID: r.ID.String(), UserID: r.UserID.String(),
		Name: r.Name, Description: r.Description, AboutMe: r.AboutMe,
		Skills: skills, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

func ResumesToDTO(list []domain.Resume) []dto.ResumeResponse {
	out := make([]dto.ResumeResponse, 0, len(list))
	for _, r := range list {
		out = append(out, ResumeToDTO(&r))
	}
	return out
}

func ResumeFromCreate(r dto.CreateResumeRequest) (*domain.Resume, error) {
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, err
	}
	skills := make([]domain.Skill, 0, len(r.SkillNames))
	for _, name := range r.SkillNames {
		skills = append(skills, domain.Skill{Name: name})
	}
	return &domain.Resume{
		UserID: userID, Name: r.Name,
		Description: r.Description, AboutMe: r.AboutMe, Skills: skills,
	}, nil
}

func AppToDTO(a *domain.Application) dto.ApplicationResponse {
	return dto.ApplicationResponse{
		ID: a.ID.String(), ResumeID: a.ResumeID.String(),
		VacancyID: a.VacancyID.String(), Status: string(a.Status),
		ResumeName: a.ResumeName,
		CreatedAt: a.CreatedAt, UpdatedAt: a.UpdatedAt,
	}
}

func AppsToDTO(list []domain.Application) []dto.ApplicationResponse {
	out := make([]dto.ApplicationResponse, 0, len(list))
	for _, a := range list {
		out = append(out, AppToDTO(&a))
	}
	return out
}
