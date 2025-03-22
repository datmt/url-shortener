package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Handler struct {
	Store    *Storage
	AdminKey string
}

type ShortenRequest struct {
	Target string `json:"target"`
	Handle string `json:"handle"`
}

func (h *Handler) CreateOrUpdateShortLink(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	var req ShortenRequest
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &req); err != nil || req.Target == "" || req.Handle == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Store.SaveLink(req.Handle, req.Target, username); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetShortLink(w http.ResponseWriter, r *http.Request) {
	handle := strings.TrimPrefix(r.URL.Path, "/shorten/")
	target, err := h.Store.GetTarget(handle)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Write([]byte(target))
}

func (h *Handler) DeleteShortLink(w http.ResponseWriter, r *http.Request) {
	handle := strings.TrimPrefix(r.URL.Path, "/delete/")
	username := r.Context().Value("username").(string)
	if err := h.Store.DeleteLink(handle, username); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	handle := strings.TrimPrefix(r.URL.Path, "/r/")
	target, err := h.Store.GetTarget(handle)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, target, http.StatusFound)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Admin-Key") != h.AdminKey {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var creds Credentials
	body, _ := io.ReadAll(r.Body)
	if err := json.Unmarshal(body, &creds); err != nil || creds.Username == "" || creds.Password == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if err := h.Store.CreateUser(creds.Username, creds.Password); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ListUserLinks(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	rows, err := h.Store.DB.Query("SELECT handle, target FROM links WHERE owner = ?", username)
	if err != nil {
		http.Error(w, "Failed to load links", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Link struct {
		Handle string `json:"handle"`
		Target string `json:"target"`
	}

	links := []Link{}
	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.Handle, &l.Target); err == nil {
			links = append(links, l)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
