package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	const query = `
		INSERT INTO tasks (title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, description, status, created_at, updated_at
	`

	row := tx.QueryRow(ctx, query, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt)
	created, err := scanTaskWithoutRec(row)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	const recurrenceQuery = `
		INSERT INTO task_recurrences (task_id, type, every_n_days, monthly_date, dates, even_odd)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(ctx, recurrenceQuery, created.ID, task.Recurrence.Type, task.Recurrence.EveryNDays, task.Recurrence.MonthlyDate, task.Recurrence.Dates, task.Recurrence.EvenOdd)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	created.Recurrence = task.Recurrence

	return created, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	const query = `
		SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at,
		       r.type, r.every_n_days, r.monthly_date, r.dates, r.even_odd
		FROM tasks t
		LEFT JOIN task_recurrences r ON r.task_id = t.id
		WHERE t.id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	found, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}

		return nil, err
	}

	return found, nil
}

func (r *Repository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	const query = `
		UPDATE tasks
		SET title = $1,
			description = $2,
			status = $3,
			updated_at = $4
		WHERE id = $5
	`

	_, err = tx.Exec(ctx, query, task.Title, task.Description, task.Status, task.UpdatedAt, task.ID)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	const recurrenceQuery = `
        INSERT INTO task_recurrences (task_id, type, every_n_days, monthly_date, dates, even_odd)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (task_id) DO UPDATE SET
            type = EXCLUDED.type,
            every_n_days = EXCLUDED.every_n_days,
            monthly_date = EXCLUDED.monthly_date,
            dates = EXCLUDED.dates,
            even_odd = EXCLUDED.even_odd
    `
	_, err = tx.Exec(ctx, recurrenceQuery, task.ID, task.Recurrence.Type, task.Recurrence.EveryNDays, task.Recurrence.MonthlyDate, task.Recurrence.Dates, task.Recurrence.EvenOdd)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	const selectQuery = `
        SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at,
               r.type, r.every_n_days, r.monthly_date, r.dates, r.even_odd
        FROM tasks t
        LEFT JOIN task_recurrences r ON r.task_id = t.id
        WHERE t.id = $1
    `

	row := tx.QueryRow(ctx, selectQuery, task.ID)
	updated, err := scanTask(row)
	if err != nil {
		tx.Rollback(ctx)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrNotFound
		}
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrNotFound
	}

	return nil
}

func (r *Repository) List(ctx context.Context) ([]taskdomain.Task, error) {
	const query = `
        SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at,
               r.type, r.every_n_days, r.monthly_date, r.dates, r.even_odd
        FROM tasks t
        LEFT JOIN task_recurrences r ON r.task_id = t.id
        ORDER BY t.id DESC
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]taskdomain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

type taskScanner interface {
	Scan(dest ...any) error
}

func scanTask(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Recurrence.Type,
		&task.Recurrence.EveryNDays,
		&task.Recurrence.MonthlyDate,
		&task.Recurrence.Dates,
		&task.Recurrence.EvenOdd,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)

	return &task, nil
}

func scanTaskWithoutRec(scanner taskScanner) (*taskdomain.Task, error) {
	var (
		task   taskdomain.Task
		status string
	)

	if err := scanner.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = taskdomain.Status(status)

	return &task, nil
}
