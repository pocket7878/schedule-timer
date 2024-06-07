package yaml

import (
	"testing"
	"time"
)

func TestReadWithoutRepeat(t *testing.T) {
	data := []byte(`
name: test project
tasks:
  - name: first task
    duration: 10
  - name: second task
    duration: 10
`)
	project, err := Read(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if project.Name() != "test project" {
		t.Errorf("Expected test project, got %v", project.Name())
	}

	if project.Repeat() {
		t.Errorf("Expected false, got %v", project.Repeat())
	}

	tasks := project.Tasks()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %v", len(tasks))
	}

	if tasks[0].Name() != "first task" {
		t.Errorf("Expected first task, got %v", tasks[0].Name())
	}

	if tasks[0].Duration() != time.Second*10 {
		t.Errorf("Expected 10, got %v", tasks[0].Duration())
	}

	if tasks[1].Name() != "second task" {
		t.Errorf("Expected second task, got %v", tasks[1].Name())
	}

	if tasks[1].Duration() != time.Second*10 {
		t.Errorf("Expected 10, got %v", tasks[1].Duration())
	}
}

func TestReadWithRepeat(t *testing.T) {
	data := []byte(`
name: test project
repeat: true
tasks:
  - name: first task
    duration: 10
`)
	project, err := Read(data)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if project.Name() != "test project" {
		t.Errorf("Expected test project, got %v", project.Name())
	}

	if project.Repeat() != true {
		t.Errorf("Expected true, got %v", project.Repeat())
	}
}
