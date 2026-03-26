package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"todos/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTodo(pool *pgxpool.Pool, title string, description string, completed bool, userID string) (*models.Todo, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = "INSERT INTO todos (title, description, completed, user_id) VALUES ($1, $2, $3, $4) RETURNING id, title, description, completed, created_at, updated_at, user_id"

	var todo models.Todo

	var err error = pool.QueryRow(ctx, query, title, description, completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func GetTodos(pool *pgxpool.Pool, userID string) ([]models.Todo, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = "SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE user_id = $1 ORDER BY created_at DESC"

	rows, err := pool.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	var todos []models.Todo

	for rows.Next() {
		var todo models.Todo
		var description sql.NullString
		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		todo.Description = description.String

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func GetTodoByID(pool *pgxpool.Pool, id int, userID string) (*models.Todo, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = "SELECT id, title, description, completed, created_at, updated_at FROM todos WHERE id = $1 AND user_id = $2"

	var todo models.Todo
	var description sql.NullString
	err := pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID)

	if err != nil {
		return nil, err
	}

	todo.Description = description.String

	return &todo, nil
}

func UpdateTodo(pool *pgxpool.Pool, id int, title *string, description *string, completed *bool, userID string) (*models.Todo, error) {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = "UPDATE todos SET title = $1, description = $2, completed = $3, updated_at = NOW() WHERE id = $4 AND user_id = $5 RETURNING id, title, description, completed, created_at, updated_at"

	var todo models.Todo
	var descriptionNull sql.NullString
	err := pool.QueryRow(ctx, query, title, description, completed, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&descriptionNull,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID)

	if err != nil {
		return nil, err
	}

	todo.Description = descriptionNull.String

	return &todo, nil
}

func DeleteTodo(pool *pgxpool.Pool, id int, userID string) error {

	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = "DELETE FROM todos WHERE id = $1 AND user_id = $2"

	var commandTag, err = pool.Exec(ctx, query, id, userID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("Todo with ID %d not found", id)
	}

	return nil
}
