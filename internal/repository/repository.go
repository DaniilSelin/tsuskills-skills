package repository

import (
	"context"
	"errors"
	"fmt"

	"tsuskills-skills/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// ════════════════ SKILLS ════════════════

func (r *Repository) SearchSkills(ctx context.Context, query string) ([]domain.Skill, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, created_at FROM skills WHERE name ILIKE '%' || $1 || '%' ORDER BY name LIMIT 50`, query)
	if err != nil {
		return nil, fmt.Errorf("search skills: %w", err)
	}
	defer rows.Close()

	var result []domain.Skill
	for rows.Next() {
		var s domain.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, rows.Err()
}

func (r *Repository) CreateSkill(ctx context.Context, name string) (*domain.Skill, error) {
	var s domain.Skill
	err := r.pool.QueryRow(ctx,
		`INSERT INTO skills (name) VALUES ($1)
		 ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, created_at`, name,
	).Scan(&s.ID, &s.Name, &s.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create skill: %w", err)
	}
	return &s, nil
}

func (r *Repository) DeleteSkill(ctx context.Context, id int) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM skills WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete skill: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ════════════════ ORGANIZATIONS ════════════════

func (r *Repository) CreateOrganization(ctx context.Context, o *domain.Organization) (uuid.UUID, error) {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO organizations (id, director_id, name, about_us, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		o.ID, o.DirectorID, o.Name, o.AboutUs, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create org: %w", err)
	}
	return o.ID, nil
}

func (r *Repository) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	var o domain.Organization
	err := r.pool.QueryRow(ctx,
		`SELECT id, director_id, name, about_us, created_at, updated_at
		 FROM organizations WHERE id = $1`, id,
	).Scan(&o.ID, &o.DirectorID, &o.Name, &o.AboutUs, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get org: %w", err)
	}
	return &o, nil
}

func (r *Repository) GetOrganizationByDirector(ctx context.Context, directorID uuid.UUID) (*domain.Organization, error) {
	var o domain.Organization
	err := r.pool.QueryRow(ctx,
		`SELECT id, director_id, name, about_us, created_at, updated_at
		 FROM organizations WHERE director_id = $1 LIMIT 1`, directorID,
	).Scan(&o.ID, &o.DirectorID, &o.Name, &o.AboutUs, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &o, nil
}

func (r *Repository) UpdateOrganization(ctx context.Context, o *domain.Organization) error {
	tag, err := r.pool.Exec(ctx,
		`UPDATE organizations SET name=$1, about_us=$2, updated_at=NOW() WHERE id=$3`,
		o.Name, o.AboutUs, o.ID)
	if err != nil {
		return fmt.Errorf("update org: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM organizations WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ════════════════ RESUMES ════════════════

func (r *Repository) CreateResume(ctx context.Context, res *domain.Resume) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO resumes (id, user_id, name, description, about_me, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		res.ID, res.UserID, res.Name, res.Description, res.AboutMe, res.CreatedAt, res.UpdatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert resume: %w", err)
	}

	for _, sk := range res.Skills {
		var skillID int
		err = tx.QueryRow(ctx,
			`INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id`,
			sk.Name).Scan(&skillID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("upsert skill: %w", err)
		}
		_, err = tx.Exec(ctx, `INSERT INTO resume_skills (resume_id, skill_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			res.ID, skillID)
		if err != nil {
			return uuid.Nil, fmt.Errorf("link skill: %w", err)
		}
	}

	return res.ID, tx.Commit(ctx)
}

