package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func showResultsDialog(app *tview.Application, pages *tview.Pages, serverURL, token, agentID string) {
	results, err := fetchResults(serverURL, token, agentID)
	if err != nil {
		showModal(app, pages, fmt.Sprintf("failed to fetch results: %v", err))
		return
	}

	text := ""
	for _, r := range results {
		text += "[yellow]═══════════════════════════════════════[-]\n"
		text += fmt.Sprintf("[cyan]Task:[-] %s\n", r.TaskID)
		text += fmt.Sprintf("[cyan]Status:[-] [green]%s[-]\n", r.Status)
		text += fmt.Sprintf("[cyan]Time:[-] %s\n", r.Timestamp.Format(time.RFC3339))
		if r.Error != "" {
			text += fmt.Sprintf("[red]Error:[-] %s\n", r.Error)
		}
		if r.Output != "" {
			text += fmt.Sprintf("[green]Output:[-]\n%s\n", r.Output)
		} else {
			text += "[yellow](no output)[-]\n"
		}
	}
	if text == "" {
		text = "[yellow]<no results>[-]"
	}

	view := tview.NewTextView().SetDynamicColors(true).SetText(text)
	view.SetBorder(true).SetTitle(fmt.Sprintf("Results (%s) - [Esc] back", agentID)).SetTitleAlign(tview.AlignLeft)

	view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			pages.RemovePage("results")
			app.SetFocus(pages)
			return nil
		}
		return event
	})

	pages.AddPage("results", view, true, true)
	app.SetFocus(view)
}

func fetchResults(serverURL, token, agentID string) ([]taskResult, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/operator/results?agent_id=%s", serverURL, agentID), nil)
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
