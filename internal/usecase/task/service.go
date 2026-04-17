package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type Service struct {
	repo Repository
	now  func() time.Time
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
		now:  func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error) {
	normalized, err := validateCreateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		Recurrence:  normalized.Recurrence,
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	normalized, err := validateUpdateInput(input)
	if err != nil {
		return nil, err
	}

	model := &taskdomain.Task{
		ID:          id,
		Title:       normalized.Title,
		Description: normalized.Description,
		Status:      normalized.Status,
		UpdatedAt:   s.now(),
		Recurrence:  normalized.Recurrence,
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: id must be positive", ErrInvalidInput)
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.repo.List(ctx)
}

func validateCreateInput(input CreateInput) (CreateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return CreateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if input.Status == "" {
		input.Status = taskdomain.StatusNew
	}

	if !input.Status.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	if !input.Recurrence.Type.Valid() {
		return CreateInput{}, fmt.Errorf("%w: invalid recurrence type", ErrInvalidInput)
	}

	switch input.Recurrence.Type {
	case taskdomain.Daily:
		if !(input.Recurrence.EveryNDays > 0) {
			return CreateInput{}, fmt.Errorf("%w: number of days must be positive", ErrInvalidInput)
		}
	case taskdomain.Monthly:
		if !(input.Recurrence.MonthlyDate >= 1 && input.Recurrence.MonthlyDate <= 30) {
			return CreateInput{}, fmt.Errorf("%w: date should be between 1 and 30", ErrInvalidInput)
		}
	case taskdomain.ExactDates:
		if !(len(input.Recurrence.Dates) >= 1 && len(input.Recurrence.Dates) <= 1000) {
			return CreateInput{}, fmt.Errorf("%w: you should specify at least 1 date and at most 1000", ErrInvalidInput)
		}
	case taskdomain.EvenOddDates:
		if !input.Recurrence.EvenOdd.Valid() {
			return CreateInput{}, fmt.Errorf("%w: even_odd_dates param is incorrect", ErrInvalidInput)
		}
	case taskdomain.NotPeriodic:
	default:
		return CreateInput{}, fmt.Errorf("%w: unknown recurrence_type", ErrInvalidInput)
	}

	return input, nil
}

func validateUpdateInput(input UpdateInput) (UpdateInput, error) {
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)

	if input.Title == "" {
		return UpdateInput{}, fmt.Errorf("%w: title is required", ErrInvalidInput)
	}

	if !input.Status.Valid() {
		return UpdateInput{}, fmt.Errorf("%w: invalid status", ErrInvalidInput)
	}

	switch input.Recurrence.Type {
	case taskdomain.Daily:
		if !(input.Recurrence.EveryNDays > 0) {
			return UpdateInput{}, fmt.Errorf("%w: number of days must be positive", ErrInvalidInput)
		}
	case taskdomain.Monthly:
		if !(input.Recurrence.MonthlyDate >= 1 && input.Recurrence.MonthlyDate <= 30) {
			return UpdateInput{}, fmt.Errorf("%w: date should be between 1 and 30", ErrInvalidInput)
		}
	case taskdomain.ExactDates:
		if !(len(input.Recurrence.Dates) >= 1 && len(input.Recurrence.Dates) <= 1000) {
			return UpdateInput{}, fmt.Errorf("%w: you should specify at least 1 date and at most 1000", ErrInvalidInput)
		}
	case taskdomain.EvenOddDates:
		if !input.Recurrence.EvenOdd.Valid() {
			return UpdateInput{}, fmt.Errorf("%w: even_odd_dates is incorrect", ErrInvalidInput)
		}
	case taskdomain.NotPeriodic:
	default:
		return UpdateInput{}, fmt.Errorf("%w: unknown recurrence_type", ErrInvalidInput)
	}

	return input, nil
}
