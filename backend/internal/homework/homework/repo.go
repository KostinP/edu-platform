package homework

import (
	"context"
	"fmt"
	"time"

	"github.com/kostinp/edu-platform-backend/pkg/db"

	"github.com/google/uuid"
)

type PostgresRepo struct{}

type Repository interface {
	Create(ctx context.Context, hw *Homework) error
	GetByID(ctx context.Context, id uuid.UUID) (*Homework, error)
	List(ctx context.Context) ([]Homework, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListForUserFiltered(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) ([]Homework, error)
	GetStatsForUser(ctx context.Context, userID uuid.UUID) (map[string]int, error)
	ListByAuthor(ctx context.Context, authorID uuid.UUID) ([]Homework, error)
	Filter(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Homework, error)
	ListByLesson(ctx context.Context, lessonID uuid.UUID) ([]Homework, error)
}

func NewPostgresRepo() Repository {
	return &PostgresRepo{}
}

func (r *PostgresRepo) Create(ctx context.Context, hw *Homework) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO homeworks
		(id, title, description, type, course_id, module_id, lesson_id, group_id, user_id, author_id, due_at, is_required, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
	`, hw.ID, hw.Title, hw.Description, hw.Type, hw.CourseID, hw.ModuleID, hw.LessonID, hw.GroupID, hw.UserID, hw.AuthorID, hw.DueAt, hw.IsRequired, time.Now(), time.Now())
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id uuid.UUID) (*Homework, error) {
	row := db.Pool.QueryRow(ctx, `
		SELECT id, title, description, type, course_id, module_id, lesson_id, group_id, user_id, author_id, due_at, is_required, created_at, updated_at
		FROM homeworks WHERE id = $1 AND deleted_at IS NULL
	`, id)

	var hw Homework
	err := row.Scan(&hw.ID, &hw.Title, &hw.Description, &hw.Type, &hw.CourseID, &hw.ModuleID, &hw.LessonID, &hw.GroupID, &hw.UserID, &hw.AuthorID, &hw.DueAt, &hw.IsRequired, &hw.CreatedAt, &hw.UpdatedAt)
	return &hw, err
}

func (r *PostgresRepo) List(ctx context.Context) ([]Homework, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, title, description, type, course_id, module_id, lesson_id, group_id, user_id, author_id, due_at, is_required, created_at, updated_at
		FROM homeworks WHERE deleted_at IS NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Homework
	for rows.Next() {
		var hw Homework
		err := rows.Scan(&hw.ID, &hw.Title, &hw.Description, &hw.Type, &hw.CourseID, &hw.ModuleID, &hw.LessonID, &hw.GroupID, &hw.UserID, &hw.AuthorID, &hw.DueAt, &hw.IsRequired, &hw.CreatedAt, &hw.UpdatedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, hw)
	}
	return result, nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.Pool.Exec(ctx, `
		UPDATE homeworks SET deleted_at = NOW() WHERE id = $1
	`, id)
	return err
}

func (r *PostgresRepo) ListForUserFiltered(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) ([]Homework, error) {
	query := `
		SELECT DISTINCT h.id, h.title, h.description, h.type, h.course_id, h.module_id,
		                h.lesson_id, h.group_id, h.user_id, h.author_id, h.due_at,
		                h.is_required, h.created_at, h.updated_at, h.deleted_at
		FROM homeworks h
		LEFT JOIN group_members gm ON h.group_id = gm.group_id
		WHERE (h.user_id = $1 OR gm.user_id = $1) AND h.deleted_at IS NULL
	`
	args := []interface{}{userID}
	argIdx := 2

	// Фильтрация
	if t, ok := filters["type"].(string); ok && t != "" {
		query += fmt.Sprintf(" AND h.type = $%d", argIdx)
		args = append(args, t)
		argIdx++
	}

	if r, ok := filters["required"].(*bool); ok && r != nil {
		query += fmt.Sprintf(" AND h.is_required = $%d", argIdx)
		args = append(args, *r)
		argIdx++
	}

	if c, ok := filters["course_id"].(uuid.UUID); ok {
		query += fmt.Sprintf(" AND h.course_id = $%d", argIdx)
		args = append(args, c)
		argIdx++
	}

	if status, ok := filters["status"].(string); ok {
		now := time.Now()
		switch status {
		case "active":
			query += fmt.Sprintf(" AND (h.due_at IS NULL OR h.due_at > $%d)", argIdx)
			args = append(args, now)
			argIdx++
		case "overdue":
			query += fmt.Sprintf(" AND h.due_at IS NOT NULL AND h.due_at <= $%d", argIdx)
			args = append(args, now)
			argIdx++
		}
	}

	if dueFrom, ok := filters["due_from"].(time.Time); ok {
		query += fmt.Sprintf(" AND h.due_at >= $%d", argIdx)
		args = append(args, dueFrom)
		argIdx++
	}

	if dueTo, ok := filters["due_to"].(time.Time); ok {
		query += fmt.Sprintf(" AND h.due_at <= $%d", argIdx)
		args = append(args, dueTo)
		argIdx++
	}

	// Сортировка
	if sortBy, ok := filters["sort_by"].(string); ok {
		switch sortBy {
		case "due_at", "created_at", "title":
			order := "ASC"
			if sortOrder, ok := filters["sort_order"].(string); ok && sortOrder == "desc" {
				order = "DESC"
			}
			query += fmt.Sprintf(" ORDER BY h.%s %s", sortBy, order)
		}
	}

	// Пагинация
	if limit, ok := filters["limit"].(int); ok && limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, limit)
		argIdx++
	}
	if offset, ok := filters["offset"].(int); ok && offset >= 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, offset)
		argIdx++
	}

	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Homework
	for rows.Next() {
		var h Homework
		err := rows.Scan(&h.ID, &h.Title, &h.Description, &h.Type, &h.CourseID, &h.ModuleID,
			&h.LessonID, &h.GroupID, &h.UserID, &h.AuthorID, &h.DueAt,
			&h.IsRequired, &h.CreatedAt, &h.UpdatedAt, &h.DeletedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PostgresRepo) GetStatsForUser(ctx context.Context, userID uuid.UUID) (map[string]int, error) {
	query := `
		SELECT 
			COUNT(*) FILTER (WHERE TRUE) AS total,
			COUNT(*) FILTER (WHERE s.status = 'submitted') AS submitted,
			COUNT(*) FILTER (WHERE s.status = 'checked') AS checked,
			COUNT(*) FILTER (
				WHERE h.due_at IS NOT NULL AND h.due_at < NOW() AND (s.status IS NULL OR s.status != 'checked')
			) AS overdue
		FROM homeworks h
		LEFT JOIN submissions s ON s.homework_id = h.id AND s.user_id = $1
		LEFT JOIN group_members gm ON gm.group_id = h.group_id AND gm.user_id = $1
		WHERE (h.user_id = $1 OR gm.user_id = $1) AND h.deleted_at IS NULL
	`
	row := db.Pool.QueryRow(ctx, query, userID)

	var total, submitted, checked, overdue int
	if err := row.Scan(&total, &submitted, &checked, &overdue); err != nil {
		return nil, err
	}

	return map[string]int{
		"total":     total,
		"submitted": submitted,
		"checked":   checked,
		"overdue":   overdue,
	}, nil
}

func (r *PostgresRepo) ListByAuthor(ctx context.Context, authorID uuid.UUID) ([]Homework, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, title, description, type, course_id, module_id,
		       lesson_id, group_id, user_id, author_id, due_at,
		       is_required, created_at, updated_at, deleted_at
		FROM homeworks
		WHERE author_id = $1 AND deleted_at IS NULL
	`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Homework
	for rows.Next() {
		var h Homework
		err := rows.Scan(&h.ID, &h.Title, &h.Description, &h.Type, &h.CourseID, &h.ModuleID,
			&h.LessonID, &h.GroupID, &h.UserID, &h.AuthorID, &h.DueAt,
			&h.IsRequired, &h.CreatedAt, &h.UpdatedAt, &h.DeletedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, nil
}

func (r *PostgresRepo) Filter(ctx context.Context, filters map[string]interface{}, limit, offset int) ([]*Homework, error) {
	baseQuery := `
	SELECT h.id, h.title, h.description, h.type, h.course_id, h.module_id, h.lesson_id, h.group_id, h.user_id, h.author_id,
	       h.due_at, h.is_required, h.created_at, h.updated_at, h.deleted_at
	FROM homeworks h
	LEFT JOIN submissions s ON s.homework_id = h.id AND s.user_id = $1
	WHERE h.deleted_at IS NULL
	`
	args := []interface{}{}
	userID, _ := filters["user_id"].(uuid.UUID)
	args = append(args, userID)

	i := 2
	for key, value := range filters {
		if key == "user_id" || key == "status" {
			continue
		}
		baseQuery += fmt.Sprintf(" AND h.%s = $%d", key, i)
		args = append(args, value)
		i++
	}

	// Фильтр по статусу выполнения
	if status, ok := filters["status"]; ok {
		if status == "done" {
			baseQuery += " AND s.submitted_at IS NOT NULL"
		} else if status == "not_done" {
			baseQuery += " AND s.submitted_at IS NULL"
		}
	}

	baseQuery += fmt.Sprintf(" ORDER BY h.due_at NULLS LAST LIMIT %d OFFSET %d", limit, offset)

	rows, err := db.Pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*Homework
	for rows.Next() {
		var hw Homework
		err := rows.Scan(
			&hw.ID, &hw.Title, &hw.Description, &hw.Type,
			&hw.CourseID, &hw.ModuleID, &hw.LessonID,
			&hw.GroupID, &hw.UserID, &hw.AuthorID,
			&hw.DueAt, &hw.IsRequired, &hw.CreatedAt,
			&hw.UpdatedAt, &hw.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &hw)
	}

	return result, nil
}

// GetByUserAndStatus возвращает домашки, которые:
// - выданы конкретному пользователю (или всем)
// - и выполнены/не выполнены в зависимости от статуса
func (r *PostgresRepo) GetByUserAndStatus(ctx context.Context, userID uuid.UUID, status string) ([]Homework, error) {
	query := `
		SELECT h.id, h.title, h.description, h.type, h.course_id, h.module_id,
		       h.lesson_id, h.group_id, h.user_id, h.author_id, h.due_at,
		       h.is_required, h.created_at, h.updated_at, h.deleted_at
		FROM homeworks h
		LEFT JOIN group_members gm ON h.group_id = gm.group_id AND gm.user_id = $1
		LEFT JOIN submissions s ON h.id = s.homework_id AND s.user_id = $1
		WHERE (h.user_id = $1 OR gm.user_id = $1) AND h.deleted_at IS NULL
	`

	if status == "done" {
		query += " AND s.submitted_at IS NOT NULL"
	} else if status == "not_done" {
		query += " AND s.submitted_at IS NULL"
	}

	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Homework
	for rows.Next() {
		var h Homework
		err := rows.Scan(&h.ID, &h.Title, &h.Description, &h.Type, &h.CourseID, &h.ModuleID,
			&h.LessonID, &h.GroupID, &h.UserID, &h.AuthorID, &h.DueAt,
			&h.IsRequired, &h.CreatedAt, &h.UpdatedAt, &h.DeletedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, h)
	}

	return result, nil
}

func (r *PostgresRepo) ListByLesson(ctx context.Context, lessonID uuid.UUID) ([]Homework, error) {
	rows, err := db.Pool.Query(ctx, `
		SELECT id, title, description, type, course_id, module_id, lesson_id, group_id, user_id, author_id, due_at, is_required, created_at, updated_at, deleted_at
		FROM homeworks
		WHERE lesson_id = $1 AND deleted_at IS NULL
	`, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var homeworks []Homework
	for rows.Next() {
		var h Homework
		err := rows.Scan(&h.ID, &h.Title, &h.Description, &h.Type, &h.CourseID, &h.ModuleID, &h.LessonID, &h.GroupID, &h.UserID, &h.AuthorID, &h.DueAt, &h.IsRequired, &h.CreatedAt, &h.UpdatedAt, &h.DeletedAt)
		if err != nil {
			return nil, err
		}
		homeworks = append(homeworks, h)
	}
	return homeworks, nil
}
