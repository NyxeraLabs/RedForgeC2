package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type dashAgent struct {
	AgentID   string
	Hostname  string
	OS        string
	Arch      string
	LastSeen  time.Time
	Registered time.Time
}

type dashTask struct {
	TaskID    string
	AgentID   string
	Command   string
	Status    string
	UpdatedAt time.Time
}

func runDashboard(ctx context.Context, pool *pgxpool.Pool) {
	t := time.NewTicker(2 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			agents := queryAgents(ctx, pool)
			tasks := queryTasks(ctx, pool)
			renderDashboard(agents, tasks)
		}
	}
}

func queryAgents(ctx context.Context, pool *pgxpool.Pool) []dashAgent {
	rows, err := pool.Query(ctx, "SELECT agent_id, hostname, os, arch, last_seen, registered FROM agents ORDER BY last_seen DESC LIMIT 12")
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := make([]dashAgent, 0)
	for rows.Next() {
		var a dashAgent
		if err := rows.Scan(&a.AgentID, &a.Hostname, &a.OS, &a.Arch, &a.LastSeen, &a.Registered); err != nil {
			continue
		}
		out = append(out, a)
	}
	return out
}

func queryTasks(ctx context.Context, pool *pgxpool.Pool) []dashTask {
	rows, err := pool.Query(ctx, "SELECT task_id, agent_id, command, status, updated_at FROM tasks ORDER BY updated_at DESC LIMIT 14")
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := make([]dashTask, 0)
	for rows.Next() {
		var t dashTask
		if err := rows.Scan(&t.TaskID, &t.AgentID, &t.Command, &t.Status, &t.UpdatedAt); err != nil {
			continue
		}
		out = append(out, t)
	}
	return out
}

func renderDashboard(agents []dashAgent, tasks []dashTask) {
	now := time.Now()

	pending := 0
	inprog := 0
	done := 0
	for _, t := range tasks {
		switch strings.ToLower(t.Status) {
		case "pending":
			pending++
		case "in-progress", "in_progress":
			inprog++
		default:
			done++
		}
	}

	// Clear screen + home cursor
	fmt.Fprint(os.Stdout, "\033[H\033[2J")

	header := fmt.Sprintf("REDFORGE-C2 TEAMSERVER  |  %s", now.Format("2006-01-02 15:04:05"))
	fmt.Println(boxLine(header, 78))
	fmt.Println()

	fmt.Println(sectionTitle("AGENTS (last heartbeat)"))
	if len(agents) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, a := range agents {
			age := now.Sub(a.LastSeen).Round(time.Second)
			sev := "OK"
			if age > 4*time.Minute {
				sev = "STALE"
			} else if age > 90*time.Second {
				sev = "WARN"
			}
			fmt.Printf("  %-16s  %-10s  %-9s  %-8s  hb=%-6s  %s\n",
				short(a.AgentID, 16),
				short(a.Hostname, 10),
				short(strings.ToUpper(a.OS), 9),
				short(strings.ToUpper(a.Arch), 8),
				short(age.String(), 6),
				sev,
			)
		}
	}
	fmt.Println()

	fmt.Println(sectionTitle("TASK QUEUE (latest)"))
	fmt.Printf("  pending=%d  in-progress=%d  done=%d\n", pending, inprog, done)
	if len(tasks) == 0 {
		fmt.Println("  (none)")
	} else {
		sort.Slice(tasks, func(i, j int) bool { return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt) })
		for _, t := range tasks {
			age := now.Sub(t.UpdatedAt).Round(time.Second)
			fmt.Printf("  %-10s  agent=%-10s  %-12s  age=%-6s  %s\n",
				short(strings.ToUpper(t.Status), 10),
				short(t.AgentID, 10),
				short(t.TaskID, 12),
				short(age.String(), 6),
				short(t.Command, 30),
			)
		}
	}

	fmt.Println()
	fmt.Println(boxLine("Tip: set REDFORGE_CLI_DASHBOARD=0 to disable", 78))
}

func sectionTitle(s string) string {
	return fmt.Sprintf("== %s ==", s)
}

func boxLine(s string, w int) string {
	if len(s) > w-4 {
		s = s[:w-4]
	}
	padding := w - 4 - len(s)
	return "┌" + strings.Repeat("─", w-2) + "┐\n" +
		"│ " + s + strings.Repeat(" ", padding) + " │\n" +
		"└" + strings.Repeat("─", w-2) + "┘"
}

func short(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

