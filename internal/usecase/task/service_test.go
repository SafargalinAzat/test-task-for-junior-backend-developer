package task

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	taskdomain "example.com/taskservice/internal/domain/task"
	mocks "example.com/taskservice/internal/repository/postgres/mocks"
)

func setupTestUsecase(t *testing.T) *Service {
	repo := mocks.NewMockRepository()
	return NewService(repo)
}

func TestService_Create(t *testing.T) {
	usecase := setupTestUsecase(t)

	input := CreateInput{
		Title:       "Test task",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	created, err := usecase.Create(context.Background(), input)
	require.NoError(t, err)
	assert.NotZero(t, created.ID)
	assert.Equal(t, input.Title, created.Title)
	assert.Equal(t, input.Description, created.Description)
	assert.Equal(t, input.Status, created.Status)
	assert.Equal(t, input.Recurrence.Type, created.Recurrence.Type)
	assert.NotZero(t, created.CreatedAt)
	assert.NotZero(t, created.UpdatedAt)
}

func TestService_Create_InvalidInput(t *testing.T) {
	usecase := setupTestUsecase(t)

	input := CreateInput{
		Title:       "",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	_, err := usecase.Create(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")

	input = CreateInput{
		Title:       "Test task",
		Description: "Test description",
		Status:      "invalid",
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	_, err = usecase.Create(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid status")

	input = CreateInput{
		Title:       "Test task",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: "invalid",
		},
	}

	_, err = usecase.Create(context.Background(), input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid recurrence type")
}

func TestService_GetByID(t *testing.T) {
	usecase := setupTestUsecase(t)

	input := CreateInput{
		Title:       "Test task",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	created, err := usecase.Create(context.Background(), input)
	require.NoError(t, err)

	found, err := usecase.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, created.Title, found.Title)
}

func TestService_GetByID_InvalidID(t *testing.T) {
	usecase := setupTestUsecase(t)

	_, err := usecase.GetByID(context.Background(), 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be positive")
}

func TestService_Update(t *testing.T) {
	usecase := setupTestUsecase(t)

	input := CreateInput{
		Title:       "Test task",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	created, err := usecase.Create(context.Background(), input)
	require.NoError(t, err)

	updateInput := UpdateInput{
		Title:       "Updated title",
		Description: "Updated description",
		Status:      taskdomain.StatusInProgress,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	updated, err := usecase.Update(context.Background(), created.ID, updateInput)
	require.NoError(t, err)
	assert.Equal(t, "Updated title", updated.Title)
	assert.Equal(t, "Updated description", updated.Description)
	assert.Equal(t, taskdomain.StatusInProgress, updated.Status)
	assert.NotZero(t, updated.UpdatedAt)
}

func TestService_Update_InvalidID(t *testing.T) {
	usecase := setupTestUsecase(t)

	updateInput := UpdateInput{
		Title:       "Updated title",
		Description: "Updated description",
		Status:      taskdomain.StatusInProgress,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	_, err := usecase.Update(context.Background(), 0, updateInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be positive")
}

func TestService_Update_InvalidInput(t *testing.T) {
	usecase := setupTestUsecase(t)

	input := CreateInput{
		Title:       "Test task",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	created, err := usecase.Create(context.Background(), input)
	require.NoError(t, err)

	updateInput := UpdateInput{
		Title:       "",
		Description: "Updated description",
		Status:      taskdomain.StatusInProgress,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	_, err = usecase.Update(context.Background(), created.ID, updateInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestService_Delete(t *testing.T) {
	usecase := setupTestUsecase(t)

	input := CreateInput{
		Title:       "Test task",
		Description: "Test description",
		Status:      taskdomain.StatusNew,
		Recurrence: taskdomain.Recurrence{
			Type: taskdomain.NotPeriodic,
		},
	}

	created, err := usecase.Create(context.Background(), input)
	require.NoError(t, err)

	err = usecase.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	_, err = usecase.GetByID(context.Background(), created.ID)
	assert.Error(t, err)
}

func TestService_Delete_InvalidID(t *testing.T) {
	usecase := setupTestUsecase(t)

	err := usecase.Delete(context.Background(), 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be positive")
}

func TestService_List(t *testing.T) {
	usecase := setupTestUsecase(t)

	inputs := []CreateInput{
		{
			Title:       "Task 1",
			Description: "Description 1",
			Status:      taskdomain.StatusNew,
			Recurrence: taskdomain.Recurrence{
				Type: taskdomain.NotPeriodic,
			},
		},
		{
			Title:       "Task 2",
			Description: "Description 2",
			Status:      taskdomain.StatusInProgress,
			Recurrence: taskdomain.Recurrence{
				Type: taskdomain.NotPeriodic,
			},
		},
	}

	for _, input := range inputs {
		_, err := usecase.Create(context.Background(), input)
		require.NoError(t, err)
	}

	allTasks, err := usecase.List(context.Background())
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(allTasks), 2)
}
