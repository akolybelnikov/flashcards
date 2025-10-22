package db

import (
	"database/sql"
	"fmt"

	"github.com/akolybelnikov/flashcards/models"

	_ "github.com/lib/pq"
)

type TodoRepository interface {
	CreateTodo(todo *models.Todo) error
	GetTodoByID(id int) (*models.Todo, error)
	GetAllTodos() ([]*models.Todo, error)
	UpdateTodo(id int, updates map[string]any) error
	DeleteTodo(id int) error
}

type PostgresTodoRepository struct {
	db *sql.DB
}

func NewPostgresTodoRepository(databaseURL string) (*PostgresTodoRepository, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresTodoRepository{db: db}, nil
}

func (r *PostgresTodoRepository) CreateTodo(todo *models.Todo) error {
	query := `
		INSERT INTO gocourse.todos (title, description, completed) 
		VALUES ($1, $2, $3) 
		RETURNING id, createdAt, updatedAt`

	row := r.db.QueryRow(query, todo.Title, todo.Description, todo.Completed)

	err := row.Scan(&todo.ID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

func (r *PostgresTodoRepository) GetTodoByID(id int) (*models.Todo, error) {
	query := `
		SELECT id, title, description, completed, createdAt, updatedAt 
		FROM gocourse.todos 
		WHERE id = $1`

	todo := &models.Todo{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}

	return todo, nil
}

func (r *PostgresTodoRepository) GetAllTodos() ([]*models.Todo, error) {
	query := `
		SELECT id, title, description, completed, createdAt, updatedAt 
		FROM gocourse.todos 
		ORDER BY createdAt DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	todos := make([]*models.Todo, 0)
	for rows.Next() {
		todo := &models.Todo{}
		err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over todos: %w", err)
	}

	return todos, nil
}

func (r *PostgresTodoRepository) UpdateTodo(id int, updates map[string]any) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	query := "UPDATE gocourse.todos SET "
	args := []any{}
	argIndex := 1

	for field, value := range updates {
		if argIndex > 1 {
			query += ", "
		}
		query += fmt.Sprintf("%s = $%d", field, argIndex)
		args = append(args, value)
		argIndex++
	}

	query += fmt.Sprintf(", updatedAt = NOW() WHERE id = $%d", argIndex)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}

func (r *PostgresTodoRepository) DeleteTodo(id int) error {
	query := "DELETE FROM gocourse.todos WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %d not found", id)
	}

	return nil
}

func (r *PostgresTodoRepository) Close() error {
	return r.db.Close()
}
