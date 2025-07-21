package tag

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kostinp/edu-platform-backend/pkg/db"
)

// Repository описывает интерфейс доступа к хранилищу тегов и связей с сущностями (курсы, уроки, тесты)
type Repository interface {
	// --- Tag CRUD ---
	Create(ctx context.Context, tag *Tag) error
	GetAll(ctx context.Context) ([]*Tag, error)
	Update(ctx context.Context, id uuid.UUID, dto UpdateTagDTO) error
	Delete(ctx context.Context, id uuid.UUID) error

	// --- Course tags ---
	AddTagToCourse(ctx context.Context, courseID, tagID uuid.UUID) error
	RemoveTagFromCourse(ctx context.Context, courseID, tagID uuid.UUID) error
	ListTagsByCourse(ctx context.Context, courseID uuid.UUID) ([]*Tag, error)

	// --- Lesson tags ---
	AddTagToLesson(ctx context.Context, lessonID, tagID uuid.UUID) error
	RemoveTagFromLesson(ctx context.Context, lessonID, tagID uuid.UUID) error
	ListTagsByLesson(ctx context.Context, lessonID uuid.UUID) ([]*Tag, error)

	// --- Test tags ---
	AddTagToTest(ctx context.Context, testID, tagID uuid.UUID) error
	RemoveTagFromTest(ctx context.Context, testID, tagID uuid.UUID) error
	ListTagsByTest(ctx context.Context, testID uuid.UUID) ([]*Tag, error)
}

// PostgresRepo реализует интерфейс Repository для PostgreSQL
type PostgresRepo struct{}

// NewPostgresRepo создаёт новый экземпляр PostgresRepo
func NewPostgresRepo() *PostgresRepo {
	return &PostgresRepo{}
}

// --- Tag CRUD ---

// Create сохраняет новый тег в базе данных
func (r *PostgresRepo) Create(ctx context.Context, tag *Tag) error {
	query := `INSERT INTO tags (id, name, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	now := time.Now()
	tag.ID = uuid.New()
	tag.CreatedAt = now
	tag.UpdatedAt = now

	_, err := db.DB().Exec(ctx, query, tag.ID, tag.Name, tag.AuthorID, tag.CreatedAt, tag.UpdatedAt)
	return err
}

// GetAll возвращает все активные (не удалённые) теги
func (r *PostgresRepo) GetAll(ctx context.Context) ([]*Tag, error) {
	rows, err := db.DB().Query(ctx, `SELECT id, name, author_id, created_at, updated_at FROM tags WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

// Update обновляет имя тега по его ID
func (r *PostgresRepo) Update(ctx context.Context, id uuid.UUID, dto UpdateTagDTO) error {
	query := `UPDATE tags SET name = $1, updated_at = $2 WHERE id = $3 AND deleted_at IS NULL`
	_, err := db.DB().Exec(ctx, query, dto.Name, time.Now(), id)
	return err
}

// Delete выполняет soft delete тега (устанавливает deleted_at)
func (r *PostgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tags SET deleted_at = $1 WHERE id = $2`
	_, err := db.DB().Exec(ctx, query, time.Now(), id)
	return err
}

// --- Course tags ---

// AddTagToCourse добавляет тег к курсу (без дубликатов)
func (r *PostgresRepo) AddTagToCourse(ctx context.Context, courseID, tagID uuid.UUID) error {
	query := `INSERT INTO course_tags (course_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.DB().Exec(ctx, query, courseID, tagID)
	return err
}

// RemoveTagFromCourse удаляет тег у курса
func (r *PostgresRepo) RemoveTagFromCourse(ctx context.Context, courseID, tagID uuid.UUID) error {
	query := `DELETE FROM course_tags WHERE course_id = $1 AND tag_id = $2`
	_, err := db.DB().Exec(ctx, query, courseID, tagID)
	return err
}

// ListTagsByCourse возвращает список тегов, связанных с курсом
func (r *PostgresRepo) ListTagsByCourse(ctx context.Context, courseID uuid.UUID) ([]*Tag, error) {
	query := `
		SELECT t.id, t.name, t.author_id, t.created_at, t.updated_at
		FROM tags t
		JOIN course_tags ct ON ct.tag_id = t.id
		WHERE ct.course_id = $1 AND t.deleted_at IS NULL
	`
	rows, err := db.DB().Query(ctx, query, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

// --- Lesson tags ---

// AddTagToLesson добавляет тег к уроку
func (r *PostgresRepo) AddTagToLesson(ctx context.Context, lessonID, tagID uuid.UUID) error {
	query := `INSERT INTO lesson_tags (lesson_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.DB().Exec(ctx, query, lessonID, tagID)
	return err
}

// RemoveTagFromLesson удаляет тег у урока
func (r *PostgresRepo) RemoveTagFromLesson(ctx context.Context, lessonID, tagID uuid.UUID) error {
	query := `DELETE FROM lesson_tags WHERE lesson_id = $1 AND tag_id = $2`
	_, err := db.DB().Exec(ctx, query, lessonID, tagID)
	return err
}

// ListTagsByLesson возвращает список тегов, связанных с уроком
func (r *PostgresRepo) ListTagsByLesson(ctx context.Context, lessonID uuid.UUID) ([]*Tag, error) {
	query := `
		SELECT t.id, t.name, t.author_id, t.created_at, t.updated_at
		FROM tags t
		JOIN lesson_tags lt ON lt.tag_id = t.id
		WHERE lt.lesson_id = $1 AND t.deleted_at IS NULL
	`
	rows, err := db.DB().Query(ctx, query, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}

// --- Test tags ---

// AddTagToTest добавляет тег к тесту
func (r *PostgresRepo) AddTagToTest(ctx context.Context, testID, tagID uuid.UUID) error {
	query := `INSERT INTO test_tags (test_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := db.DB().Exec(ctx, query, testID, tagID)
	return err
}

// RemoveTagFromTest удаляет тег у теста
func (r *PostgresRepo) RemoveTagFromTest(ctx context.Context, testID, tagID uuid.UUID) error {
	query := `DELETE FROM test_tags WHERE test_id = $1 AND tag_id = $2`
	_, err := db.DB().Exec(ctx, query, testID, tagID)
	return err
}

// ListTagsByTest возвращает список тегов, связанных с тестом
func (r *PostgresRepo) ListTagsByTest(ctx context.Context, testID uuid.UUID) ([]*Tag, error) {
	query := `
		SELECT t.id, t.name, t.author_id, t.created_at, t.updated_at
		FROM tags t
		JOIN test_tags tt ON tt.tag_id = t.id
		WHERE tt.test_id = $1 AND t.deleted_at IS NULL
	`
	rows, err := db.DB().Query(ctx, query, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		var t Tag
		if err := rows.Scan(&t.ID, &t.Name, &t.AuthorID, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, &t)
	}
	return tags, nil
}
