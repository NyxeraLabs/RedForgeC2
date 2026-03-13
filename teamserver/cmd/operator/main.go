package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const defaultServerURL = "http://localhost:9080"

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type agentEntry struct {
	AgentID    string            `json:"agent_id"`
	Token      string            `json:"token"`
	OS         string            `json:"os"`
	Arch       string            `json:"arch"`
	Hostname   string            `json:"hostname"`
	Version    string            `json:"version"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	LastSeen   string            `json:"last_seen"`
	Registered string            `json:"registered"`
}

type taskRequest struct {
	AgentID       string   `json:"agent_id"`
	Command       string   `json:"command"`
	Args          []string `json:"args"`
	TimeoutSecond int      `json:"timeout_seconds"`
}

type taskResult struct {
	AgentID   string    `json:"agent_id"`
	TaskID    string    `json:"task_id"`
	Status    string    `json:"status"`
	Output    string    `json:"output"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorBlack
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorBlack
	tview.Styles.BorderColor = tcell.ColorDarkRed
	tview.Styles.TitleColor = tcell.ColorRed
	tview.Styles.GraphicsColor = tcell.ColorRed
	tview.Styles.PrimaryTextColor = tcell.ColorWhite
	tview.Styles.SecondaryTextColor = tcell.ColorGray
	tview.Styles.TertiaryTextColor = tcell.ColorDarkGray
	tview.Styles.InverseTextColor = tcell.ColorBlack

	serverURL := os.Getenv("REDFORGE_SERVER_URL")
	if serverURL == "" {
		serverURL = defaultServerURL
	}

	var token string

	app := tview.NewApplication()
	pages := tview.NewPages()

	agentTable, refreshAgents, focusAgents := newAgentTable(app, pages, serverURL, &token)
	pages.AddPage("agents", agentTable, true, false)

	loginForm := newLoginForm(app, pages, serverURL, &token, func() {
		if refreshAgents != nil {
			refreshAgents()
		}
		if focusAgents != nil {
			focusAgents()
		}
	})
	pages.AddPage("login", loginForm, true, true)

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("failed to run operator UI: %v", err)
	}
}

