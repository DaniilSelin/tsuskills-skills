CREATE SCHEMA IF NOT EXISTS skills;
SET search_path TO skills;

-- ── Справочник навыков ───────────────────
CREATE TABLE IF NOT EXISTS skills (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ── Организации ──────────────────────────
CREATE TABLE IF NOT EXISTS organizations (
    id          UUID PRIMARY KEY,
    director_id UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL,
    about_us    TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_org_director ON organizations(director_id);

-- ── Резюме ───────────────────────────────
CREATE TABLE IF NOT EXISTS resumes (
    id          UUID PRIMARY KEY,
    user_id     UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    about_me    TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_resume_user ON resumes(user_id);

CREATE TABLE IF NOT EXISTS resume_skills (
    resume_id UUID NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    skill_id  INT  NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (resume_id, skill_id)
);

-- ── Отклики ──────────────────────────────
CREATE TABLE IF NOT EXISTS applications (
    id         UUID PRIMARY KEY,
    resume_id  UUID        NOT NULL REFERENCES resumes(id) ON DELETE CASCADE,
    vacancy_id UUID        NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_app_resume_vacancy UNIQUE (resume_id, vacancy_id)
);

CREATE INDEX IF NOT EXISTS idx_app_vacancy ON applications(vacancy_id);
CREATE INDEX IF NOT EXISTS idx_app_resume  ON applications(resume_id);