func (r *Repository) GetResumeByID(ctx context.Context, id uuid.UUID) (*domain.Resume, error) {
	var res domain.Resume
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, description, about_me, created_at, updated_at
		 FROM resumes WHERE id = $1`, id,
	).Scan(&res.ID, &res.UserID, &res.Name, &res.Description, &res.AboutMe, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	skills, _ := r.getResumeSkills(ctx, id)
	res.Skills = skills
	return &res, nil
}

func (r *Repository) ListResumesByUser(ctx context.Context, userID uuid.UUID) ([]domain.Resume, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, description, about_me, created_at, updated_at
		 FROM resumes WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Resume
	for rows.Next() {
		var res domain.Resume
		if err := rows.Scan(&res.ID, &res.UserID, &res.Name, &res.Description, &res.AboutMe, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}
		skills, _ := r.getResumeSkills(ctx, res.ID)
		res.Skills = skills
		result = append(result, res)
	}
	return result, rows.Err()
}

func (r *Repository) UpdateResume(ctx context.Context, res *domain.Resume) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx,
		`UPDATE resumes SET name=$1, description=$2, about_me=$3, updated_at=NOW() WHERE id=$4`,
		res.Name, res.Description, res.AboutMe, res.ID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	_, _ = tx.Exec(ctx, `DELETE FROM resume_skills WHERE resume_id = $1`, res.ID)
	for _, sk := range res.Skills {
		var skillID int
		err = tx.QueryRow(ctx,
			`INSERT INTO skills (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name RETURNING id`,
			sk.Name).Scan(&skillID)
		if err != nil {
			return err
		}
		_, _ = tx.Exec(ctx, `INSERT INTO resume_skills (resume_id, skill_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, res.ID, skillID)
	}

	return tx.Commit(ctx)
}

func (r *Repository) DeleteResume(ctx context.Context, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM resumes WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) getResumeSkills(ctx context.Context, resumeID uuid.UUID) ([]domain.Skill, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT s.id, s.name, s.created_at FROM skills s
		 JOIN resume_skills rs ON rs.skill_id = s.id WHERE rs.resume_id = $1 ORDER BY s.name`, resumeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var skills []domain.Skill
	for rows.Next() {
		var s domain.Skill
		if err := rows.Scan(&s.ID, &s.Name, &s.CreatedAt); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}
	return skills, nil
}

// ════════════════ APPLICATIONS ════════════════

func (r *Repository) CreateApplication(ctx context.Context, app *domain.Application) (uuid.UUID, error) {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO applications (id, resume_id, vacancy_id, status, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`,
		app.ID, app.ResumeID, app.VacancyID, string(app.Status), app.CreatedAt, app.UpdatedAt)
	if err != nil {
		return uuid.Nil, fmt.Errorf("create app: %w", err)
	}
	return app.ID, nil
}

func (r *Repository) GetApplicationsByVacancy(ctx context.Context, vacancyID uuid.UUID) ([]domain.Application, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT a.id, a.resume_id, a.vacancy_id, a.status, a.created_at, a.updated_at,
		        COALESCE(res.name, '') AS resume_name
		 FROM applications a
		 LEFT JOIN resumes res ON res.id = a.resume_id
		 WHERE a.vacancy_id = $1
		 ORDER BY a.created_at DESC`, vacancyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanApplications(rows)
}

func (r *Repository) GetApplicationsByUser(ctx context.Context, userID uuid.UUID) ([]domain.Application, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT a.id, a.resume_id, a.vacancy_id, a.status, a.created_at, a.updated_at,
		        COALESCE(res.name, '') AS resume_name
		 FROM applications a
		 JOIN resumes res ON res.id = a.resume_id
		 WHERE res.user_id = $1
		 ORDER BY a.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.scanApplications(rows)
}

func (r *Repository) GetApplicationByID(ctx context.Context, id uuid.UUID) (*domain.Application, error) {
	var a domain.Application
	err := r.pool.QueryRow(ctx,
		`SELECT a.id, a.resume_id, a.vacancy_id, a.status, a.created_at, a.updated_at,
		        COALESCE(res.name, '') AS resume_name
		 FROM applications a
		 LEFT JOIN resumes res ON res.id = a.resume_id
		 WHERE a.id = $1`, id,
	).Scan(&a.ID, &a.ResumeID, &a.VacancyID, &a.Status, &a.CreatedAt, &a.UpdatedAt, &a.ResumeName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) scanApplications(rows pgx.Rows) ([]domain.Application, error) {
	var result []domain.Application
	for rows.Next() {
		var a domain.Application
		if err := rows.Scan(&a.ID, &a.ResumeID, &a.VacancyID, &a.Status, &a.CreatedAt, &a.UpdatedAt, &a.ResumeName); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, rows.Err()
}
