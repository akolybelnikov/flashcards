package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akolybelnikov/flashcards/models"
	"github.com/akolybelnikov/flashcards/services"

	"github.com/gorilla/mux"
)

type TodoHandler struct {
	service *services.TodoService
}

func NewTodoHandler(service *services.TodoService) *TodoHandler {
	return &TodoHandler{service: service}
}

func (h *TodoHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/todos", h.CreateTodo).Methods("POST")
	router.HandleFunc("/todos", h.GetAllTodos).Methods("GET")
	router.HandleFunc("/todos/{id:[0-9]+}", h.GetTodoByID).Methods("GET")
	router.HandleFunc("/todos/{id:[0-9]+}", h.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id:[0-9]+}", h.DeleteTodo).Methods("DELETE")
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	todo, err := h.service.CreateTodo(&req)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, todo)
}

func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetAllTodos()
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todos")
		return
	}

	h.writeJSONResponse(w, http.StatusOK, todos)
}

func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		if containsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve todo")
		}
		return
	}

	h.writeJSONResponse(w, http.StatusOK, todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	var req models.UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	todo, err := h.service.UpdateTodo(id, &req)
	if err != nil {
		if containsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	h.writeJSONResponse(w, http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid todo ID")
		return
	}

	err = h.service.DeleteTodo(id)
	if err != nil {
		if containsNotFound(err.Error()) {
			h.writeErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to delete todo")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (h *TodoHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func containsNotFound(message string) bool {
	return len(message) > 0 && (message[len(message)-9:] == "not found" ||
		message[:len("todo with id")] == "todo with id")
}
