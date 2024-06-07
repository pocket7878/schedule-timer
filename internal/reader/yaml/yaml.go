package yaml

import (
	"schedule-timer/internal/model"
	"time"

	"gopkg.in/yaml.v3"
)

type taskYAML struct {
	Name              string
	DurationInSeconds int `yaml:"duration"`
}

type projectYAML struct {
	Name  string
	Tasks []taskYAML
}

// Reader is an interface for reading a project from a file.
func Read(data []byte) (*model.Project, error) {
	var projectYAML projectYAML
	err := yaml.Unmarshal(data, &projectYAML)

	if err != nil {
		return nil, err
	}

	tasks := make([]model.Task, len(projectYAML.Tasks))
	for i, task := range projectYAML.Tasks {
		tasks[i] = model.NewTask(task.Name, time.Second*time.Duration(task.DurationInSeconds))
	}

	project := model.NewProject(projectYAML.Name, tasks)
	return &project, nil
}
