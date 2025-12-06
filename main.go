package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	startTime        time.Time
	lastBeatTime     time.Time
	beatsSinceStart  int
	currentBPM       float64
	secondsSinceBeat int
}

func InitialModel() model {
	return model{}
}

type secondPassedMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return secondPassedMsg{}
}

func waitSecond() tea.Cmd {
	return tick
}

func (m model) Init() tea.Cmd {
	return waitSecond()
}

func updateBPM(m model) model {
	now := time.Now()
	m.secondsSinceBeat = 0

	if m.beatsSinceStart == 0 {
		m.beatsSinceStart++
		m.startTime = now
		m.lastBeatTime = now
		m.currentBPM = 0
		return m
	}
	m.lastBeatTime = now
	m.beatsSinceStart++

	elapsedMinutes := now.Sub(m.startTime).Minutes()
	if elapsedMinutes > 0 {
		m.currentBPM = float64(m.beatsSinceStart-1) / elapsedMinutes
	}

	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "r":
			return InitialModel(), waitSecond()
		case " ":
			m = updateBPM(m)
			return m, nil
		}

	case secondPassedMsg:
		m.secondsSinceBeat++
		if m.secondsSinceBeat > 5 {
			m = InitialModel()
		}
		return m, waitSecond()
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("space: tap tempo | r: reset (or 5 seconds without keypress) | q: quit\n\nCurrent BPM: %.1f\n", m.currentBPM)
}

func main() {
	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
