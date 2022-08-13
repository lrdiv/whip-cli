package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var platforms = []string{
	"spotify",
	"itunes",
	"youtube",
	"tidal",
	"amazonMusic",
	"pandora",
	"deezer",
	"audiomack",
	"qobuz",
}

const (
	getSource State = 0
	choosePlatform State = 1
	fetching State = 2
	done State = 3
	hasError State = 4
)

type State int

type model struct {
	state State
	source textinput.Model
	cursor int
	platform int
	url string
	errMsg error
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Oh no! An error! :( %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "https://open.spotify.com/track/Sfsdj3924hjd"
	ti.Focus()
	return model{
		state: getSource,
		source: ti,
		cursor: 0,
		platform: 0,
		url: "",
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) updatePlatform(msg tea.Msg) (tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(platforms)-1 {
				m.cursor++
			}
		case "enter":
			m.state = fetching
			m.platform = m.cursor
		}
	}

	return nil
}

func (m *model) updateSource(msg tea.Msg) (tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.state = choosePlatform
			return nil

		case tea.KeyCtrlC, tea.KeyEsc:
			return tea.Quit
		}

	case error:
		m.errMsg = msg
		return nil
	}

	m.source, cmd = m.source.Update(msg)
	return cmd
}

func (m *model) getUrl() (tea.Cmd) {
	var jsonData = []byte(`{
			"url": {{.m.source}}
		}`)
	client := http.Client{}
	req, _ := http.NewRequest("POST", "https://songwhip.com/", bytes.NewBuffer(jsonData))
	req.Header = http.Header{
		"Content-Type": {"application/json"},
	}
	res, err := client.Do(req)
	if err != nil {
		m.state = hasError
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
			m.state = hasError
		}
		m.url = string(bodyBytes)
		m.state = done
	}
	
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var ret tea.Cmd

	var keyMsg, isKeyMsg = msg.(tea.KeyMsg)
	if (isKeyMsg && keyMsg.Type == tea.KeyCtrlC || keyMsg.Type == tea.KeyEsc) {
		return m, tea.Quit
	}

	switch m.state {
	case getSource:
		ret = m.updateSource(msg)
	case choosePlatform:
		ret = m.updatePlatform(msg)
	case fetching:
		m.getUrl()
	}

	return m, ret
}

func (m model) View() string {
	s := ""

	if (m.state == getSource) {
		return fmt.Sprintf("Enter a track or album URL from any supported platform...\n\n%s\n\n%s",
			m.source.View(),
			"(esc to quit)",
		)
	}

	if (m.state == choosePlatform) {
		s += "Which platform do you want a link for?\n\n"
		for i, platform := range platforms {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, strings.Title(platform))
		}

		s += "\nPress q to quit.\n"
	}

	if (m.state == fetching) {
		s += "Fetching links from Songwhip..."
	}

	if (m.state == done) {
		s += m.url
	}

	if (m.state == hasError) {
		s += "aw fudge we got error."
	}

	return s
}
