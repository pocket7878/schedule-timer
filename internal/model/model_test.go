package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestAddTask(t *testing.T) {
	project := NewProject("test project", []Task{}, false)
	task := NewTask("test task", time.Second*10)

	project.AddTask(task)
	if len(project.Tasks()) != 1 {
		t.Errorf("Expected 1 task, got %v", len(project.Tasks()))
	}
}

func TestNewProject(t *testing.T) {
	project := NewProject("test project", []Task{}, false)
	if project.Name() != "test project" {
		t.Errorf("Expected test project, got %v", project.Name())
	}
}

func TestProjectDuration(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*20)
	project := NewProject("test project", []Task{firstTask, secondTask}, false)

	if project.Duration() != firstTask.Duration()+secondTask.Duration() {
		t.Errorf("Expected %v, got %v", firstTask.Duration()+secondTask.Duration(), project.Duration())
	}
}

func TestTaskProgressRatio(t *testing.T) {
	t.Parallel()
	tests := []struct {
		totalDuration time.Duration
		elapsed       time.Duration
		expected      float64
	}{
		{time.Second * 12, time.Second * 0, 0},
		{time.Second * 12, time.Second * 3, 0.25},
		{time.Second * 12, time.Second * 6, 0.5},
		{time.Second * 12, time.Second * 9, 0.75},
		{time.Second * 12, time.Second * 12, 1},
	}
	for _, test := range tests {
		test := test
		t.Run("", func(t *testing.T) {
			t.Parallel()
			taskProgress := NewTaskProgress(uuid.New(), 1, "test task", test.elapsed, test.totalDuration)
			if taskProgress.ProgressRatio() != test.expected {
				t.Errorf("Expected %v when %s / %s, got %v", test.expected, test.elapsed.String(), test.totalDuration.String(), taskProgress.ProgressRatio())
			}
		})
	}
}

func TestNonRepeatProjectProgress(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*10)
	project := NewProject("test project", []Task{firstTask, secondTask}, false)

	tests := []struct {
		elapsed time.Duration

		expectedTaskID          *uuid.UUID
		expectedTaskIndex       int
		expectedName            string
		expectedElapsedDuration time.Duration
		expectedTotalDuration   time.Duration
	}{
		{time.Second * 0, makePointer(firstTask.ID()), 0, firstTask.Name(), time.Second * 0, firstTask.Duration()},
		{time.Second * 5, makePointer(firstTask.ID()), 0, firstTask.Name(), time.Second * 5, firstTask.Duration()},
		{time.Second * 10, makePointer(secondTask.ID()), 1, secondTask.Name(), time.Second * 0, secondTask.Duration()},
		{time.Second * 15, makePointer(secondTask.ID()), 1, secondTask.Name(), time.Second * 5, secondTask.Duration()},
		{time.Second * 20, nil, -1, "", time.Second * 0, time.Second * 0},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Elapsed %s", test.elapsed.String()), func(t *testing.T) {
			progress := project.Progress(test.elapsed)
			if progress == nil {
				if test.expectedTaskID != nil {
					t.Errorf("Expected %v, got nil", test.expectedTaskID)
				}
				return
			}
			if progress.TaskId() != *test.expectedTaskID {
				t.Errorf("Expected %v, got %v", test.expectedTaskID, progress.TaskId())
			}
			if progress.TaskIndex() != test.expectedTaskIndex {
				t.Errorf("Expected %v, got %v", test.expectedTaskIndex, progress.TaskIndex())
			}
			if progress.Name() != test.expectedName {
				t.Errorf("Expected %v, got %v", test.expectedName, progress.Name())
			}
			if progress.ElapsedDuration() != test.expectedElapsedDuration {
				t.Errorf("Expected %v, got %v", test.expectedElapsedDuration, progress.ElapsedDuration())
			}
			if progress.TotalDuration() != test.expectedTotalDuration {
				t.Errorf("Expected %v, got %v", test.expectedTotalDuration, progress.TotalDuration())
			}
		})
	}
}

func TestRepeatProjectProgress(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*10)
	project := NewProject("test project", []Task{firstTask, secondTask}, true)

	tests := []struct {
		elapsed time.Duration

		expectedTaskID          *uuid.UUID
		expectedTaskIndex       int
		expectedName            string
		expectedElapsedDuration time.Duration
		expectedTotalDuration   time.Duration
	}{
		{time.Second * 0, makePointer(firstTask.ID()), 0, firstTask.Name(), time.Second * 0, firstTask.Duration()},
		{time.Second * 5, makePointer(firstTask.ID()), 0, firstTask.Name(), time.Second * 5, firstTask.Duration()},
		{time.Second * 10, makePointer(secondTask.ID()), 1, secondTask.Name(), time.Second * 0, secondTask.Duration()},
		{time.Second * 15, makePointer(secondTask.ID()), 1, secondTask.Name(), time.Second * 5, secondTask.Duration()},
		{time.Second * 20, makePointer(firstTask.ID()), 0, firstTask.Name(), time.Second * 0, firstTask.Duration()},
		{time.Second * 25, makePointer(firstTask.ID()), 0, firstTask.Name(), time.Second * 5, firstTask.Duration()},
		{time.Second * 30, makePointer(secondTask.ID()), 1, secondTask.Name(), time.Second * 0, secondTask.Duration()},
		{time.Second * 35, makePointer(secondTask.ID()), 1, secondTask.Name(), time.Second * 5, secondTask.Duration()},
	}
	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("Elapsed %s", test.elapsed.String()), func(t *testing.T) {
			progress := project.Progress(test.elapsed)
			if progress == nil {
				if test.expectedTaskID != nil {
					t.Errorf("Expected %v, got nil", test.expectedTaskID)
				}
				return
			}
			if progress.TaskId() != *test.expectedTaskID {
				t.Errorf("Expected %v, got %v", test.expectedTaskID, progress.TaskId())
			}
			if progress.TaskIndex() != test.expectedTaskIndex {
				t.Errorf("Expected %v, got %v", test.expectedTaskIndex, progress.TaskIndex())
			}
			if progress.Name() != test.expectedName {
				t.Errorf("Expected %v, got %v", test.expectedName, progress.Name())
			}
			if progress.ElapsedDuration() != test.expectedElapsedDuration {
				t.Errorf("Expected %v, got %v", test.expectedElapsedDuration, progress.ElapsedDuration())
			}
			if progress.TotalDuration() != test.expectedTotalDuration {
				t.Errorf("Expected %v, got %v", test.expectedTotalDuration, progress.TotalDuration())
			}
		})
	}
}

func makePointer[T any](t T) *T {
	return &t
}
