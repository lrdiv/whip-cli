package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gocolly/colly"
)

// Enum for managing TUI state
type state int64

const (
	GetOriginalUrl state = iota
	ChoosePlatform
	FetchingSongwhip
	CrawlingSongwhip
	Done
	HasError
)

// Message types to be used in bubbletea commands
type errorMsg struct{}
type songwhipReadyMsg struct{}
type songwhipDoneMsg struct {
	url string
}
type songwhipCrawlMsg struct {
	url string
}

type platform struct {
	Slug     string
	Title    string
	HelpText string
}

var platforms [10]platform = [10]platform{
	{
		Slug:     "songwhip",
		Title:    "Songwhip",
		HelpText: "Get a Songwhip URL with links to all available platforms.",
	}, {
		Slug:  "spotify",
		Title: "Spotify",
	}, {
		Slug:  "itunes",
		Title: "Apple Music",
	}, {
		Slug:  "youtube",
		Title: "YouTube Music",
	}, {
		Slug:  "tidal",
		Title: "Tidal",
	}, {
		Slug:  "amazonMusic",
		Title: "Amazon Music",
	}, {
		Slug:  "pandora",
		Title: "Pandora",
	}, {
		Slug:  "deezer",
		Title: "Deezer",
	}, {
		Slug:  "audiomack",
		Title: "AudioMack",
	}, {
		Slug:  "qobuz",
		Title: "Qobuz",
	},
}

type model struct {
	Log            *os.File
	OriginalUrl    textinput.Model
	Platform       platform
	PlatformCursor int
	PlatformUrl    string
	SongwhipData   songwhipResponse
	Spinner        spinner.Model
	State          state
}

var p *tea.Program

func main() {
	model := initialModel()
	p = tea.NewProgram(model)

	if err := p.Start(); err != nil {
		fmt.Printf("Oh no! An error! :( %v", err)
		os.Exit(1)
	}
}

func initialModel() model {
	return model{
		Log:            openLogFile(),
		OriginalUrl:    makeUrlTextinput(),
		PlatformCursor: 0,
		PlatformUrl:    "",
		Spinner:        makeSpinner(),
		State:          GetOriginalUrl,
	}
}

func makeUrlTextinput() textinput.Model {
	input := textinput.New()
	input.Placeholder = "https://open.spotify.com/track/4uLU6hMCjMI75M1A2tKUQC"
	input.Focus()
	return input
}

func makeSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	return s
}

func openLogFile() *os.File {
	path := "./tmp/whip.log"
	os.Truncate(path, 0)
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil || logFile == nil {
		return nil
	}
	return logFile
}

func (m model) writeToLog(text string) {
	logger := log.New(m.Log, "[info]", log.LstdFlags|log.Lmicroseconds)
	logger.Println(text)
}

func quit(log *os.File) tea.Cmd {
	if log != nil {
		defer log.Close()
	}
	return tea.Quit
}

func platformSelectionView(platformCursor int) string {
	var sb strings.Builder
	sb.WriteString("Which platform do you want a link for?\n\n")
	for i, platform := range platforms {
		cursor := " "
		if platformCursor == i {
			cursor = ">"
			style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
			sb.WriteString(style.Render(fmt.Sprintf("%s %s", cursor, platform.Title)))
			if platform.HelpText != "" {
				sb.WriteString(fmt.Sprintf(": %s", platform.HelpText))
			}
		} else {
			sb.WriteString(fmt.Sprintf("%s %s", cursor, platform.Title))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("\nPress ctrl+c or esc to quit.\n")
	return sb.String()
}

func getSongwhipData(url string) {
	var songwhipData songwhipResponse
	var jsonData = bytes.NewBuffer([]byte(fmt.Sprintf(`{"url": "%s"}`, url)))
	res, err := http.Post("https://songwhip.com/", "application/json", jsonData)

	if err != nil {
		p.Send(errorMsg{})
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			p.Send(errorMsg{})
		}

		jsonErr := json.Unmarshal(bodyBytes, &songwhipData)
		if jsonErr != nil {
			p.Send(errorMsg{})
		}
		p.Send(songwhipDoneMsg{
			url: songwhipData.Url,
		})
	}

}

func crawlSongwhip(url string, platform string) {
	var platformUrl string

	if platform == "songwhip" {
		p.Send(songwhipCrawlMsg{url: url})
	}

	c := colly.NewCollector(colly.AllowedDomains("songwhip.com"))
	selector := fmt.Sprintf("a[data-testid=\"ServiceButton %s itemLinkButton %sItemLinkButton\"]", platform, platform)

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		platformUrl = e.Attr("href")
	})

	c.Visit(url)
	p.Send(songwhipCrawlMsg{
		url: platformUrl,
	})
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Keypress messages can come from either the GetOriginalUrl or ChoosePlatform states so we handle those before
	// focusing on custom message structs
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.State == Done {
			return m, quit(m.Log)
		}
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, quit(m.Log)
		case tea.KeyEnter:
			switch m.State {
			case GetOriginalUrl:
				m.State = ChoosePlatform
				return m, nil
			case ChoosePlatform:
				m.Platform = platforms[m.PlatformCursor]
				return m.Update(songwhipReadyMsg{})
			}
		case tea.KeyUp:
			if m.State == ChoosePlatform && m.PlatformCursor > 0 {
				m.PlatformCursor--
			}
		case tea.KeyDown:
			if m.State == ChoosePlatform && m.PlatformCursor < len(platforms)-1 {
				m.PlatformCursor++
			}
		default:
			if m.State == GetOriginalUrl {
				originalUrl, cmd := m.OriginalUrl.Update(msg)
				m.OriginalUrl = originalUrl
				return m, cmd
			}
		}
	case songwhipReadyMsg:
		m.State = FetchingSongwhip
		go getSongwhipData(m.OriginalUrl.Value())
		return m, m.Spinner.Tick
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case songwhipDoneMsg:
		m.State = CrawlingSongwhip
		go crawlSongwhip(msg.url, m.Platform.Slug)
		return m, m.Spinner.Tick
	case songwhipCrawlMsg:
		clipboard.WriteAll(msg.url)
		m.PlatformUrl = msg.url
		m.State = Done
	case errorMsg:
		m.State = HasError
	}
	return m, nil
}

func (m model) View() string {
	switch m.State {
	case HasError:
		return "Uh oh! We've encountered an error :("
	case GetOriginalUrl:
		return fmt.Sprintf("Enter a track or album URL from any supported platform...\n\n%s\n\n%s",
			m.OriginalUrl.View(),
			"(ctrl+c or esc to quit)",
		)
	case ChoosePlatform:
		return platformSelectionView(m.PlatformCursor)
	case FetchingSongwhip:
		return fmt.Sprintf("%s Getting Songwhip Data...", m.Spinner.View())
	case CrawlingSongwhip:
		return fmt.Sprintf("%s Getting %s URL...", m.Spinner.View(), m.Platform.Title)
	case Done:
		if len(m.PlatformUrl) == 0 {
			return fmt.Sprintf("Oh no! Could not find a URL for %s :(", m.Platform.Title)
		} else {
			return fmt.Sprintf(
				"Here's your %s URL! The link has been copied to your clipboard.\n\n%s\n\n(press any key to quit)",
				m.Platform.Title,
				m.PlatformUrl,
			)
		}
	}

	return ""
}
