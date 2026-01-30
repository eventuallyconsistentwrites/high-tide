package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eventuallyconsistentwrites/high-tide-server/internal/domain"
	"github.com/eventuallyconsistentwrites/high-tide-server/internal/service"
)

type PostRoutes struct {
	service *service.PostService
}

func NewPostRoutes(s *service.PostService) *PostRoutes {
	return &PostRoutes{service: s}
}

// Below defined methods call corresponding service methods to execute business logic determined by the request

func (h *PostRoutes) CreatePost(w http.ResponseWriter, r *http.Request) {
	var p domain.Post

	// 1. Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Call Service
	if err := h.service.CreatePost(r.Context(), &p); err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// 3. Respond
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *PostRoutes) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostRoutes) GetPost(w http.ResponseWriter, r *http.Request) {
	idUrl := r.URL.Query().Get("id") // The id is expected as a query param

	if idUrl == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idUrl, 10, 64)

	post, err := h.service.GetPost(r.Context(), id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostRoutes) DeletePost(w http.ResponseWriter, r *http.Request) {
	idUrl := r.URL.Query().Get("id") // The id is expected as a query param

	if idUrl == "" {
		http.Error(w, "Missing ID parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idUrl, 10, 64)

	post, err := h.service.GetPost(r.Context(), id)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	} else {
		err := h.service.DeletePost(r.Context(), post)
		if err != nil {
			http.Error(w, "Post deletion was unsuccessful", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}

// RegisterRoutes which will be used by server defined in main.go
func (h *PostRoutes) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /posts", h.CreatePost)
	mux.HandleFunc("GET /posts", h.GetAllPosts)
	mux.HandleFunc("GET /post", h.GetPost)       // id as query parameter expected
	mux.HandleFunc("DELETE /post", h.DeletePost) // id as query parameter expected
}