func newLoginForm(app *tview.Application, pages *tview.Pages, serverURL string, token *string, onLogin func()) *tview.Form {
	username := "admin"
	password := "redforge"

	form := tview.NewForm().
		AddInputField("Username", username, 20, nil, func(text string) { username = text }).
		AddPasswordField("Password", password, 20, '*', func(text string) { password = text }).
		AddButton("Login", func() {
			// Avoid blocking the UI thread with network I/O.
			showModal(app, pages, "Logging in...")
			go func(u, p string) {
				t, err := doLogin(serverURL, u, p)
				app.QueueUpdateDraw(func() {
					pages.RemovePage("modal")
					if err != nil {
						showModal(app, pages, fmt.Sprintf("login failed: %v", err))
						return
					}
					*token = t
					pages.HidePage("login")
					pages.SwitchToPage("agents")
					if onLogin != nil {
						onLogin()
					}
				})
			}(username, password)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).
		SetBorderColor(tcell.ColorDarkRed).
		SetTitle("REDFORGE-C2 OPERATOR LOGIN").
		SetTitleColor(tcell.ColorRed).
		SetTitleAlign(tview.AlignLeft)
	form.SetFieldTextColor(tcell.ColorWhite).
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetButtonBackgroundColor(tcell.ColorDarkRed).
		SetButtonTextColor(tcell.ColorWhite)
	return form
}

func newAgentTable(app *tview.Application, pages *tview.Pages, serverURL string, token *string) (*tview.Flex, func(), func()) {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("[red::b]REDFORGE[-] [gray]TUI v1.0[-]  [::d]R refresh  T task  V results  Q quit[-]")
	header.SetBorder(true).SetBorderColor(tcell.ColorDarkRed)

	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetBorderColor(tcell.ColorDarkRed).SetTitle("SESSIONS").SetTitleColor(tcell.ColorRed)
	table.SetSelectedStyle(tcell.StyleDefault.Background(tcell.ColorDarkRed).Foreground(tcell.ColorWhite))

	status := tview.NewTextView().SetDynamicColors(true)
	status.SetText("[gray]Ready.[-]")
	status.SetBorder(true).SetBorderColor(tcell.ColorDarkRed).SetTitle("STATUS").SetTitleColor(tcell.ColorRed)

	refresh := func() {
		if strings.TrimSpace(*token) == "" {
			status.SetText("[yellow]Not authenticated. Log in first.[-]")
			return
		}
		status.SetText(fmt.Sprintf("[gray]Refreshing from %s...[-]", serverURL))
		currentToken := strings.TrimSpace(*token)
		go func(tok string) {
			agents, err := fetchAgents(serverURL, tok)
			app.QueueUpdateDraw(func() {
				// Token may have been cleared/rotated while we were fetching.
				if strings.TrimSpace(*token) == "" || strings.TrimSpace(*token) != tok {
					return
				}
				if err != nil {
					status.SetText(fmt.Sprintf("[red]refresh failed: %v", err))
					return
				}
				table.Clear()
				headers := []string{"ID", "Host", "OS", "Arch", "Version", "Last Seen"}
				for i, h := range headers {
					table.SetCell(0, i, tview.NewTableCell(h).
						SetSelectable(false).
						SetAttributes(tcell.AttrBold).
						SetTextColor(tcell.ColorRed))
				}
				for r, a := range agents {
					row := r + 1
					table.SetCell(row, 0, tview.NewTableCell(a.AgentID).SetTextColor(tcell.ColorWhite))
					table.SetCell(row, 1, tview.NewTableCell(a.Hostname).SetTextColor(tcell.ColorWhite))
					table.SetCell(row, 2, tview.NewTableCell(a.OS).SetTextColor(tcell.ColorGray))
					table.SetCell(row, 3, tview.NewTableCell(a.Arch).SetTextColor(tcell.ColorGray))
					table.SetCell(row, 4, tview.NewTableCell(a.Version).SetTextColor(tcell.ColorGray))
					table.SetCell(row, 5, tview.NewTableCell(a.LastSeen).SetTextColor(tcell.ColorGray))
				}
				status.SetText(fmt.Sprintf("[green]Connected[-] to [white]%s[-]  Agents: [white]%d[-]", serverURL, len(agents)))
			})
		}(currentToken)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(header, 3, 0, false)
	flex.AddItem(table, 0, 1, true)
	flex.AddItem(status, 3, 0, false)

	table.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
		case 'r', 'R':
			refresh()
		case 't', 'T':
			row, _ := table.GetSelection()
			if row <= 0 {
				return event
			}
			agentID := table.GetCell(row, 0).Text
			doTaskDialog(app, pages, serverURL, *token, agentID)
		case 'v', 'V':
			row, _ := table.GetSelection()
			if row <= 0 {
				return event
			}
			agentID := table.GetCell(row, 0).Text
			showResultsDialog(app, pages, serverURL, *token, agentID)
		}
		return event
	})

	refresh()
	return flex, refresh, func() { app.SetFocus(table) }
}

func showModal(app *tview.Application, pages *tview.Pages, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("modal")
		})
	modal.SetTitle("REDFORGE").SetTitleColor(tcell.ColorRed)
	pages.AddPage("modal", modal, true, true)
	app.SetFocus(modal)
}

func doLogin(serverURL, username, password string) (string, error) {
	body, _ := json.Marshal(loginRequest{Username: username, Password: password})
	resp, err := http.Post(serverURL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}
	var out loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.Token, nil
}

func fetchAgents(serverURL, token string) ([]agentEntry, error) {
	req, _ := http.NewRequest("GET", serverURL+"/api/operator/agents", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var agents []agentEntry
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		return nil, err
	}
	return agents, nil
}

func doTaskDialog(app *tview.Application, pages *tview.Pages, serverURL, token, agentID string) {
	var cmd, args string
	form := tview.NewForm().
		AddInputField("Command", "ls", 40, nil, func(text string) { cmd = text }).
		AddInputField("Args (space separated)", "", 40, nil, func(text string) { args = text }).
		AddButton("Send", func() {
			argSlice := []string{}
			if strings.TrimSpace(args) != "" {
				argSlice = strings.Fields(args)
			}
			task := taskRequest{AgentID: agentID, Command: cmd, Args: argSlice, TimeoutSecond: 30}
			if err := submitTask(serverURL, token, task); err != nil {
				showModal(app, pages, fmt.Sprintf("task failed: %v", err))
				return
			}
			showModal(app, pages, "task queued")
		}).
		AddButton("Cancel", func() {
			pages.SwitchToPage("agents")
		})
	form.SetBorder(true).SetTitle(fmt.Sprintf("Send Task to %s", agentID))

	pages.AddPage("task", form, true, true)
	pages.SwitchToPage("task")
}

func submitTask(serverURL, token string, task taskRequest) error {
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", serverURL+"/api/operator/task", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}
