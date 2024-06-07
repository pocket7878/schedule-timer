package tui

import (
	smodel "schedule-timer/internal/model"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

func Run(project smodel.Project) error {
	m := initialModel(project)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}

	return nil
}

type model struct {
	progress  progress.Model
	taskTable table.Model

	elapsedDuration    time.Duration
	project            smodel.Project
	latestTaskProgress *smodel.TaskProgress
	screenWidth        int
	screenHeight       int
	done               bool
}

const (
	// 1秒ごとに進行状況を更新する
	tickInterval = time.Second * 1
)

type tickMsg struct{}

func initialModel(project smodel.Project) model {
	latestTaskProgress := project.Progress(time.Second * 0)
	taskTable := buildTaskTables(project, latestTaskProgress.TaskIndex())
	return model{
		progress:           progress.New(progress.WithGradient("#FF0000", "#00FF00"), progress.WithoutPercentage()),
		project:            project,
		elapsedDuration:    time.Second * 0,
		screenWidth:        0,
		screenHeight:       0,
		latestTaskProgress: latestTaskProgress,
		taskTable:          taskTable,
		done:               false,
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.String() == "q" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
		m.progress.Width = msg.Width
		return m, nil
	case tickMsg:
		m.elapsedDuration += tickInterval
		taskProgress := m.project.Progress(m.elapsedDuration)
		// まだ進行中の場合は時間を進める
		if taskProgress != nil {
			m.latestTaskProgress = taskProgress
			m.taskTable.SetCursor(taskProgress.TaskIndex())
			return m, tickCmd()
		} else {
			m.done = true
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	taskProgress := m.latestTaskProgress
	if taskProgress == nil {
		return "Done"
	}

	screen := ""
	screen += "Project: " + m.project.Name() + "\n"
	screen += "Task: " + taskProgress.Name() + "\n"
	progressRatio := taskProgress.ProgressRatio()
	// 表示するときに、バーが減って行くように表示する
	screen += m.progress.ViewAs(1.0-progressRatio) + "\n"
	screen += m.progress.ViewAs(1.0-progressRatio) + "\n"
	screen += m.progress.ViewAs(1.0-progressRatio) + "\n"
	restDuration := taskProgress.TotalDuration() - taskProgress.ElapsedDuration()
	screen += restDuration.String() + " / " + taskProgress.TotalDuration().String() + "\n"
	screen += strings.Repeat("-", m.screenWidth) + "\n"
	screen += m.taskTable.View() + "\n"

	return screen
}

func buildTaskTables(project smodel.Project, activeTaskIndex int) table.Model {
	rows := make([]table.Row, 0, len(project.Tasks()))
	for i, task := range project.Tasks() {
		rows = append(rows, table.Row{strconv.Itoa(i), task.Name(), task.Duration().String()})
	}
	var maxIDWidth int = 3
	var maxNameWidth int = 4
	var maxDurationWidth int = 8
	for _, row := range rows {
		if len(row[0]) > maxIDWidth {
			maxIDWidth = len(row[0])
		}
		if len(row[1]) > maxNameWidth {
			maxNameWidth = len(row[1])
		}
		if len(row[2]) > maxDurationWidth {
			maxDurationWidth = len(row[2])
		}
	}

	columns := []table.Column{
		{Title: "ID", Width: maxIDWidth},
		{Title: "Name", Width: maxNameWidth},
		{Title: "Duration", Width: maxDurationWidth},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)
	t.SetCursor(activeTaskIndex)

	return t
}

var _ tea.Model = model{}

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}
