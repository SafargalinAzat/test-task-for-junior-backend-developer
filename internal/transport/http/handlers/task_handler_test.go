package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	taskdomain "example.com/taskservice/internal/domain/task"
	taskusecase "example.com/taskservice/internal/usecase/task"
)

type mockUsecase struct {
	taskusecase.Usecase
	createFunc func(ctx context.Context, input taskusecase.CreateInput) (*taskdomain.Task, error)
	getFunc    func(ctx context.Context, id int64) (*taskdomain.Task, error)
	updateFunc func(ctx context.Context, id int64, input taskusecase.UpdateInput) (*taskdomain.Task, error)
	deleteFunc func(ctx context.Context, id int64) error
	listFunc   func(ctx context.Context) ([]taskdomain.Task, error)
}

func (m *mockUsecase) Create(ctx context.Context, input taskusecase.CreateInput) (*taskdomain.Task, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, input)
	}
	return nil, nil
}

func (m *mockUsecase) GetByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUsecase) Update(ctx context.Context, id int64, input taskusecase.UpdateInput) (*taskdomain.Task, error) {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, id, input)
	}
	return nil, nil
}

func (m *mockUsecase) Delete(ctx context.Context, id int64) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

func (m *mockUsecase) List(ctx context.Context) ([]taskdomain.Task, error) {
	if m.listFunc != nil {
		return m.listFunc(ctx)
	}
	return nil, nil
}

func setupTestHandler(t *testing.T) (*TaskHandler, *mockUsecase) {
	mock := &mockUsecase{}
	handler := NewTaskHandler(mock)
	return handler, mock
}

func TestTaskHandler_Create(t *testing.T) {
	handler, mock := setupTestHandler(t)

	mock.createFunc = func(ctx context.Context, input taskusecase.CreateInput) (*taskdomain.Task, error) {
		return &taskdomain.Task{
			ID:          1,
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Recurrence:  input.Recurrence,
		}, nil
	}

	reqBody := taskMutationDTO{
		Title:       "Test task",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Test task", response["title"])
	assert.Equal(t, "Test description", response["description"])
}

func TestTaskHandler_Create_InvalidJSON(t *testing.T) {
	handler, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTaskHandler_GetByID(t *testing.T) {
	handler, mock := setupTestHandler(t)

	mock.getFunc = func(ctx context.Context, id int64) (*taskdomain.Task, error) {
		return &taskdomain.Task{
			ID:          id,
			Title:       "Test task",
			Description: "Test description",
			Status:      taskdomain.StatusNew,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Recurrence: taskdomain.Recurrence{
				Type: taskdomain.NotPeriodic,
			},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Test task", response["title"])
}

func TestTaskHandler_GetByID_InvalidID(t *testing.T) {
	handler, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/tasks/invalid", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTaskHandler_Update(t *testing.T) {
	handler, mock := setupTestHandler(t)

	mock.updateFunc = func(ctx context.Context, id int64, input taskusecase.UpdateInput) (*taskdomain.Task, error) {
		return &taskdomain.Task{
			ID:          id,
			Title:       input.Title,
			Description: input.Description,
			Status:      input.Status,
			UpdatedAt:   time.Now(),
			Recurrence:  input.Recurrence,
		}, nil
	}

	reqBody := taskMutationDTO{
		Title:       "Updated title",
		Description: "Updated description",
		Status:      taskdomain.StatusInProgress,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}
	reqBytes, err := json.Marshal(reqBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer(reqBytes))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Updated title", response["title"])
}

func TestTaskHandler_Update_InvalidJSON(t *testing.T) {
	handler, _ := setupTestHandler(t)

	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewBuffer([]byte("{invalid}")))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestTaskHandler_Delete(t *testing.T) {
	handler, mock := setupTestHandler(t)

	mock.deleteFunc = func(ctx context.Context, id int64) error {
		return nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestTaskHandler_List(t *testing.T) {
	handler, mock := setupTestHandler(t)

	mock.listFunc = func(ctx context.Context) ([]taskdomain.Task, error) {
		return []taskdomain.Task{
			{
				ID:          1,
				Title:       "Task 1",
				Description: "Description 1",
				Status:      taskdomain.StatusNew,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Recurrence: taskdomain.Recurrence{
					Type: taskdomain.NotPeriodic,
				},
			},
			{
				ID:          2,
				Title:       "Task 2",
				Description: "Description 2",
				Status:      taskdomain.StatusInProgress,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Recurrence: taskdomain.Recurrence{
					Type: taskdomain.NotPeriodic,
				},
			},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 2)
}
