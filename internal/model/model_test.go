package model

import (
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

func TestTaskProgressRatioAtZero(t *testing.T) {
	taskProgress := NewTaskProgress(uuid.New(), 1, "test task", time.Second*10, time.Second*0)
	if taskProgress.ProgressRatio() != 0 {
		t.Errorf("Expected 0, got %v", taskProgress.ProgressRatio())
	}
}

func TestTaskProgressRatioAtQuarter(t *testing.T) {
	taskProgress := NewTaskProgress(uuid.New(), 1, "test task", time.Second*3, time.Second*12)
	if taskProgress.ProgressRatio() != 0.75 {
		t.Errorf("Expected 0.75, got %v", taskProgress.ProgressRatio())
	}
}

func TestTaskProgressRatioAtHalf(t *testing.T) {
	taskProgress := NewTaskProgress(uuid.New(), 1, "test task", time.Second*5, time.Second*10)
	if taskProgress.ProgressRatio() != 0.5 {
		t.Errorf("Expected 0.5, got %v", taskProgress.ProgressRatio())
	}
}

func TestTaskProgressRatioAtFull(t *testing.T) {
	taskProgress := NewTaskProgress(uuid.New(), 1, "test task", time.Second*0, time.Second*10)
	if taskProgress.ProgressRatio() != 1 {
		t.Errorf("Expected 1, got %v", taskProgress.ProgressRatio())
	}
}

func TestProgressReturnsFirstTaskOnZero(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*10)
	project := NewProject("test project", []Task{firstTask, secondTask}, false)

	progress := project.Progress(time.Second * 0)
	if progress.TaskId() != firstTask.ID() {
		t.Errorf("Expected first task id, got %v", progress.TaskId())
	}
}

func TestProgressReturnsFirstTaskOnFirstTaskEnd(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*10)
	project := NewProject("test project", []Task{firstTask, secondTask}, false)

	progress := project.Progress(firstTask.Duration())
	if progress.TaskId() != firstTask.ID() {
		t.Errorf("Expected first task id, got %v", progress.TaskId())
	}
}

func TestProgressReturnsSecondTaskOnFirstTaskEndPlusOne(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*10)
	project := NewProject("test project", []Task{firstTask, secondTask}, false)

	progress := project.Progress(firstTask.Duration() + time.Second)
	if progress.TaskId() != secondTask.ID() {
		t.Errorf("Expected second task id, got %v", progress.TaskId())
	}
}

func TestProgressReturnsSecondTaskOnFirstTaskEndPlusDuration(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*10)
	project := NewProject("test project", []Task{firstTask, secondTask}, false)

	progress := project.Progress(firstTask.Duration() + secondTask.Duration())
	if progress.TaskId() != secondTask.ID() {
		t.Errorf("Expected second task id, got %v", progress.TaskId())
	}
}

func TestProgressReturnsNilOnFirstTaskEndPlusDurationPlusOne(t *testing.T) {
	firstTask := NewTask("first task", time.Second*10)
	secondTask := NewTask("second task", time.Second*10)
	project := NewProject("test project", []Task{firstTask, secondTask}, false)

	progress := project.Progress(firstTask.Duration() + secondTask.Duration() + time.Second)
	if progress != nil {
		t.Errorf("Expected nil, got %v", progress)
	}
}
