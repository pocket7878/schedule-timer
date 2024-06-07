package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	id       uuid.UUID
	name     string
	duration time.Duration
}

func NewTask(name string, duration time.Duration) Task {
	taskId := uuid.New()
	return Task{id: taskId, name: name, duration: duration}
}

func (t Task) ID() uuid.UUID {
	return t.id
}

func (t Task) Name() string {
	return t.name
}

func (t Task) Duration() time.Duration {
	return t.duration
}

type Project struct {
	id     uuid.UUID
	name   string
	tasks  []Task
	repeat bool
}

func NewProject(name string, tasks []Task, repeat bool) Project {
	projectId := uuid.New()
	return Project{id: projectId, name: name, tasks: tasks, repeat: repeat}
}

func (p Project) ID() uuid.UUID {
	return p.id
}

func (p Project) Name() string {
	return p.name
}

func (p Project) Repeat() bool {
	return p.repeat
}

func (p Project) Tasks() []Task {
	return p.tasks
}

func (p *Project) AddTask(task Task) {
	p.tasks = append(p.tasks, task)
}

func (p Project) Duration() time.Duration {
	var total time.Duration
	for _, task := range p.tasks {
		total += task.Duration()
	}
	return total
}

type TaskProgress struct {
	taskId        uuid.UUID
	taskIndex     int
	name          string
	restDuration  time.Duration
	totalDuration time.Duration
}

func NewTaskProgress(taskId uuid.UUID, taskIndex int, name string, restDuration time.Duration, totalDuration time.Duration) TaskProgress {
	return TaskProgress{taskId: taskId, taskIndex: taskIndex, name: name, restDuration: restDuration, totalDuration: totalDuration}
}

func (tp TaskProgress) TaskId() uuid.UUID {
	return tp.taskId
}

func (tp TaskProgress) Name() string {
	return tp.name
}

func (tp TaskProgress) TaskIndex() int {
	return tp.taskIndex
}

func (tp TaskProgress) RestDuration() time.Duration {
	return tp.restDuration
}

func (tp TaskProgress) TotalDuration() time.Duration {
	return tp.totalDuration
}

func (tp TaskProgress) ProgressRatio() float64 {
	if tp.totalDuration == 0 {
		return 0
	}
	if tp.restDuration == 0 {
		return 1
	}

	return 1.0 - float64(tp.restDuration)/float64(tp.totalDuration)
}

func (p Project) Progress(elapsedDuration time.Duration) *TaskProgress {
	// Find active task
	var activeTask *Task
	var activeTaskIndex int
	var totalDuration time.Duration
	for i, task := range p.tasks {
		totalDuration += task.Duration()
		if elapsedDuration <= totalDuration {
			activeTask = &task
			activeTaskIndex = i
			break
		}
	}

	if activeTask == nil {
		return nil
	}

	restDuration := totalDuration - elapsedDuration
	activeTaskProgress := NewTaskProgress(activeTask.ID(), activeTaskIndex, activeTask.Name(), restDuration, activeTask.Duration())

	return &activeTaskProgress
}
