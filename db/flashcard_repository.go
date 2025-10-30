package db

import (
	"database/sql"
	"fmt"

	"github.com/akolybelnikov/flashcards/models"
)

type FlashcardRepository interface {
	Create(req *models.CreateFlashcardRequest) (*models.Flashcard, error)
	GetAll() ([]*models.Flashcard, error)
	GetByID(id int) (*models.Flashcard, error)
	Update(id int, req *models.UpdateFlashcardRequest) (*models.Flashcard, error)
	Delete(id int) error
	GetRandom() (*models.Flashcard, error)
}

type PostgresFlashcardRepository struct {
	db *sql.DB
}

func NewPostgresFlashcardRepository(db *sql.DB) *PostgresFlashcardRepository {
	return &PostgresFlashcardRepository{db: db}
}

func (r *PostgresFlashcardRepository) Create(req *models.CreateFlashcardRequest) (*models.Flashcard, error) {
	query := `INSERT INTO flashcards (question, answer) VALUES ($1, $2) RETURNING id, question, answer, created_at, updated_at`

	var flashcard models.Flashcard
	err := r.db.QueryRow(query, req.Question, req.Answer).Scan(
		&flashcard.ID,
		&flashcard.Question,
		&flashcard.Answer,
		&flashcard.CreatedAt,
		&flashcard.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &flashcard, nil
}

func (r *PostgresFlashcardRepository) GetAll() ([]*models.Flashcard, error) {
	query := `SELECT id, question, answer, created_at, updated_at FROM flashcards ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flashcards []*models.Flashcard
	for rows.Next() {
		var flashcard models.Flashcard
		err := rows.Scan(
			&flashcard.ID,
			&flashcard.Question,
			&flashcard.Answer,
			&flashcard.CreatedAt,
			&flashcard.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, &flashcard)
	}

	return flashcards, nil
}

func (r *PostgresFlashcardRepository) GetByID(id int) (*models.Flashcard, error) {
	query := `SELECT id, question, answer, created_at, updated_at FROM flashcards WHERE id = $1`

	var flashcard models.Flashcard
	err := r.db.QueryRow(query, id).Scan(
		&flashcard.ID,
		&flashcard.Question,
		&flashcard.Answer,
		&flashcard.CreatedAt,
		&flashcard.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("flashcard with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &flashcard, nil
}

func (r *PostgresFlashcardRepository) Update(id int, req *models.UpdateFlashcardRequest) (*models.Flashcard, error) {
	query := `UPDATE flashcards SET question = COALESCE($1, question), answer = COALESCE($2, answer), updated_at = CURRENT_TIMESTAMP WHERE id = $3 RETURNING id, question, answer, created_at, updated_at`

	var flashcard models.Flashcard
	err := r.db.QueryRow(query, req.Question, req.Answer, id).Scan(
		&flashcard.ID,
		&flashcard.Question,
		&flashcard.Answer,
		&flashcard.CreatedAt,
		&flashcard.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("flashcard with id %d not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &flashcard, nil
}

func (r *PostgresFlashcardRepository) Delete(id int) error {
	query := `DELETE FROM flashcards WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("flashcard with id %d not found", id)
	}

	return nil
}

func (r *PostgresFlashcardRepository) GetRandom() (*models.Flashcard, error) {
	query := `SELECT id, question, answer, created_at, updated_at FROM flashcards ORDER BY RANDOM() LIMIT 1`

	var flashcard models.Flashcard
	err := r.db.QueryRow(query).Scan(
		&flashcard.ID,
		&flashcard.Question,
		&flashcard.Answer,
		&flashcard.CreatedAt,
		&flashcard.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no flashcards found")
	}
	if err != nil {
		return nil, err
	}

	return &flashcard, nil
}
