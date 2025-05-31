package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"backend-blog/backend/config"
	"backend-blog/backend/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetPosts(w http.ResponseWriter, r *http.Request) {
	resp, err := config.MakeRequest("GET", "/posts", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Erro ao buscar posts", resp.StatusCode)
		return
	}

	var posts []models.Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	resp, err := config.MakeRequest("GET", "/posts?id=eq."+params["id"], nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Post não encontrado", http.StatusNotFound)
		return
	}

	var posts []models.Post
	if err := json.NewDecoder(resp.Body).Decode(&posts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(posts) == 0 {
		http.Error(w, "Post não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts[0])
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Gerar UUID para o novo post
	post.ID = uuid.New().String()
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	log.Printf("Tentando criar post: %+v", post)

	resp, err := config.MakeRequest("POST", "/posts", post)
	if err != nil {
		log.Printf("Erro ao fazer request para Supabase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Erro do Supabase (Status %d): %s", resp.StatusCode, string(body))
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	var result []models.Post
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Erro ao decodificar resposta: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result[0])
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "ID do post não fornecido", http.StatusBadRequest)
		return
	}

	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post.ID = id // Garantir que o ID seja o mesmo da URL
	post.UpdatedAt = time.Now()

	log.Printf("Atualizando post %s: %+v", id, post)

	resp, err := config.MakeRequest("PATCH", "/posts?id=eq."+id, post)
	if err != nil {
		log.Printf("Erro ao fazer request para Supabase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Erro do Supabase (Status %d): %s", resp.StatusCode, string(body))
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	var result []models.Post
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Erro ao decodificar resposta: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(result) == 0 {
		http.Error(w, "Post não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result[0])
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "ID do post não fornecido", http.StatusBadRequest)
		return
	}

	log.Printf("Deletando post %s", id)

	resp, err := config.MakeRequest("DELETE", "/posts?id=eq."+id, nil)
	if err != nil {
		log.Printf("Erro ao fazer request para Supabase: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Erro do Supabase (Status %d): %s", resp.StatusCode, string(body))
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
