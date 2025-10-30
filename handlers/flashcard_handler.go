package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/akolybelnikov/flashcards/models"
	"github.com/akolybelnikov/flashcards/services"

	"github.com/gorilla/mux"
)

type FlashcardHandler struct {
	service services.FlashcardServiceInterface
}

func NewFlashcardHandler(service services.FlashcardServiceInterface) *FlashcardHandler {
	if service == nil {
		panic("service is nil")
	}
	return &FlashcardHandler{service: service}
}

func (h *FlashcardHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/flashcards", h.CreateFlashcard).Methods("POST")
	router.HandleFunc("/flashcards", h.GetAllFlashcards).Methods("GET")
	router.HandleFunc("/flashcards/random", h.GetRandomFlashcard).Methods("GET")
	router.HandleFunc("/flashcards/{id:[0-9]+}", h.GetFlashcardByID).Methods("GET")
	router.HandleFunc("/flashcards/{id:[0-9]+}", h.UpdateFlashcard).Methods("PUT")
	router.HandleFunc("/flashcards/{id:[0-9]+}", h.DeleteFlashcard).Methods("DELETE")
}

func (h *FlashcardHandler) CreateFlashcard(w http.ResponseWriter, r *http.Request) {
	var req models.CreateFlashcardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	flashcard, err := h.service.CreateFlashcard(&req)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, flashcard)
}

func (h *FlashcardHandler) GetAllFlashcards(w http.ResponseWriter, _ *http.Request) {
	flashcards, err := h.service.GetAllFlashcards()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve flashcards")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, flashcards)
}

func (h *FlashcardHandler) GetRandomFlashcard(w http.ResponseWriter, r *http.Request) {
	// Optional query param 'lang' for desired AI hint language (e.g., 'el' for Greek)
	lang := r.URL.Query().Get("lang")

	flashcard, err := h.service.GetRandomFlashcard()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve random flashcard")
		return
	}

	var aiHint *string
	// Try to generate an AI hint but don't fail the request if AI is unavailable
	aiHint = h.service.GenerateAIHint(flashcard, lang)

	resp := models.RandomFlashcardResponse{
		Flashcard: flashcard,
		AIHint:    aiHint,
	}

	h.writeJSONResponse(w, http.StatusOK, resp)
}

func (h *FlashcardHandler) GetFlashcardByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid flashcard ID")
		return
	}

	flashcard, err := h.service.GetFlashcardByID(id)
	if err != nil {
		if containsNotFoundFlashcard(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve flashcard")
		}
		return
	}

	h.writeJSONResponse(w, http.StatusOK, flashcard)
}

func (h *FlashcardHandler) UpdateFlashcard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid flashcard ID")
		return
	}

	var req models.UpdateFlashcardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	flashcard, err := h.service.UpdateFlashcard(id, &req)
	if err != nil {
		if containsNotFoundFlashcard(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.writeJSONResponse(w, http.StatusOK, flashcard)
}

func (h *FlashcardHandler) DeleteFlashcard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid flashcard ID")
		return
	}

	err = h.service.DeleteFlashcard(id)
	if err != nil {
		if containsNotFoundFlashcard(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete flashcard")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FlashcardHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func (h *FlashcardHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(map[string]string{"error": message})
	if err != nil {
		return
	}
}

func containsNotFoundFlashcard(message string) bool {
	return strings.Contains(message, "not found") || strings.Contains(message, "flashcard with id")
}
