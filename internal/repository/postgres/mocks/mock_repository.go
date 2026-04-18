package postgres

import (
	"context"
	"errors"
	"sync"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type MockRepository struct {
	tasks    map[int64]*taskdomain.Task
	nextID   int64
	mu       sync.RWMutex
	failNext bool
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		tasks: make(map[int64]*taskdomain.Task),
	}
}

func (m *MockRepository) Create(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.failNext {
		m.failNext = false
		return nil, errors.New("mock database error")
	}

	m.nextID++
	task.ID = m.nextID
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	m.tasks[m.nextID] = task

	return task, nil
}

func (m *MockRepository) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.failNext {
		m.failNext = false
		return nil, errors.New("mock database error")
	}

	task, exists := m.tasks[id]
	if !exists {
		return nil, taskdomain.ErrNotFound
	}

	return task, nil
}

func (m *MockRepository) Update(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.failNext {
		m.failNext = false
		return nil, errors.New("mock database error")
	}

	existing, exists := m.tasks[task.ID]
	if !exists {
		return nil, taskdomain.ErrNotFound
	}

	existing.Title = task.Title
	existing.Description = task.Description
	existing.Status = task.Status
	existing.UpdatedAt = time.Now()
	existing.Recurrence = task.Recurrence

	return existing, nil
}

func (m *MockRepository) Delete(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.failNext {
		m.failNext = false
		return errors.New("mock database error")
	}

	if _, exists := m.tasks[id]; !exists {
		return taskdomain.ErrNotFound
	}

	delete(m.tasks, id)
	return nil
}

func (m *MockRepository) List(ctx context.Context) ([]taskdomain.Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.failNext {
		m.failNext = false
		return nil, errors.New("mock database error")
	}

	tasks := make([]taskdomain.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (m *MockRepository) FailNext() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failNext = true
}
