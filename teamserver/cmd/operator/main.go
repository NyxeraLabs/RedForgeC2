package main

import (
	"bytes"
	"crypto/tls"
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

const defaultServerURL = "https://localhost:9080"

// secureClient is an HTTP client configured to handle self-signed TLS certificates
// commonly used in C2 development environments.
var secureClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
	Timeout: 15 * time.Second,
}

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

	agentView, refreshAgents, focusAgents, selectedAgentID := newAgentTable(app, pages, serverURL, &token)
	pages.AddPage("agents", agentView, true, false)

	tasksView, refreshTasks, focusTasks := newTasksTable(app, pages, serverURL, &token)
	pages.AddPage("tasks", tasksView, true, false)

	loginForm := newLoginForm(app, pages, serverURL, &token, func() {
		if refreshAgents != nil {
			refreshAgents()
		}
		if focusAgents != nil {
			focusAgents()
		}
		if refreshTasks != nil {
			refreshTasks()
		}
	})
	pages.AddPage("login", loginForm, true, true)

	active := "agents"
	switchTo := func(name string) {
		active = name
		pages.SwitchToPage(name)
		switch name {
		case "agents":
			if focusAgents != nil {
				focusAgents()
			}
		case "tasks":
			if focusTasks != nil {
				focusTasks()
			}
		}
	}

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Let table/form navigation work normally
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyBacktab ||
			event.Key() == tcell.KeyUp || event.Key() == tcell.KeyDown ||
			event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight ||
			event.Key() == tcell.KeyEnter {
			return event
		}

		if event.Key() == tcell.KeyCtrlC {
			app.Stop()
			return nil
		}

		focus := app.GetFocus()
		_, isInput := focus.(*tview.InputField)
		_, isForm := focus.(*tview.Form)

		// If in an input field or form, only handle Ctrl+C and Escape
		if isInput || isForm {
			if event.Key() == tcell.KeyEscape {
				app.Stop()
				return nil
			}
			return event
		}

		// Handle shortcuts only when not in input fields
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
			return nil
		case '1':
			switchTo("agents")
			return nil
		case '2':
			switchTo("tasks")
			return nil
		case 'r', 'R':
			if active == "tasks" && refreshTasks != nil {
				refreshTasks()
				return nil
			}
			if refreshAgents != nil {
				refreshAgents()
				return nil
			}
		case 't', 'T':
			if active != "agents" {
				return event
			}
			agentID := ""
			if selectedAgentID != nil {
				agentID = selectedAgentID()
			}
			if agentID != "" {
				doTaskDialog(app, pages, serverURL, token, agentID)
				return nil
			}
		case 'v', 'V':
			if active != "agents" {
				return event
			}
			agentID := ""
			if selectedAgentID != nil {
				agentID = selectedAgentID()
			}
			if agentID != "" {
				showResultsDialog(app, pages, serverURL, token, agentID)
				return nil
			}
		}
		return event
	})

	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		log.Fatalf("failed to run operator UI: %v", err)
	}
}

func newLoginForm(app *tview.Application, pages *tview.Pages, serverURL string, token *string, onLogin func()) *tview.Form {
	username := "admin"
	password := ""

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

func newAgentTable(app *tview.Application, pages *tview.Pages, serverURL string, token *string) (*tview.Flex, func(), func(), func() string) {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("[red::b]REDFORGE[-] [gray]TUI v1.0[-]  [::d]1 Agents  2 Tasks  R refresh  T task  V results  Q quit[-]")
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

	// Don't close on Escape; only 'q'/'Q' will quit
	refresh()
	return flex, refresh, func() { app.SetFocus(table) }, func() string {
		row, _ := table.GetSelection()
		if row <= 0 {
			return ""
		}
		return table.GetCell(row, 0).Text
	}
}

type taskEntry struct {
	TaskID    string    `json:"task_id"`
	AgentID   string    `json:"agent_id"`
	Command   string    `json:"command"`
	Args      []string  `json:"args"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

func newTasksTable(app *tview.Application, pages *tview.Pages, serverURL string, token *string) (*tview.Flex, func(), func()) {
	header := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft).
		SetText("[red::b]REDFORGE[-] [gray]TUI v1.0[-]  [::d]1 Agents  2 Tasks  R refresh  Q quit[-]")
	header.SetBorder(true).SetBorderColor(tcell.ColorDarkRed)

	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetBorderColor(tcell.ColorDarkRed).SetTitle("TASKS").SetTitleColor(tcell.ColorRed)
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
			tasks, err := fetchTasks(serverURL, tok)
			app.QueueUpdateDraw(func() {
				if strings.TrimSpace(*token) == "" || strings.TrimSpace(*token) != tok {
					return
				}
				if err != nil {
					status.SetText(fmt.Sprintf("[red]refresh failed: %v", err))
					return
				}
				table.Clear()
				headers := []string{"Status", "Agent", "Task", "Command", "Updated"}
				for i, h := range headers {
					table.SetCell(0, i, tview.NewTableCell(h).
						SetSelectable(false).
						SetAttributes(tcell.AttrBold).
						SetTextColor(tcell.ColorRed))
				}
				for r, t := range tasks {
					row := r + 1
					statusColor := tcell.ColorGray
					switch strings.ToLower(t.Status) {
					case "pending":
						statusColor = tcell.ColorYellow
					case "in-progress", "in_progress":
						statusColor = tcell.ColorLightCyan
					case "complete", "completed", "done":
						statusColor = tcell.ColorGreen
					case "failed", "error":
						statusColor = tcell.ColorRed
					}
					table.SetCell(row, 0, tview.NewTableCell(strings.ToUpper(t.Status)).SetTextColor(statusColor))
					table.SetCell(row, 1, tview.NewTableCell(shortText(t.AgentID, 10)).SetTextColor(tcell.ColorWhite))
					table.SetCell(row, 2, tview.NewTableCell(shortText(t.TaskID, 12)).SetTextColor(tcell.ColorWhite))
					table.SetCell(row, 3, tview.NewTableCell(t.Command).SetTextColor(tcell.ColorGray))
					table.SetCell(row, 4, tview.NewTableCell(t.UpdatedAt.Format(time.RFC3339)).SetTextColor(tcell.ColorGray))
				}
				status.SetText(fmt.Sprintf("[green]Connected[-] to [white]%s[-]  Tasks: [white]%d[-]", serverURL, len(tasks)))
			})
		}(currentToken)
	}

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(header, 3, 0, false)
	flex.AddItem(table, 0, 1, true)
	flex.AddItem(status, 3, 0, false)

	refresh()
	return flex, refresh, func() { app.SetFocus(table) }
}

func fetchTasks(serverURL, token string) ([]taskEntry, error) {
	req, _ := http.NewRequest("GET", serverURL+"/api/operator/tasks", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := secureClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var tasks []taskEntry
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func shortText(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
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
	resp, err := secureClient.Post(serverURL+"/api/login", "application/json", bytes.NewReader(body))
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
	resp, err := secureClient.Do(req)
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
	resp, err := secureClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}

func fetchAgentResults(serverURL, token, agentID string) ([]taskResult, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/operator/results/%s", serverURL, agentID), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := secureClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	var results []taskResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, err
	}
	return results, nil
}
