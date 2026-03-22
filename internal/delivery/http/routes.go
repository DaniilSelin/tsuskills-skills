package http

import (
	httpBase "net/http"
	"tsuskills-skills/internal/logger"

	"github.com/gorilla/mux"
)

type IHandler interface {
	// Skills
	SearchSkills(w httpBase.ResponseWriter, r *httpBase.Request)
	CreateSkill(w httpBase.ResponseWriter, r *httpBase.Request)
	DeleteSkill(w httpBase.ResponseWriter, r *httpBase.Request)

	// Organizations
	CreateOrganization(w httpBase.ResponseWriter, r *httpBase.Request)
	GetOrganization(w httpBase.ResponseWriter, r *httpBase.Request)
	GetMyOrganization(w httpBase.ResponseWriter, r *httpBase.Request)
	UpdateOrganization(w httpBase.ResponseWriter, r *httpBase.Request)
	DeleteOrganization(w httpBase.ResponseWriter, r *httpBase.Request)

	// Resumes
	CreateResume(w httpBase.ResponseWriter, r *httpBase.Request)
	GetResume(w httpBase.ResponseWriter, r *httpBase.Request)
	ListMyResumes(w httpBase.ResponseWriter, r *httpBase.Request)
	UpdateResume(w httpBase.ResponseWriter, r *httpBase.Request)
	DeleteResume(w httpBase.ResponseWriter, r *httpBase.Request)

	// Applications
	CreateApplication(w httpBase.ResponseWriter, r *httpBase.Request)
	ListApplications(w httpBase.ResponseWriter, r *httpBase.Request)
	GetApplication(w httpBase.ResponseWriter, r *httpBase.Request)
}

func NewRouter(h IHandler, log logger.Logger) *mux.Router {
	r := mux.NewRouter()

	r.Use(RequestIDMiddleware)
	r.Use(CORSMiddleware)
	r.Use(LoggingMiddleware(log))
	r.Use(RecoveryMiddleware(log))

	// ── Skills ────────────────────
	sk := r.PathPrefix("/api/v1/skills").Subrouter()
	sk.HandleFunc("", h.SearchSkills).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	sk.HandleFunc("", h.CreateSkill).Methods(httpBase.MethodPost, httpBase.MethodOptions)
	sk.HandleFunc("/{id:[0-9]+}", h.DeleteSkill).Methods(httpBase.MethodDelete, httpBase.MethodOptions)

	// ── Organizations ─────────────
	org := r.PathPrefix("/api/v1/organizations").Subrouter()
	org.HandleFunc("", h.CreateOrganization).Methods(httpBase.MethodPost, httpBase.MethodOptions)
	org.HandleFunc("/my", h.GetMyOrganization).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	org.HandleFunc("/{id}", h.GetOrganization).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	org.HandleFunc("/{id}", h.UpdateOrganization).Methods(httpBase.MethodPut, httpBase.MethodOptions)
	org.HandleFunc("/{id}", h.DeleteOrganization).Methods(httpBase.MethodDelete, httpBase.MethodOptions)

	// ── Resumes ───────────────────
	res := r.PathPrefix("/api/v1/resumes").Subrouter()
	res.HandleFunc("", h.CreateResume).Methods(httpBase.MethodPost, httpBase.MethodOptions)
	res.HandleFunc("", h.ListMyResumes).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	res.HandleFunc("/{id}", h.GetResume).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	res.HandleFunc("/{id}", h.UpdateResume).Methods(httpBase.MethodPut, httpBase.MethodOptions)
	res.HandleFunc("/{id}", h.DeleteResume).Methods(httpBase.MethodDelete, httpBase.MethodOptions)

	// ── Applications ──────────────
	app := r.PathPrefix("/api/v1/applications").Subrouter()
	app.HandleFunc("", h.CreateApplication).Methods(httpBase.MethodPost, httpBase.MethodOptions)
	app.HandleFunc("", h.ListApplications).Methods(httpBase.MethodGet, httpBase.MethodOptions)
	app.HandleFunc("/{id}", h.GetApplication).Methods(httpBase.MethodGet, httpBase.MethodOptions)

	// Health
	r.HandleFunc("/health", func(w httpBase.ResponseWriter, r *httpBase.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpBase.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods(httpBase.MethodGet)

	return r
}
