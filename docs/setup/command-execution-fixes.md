# Command Execution Fixes

## Issues Found and Fixed

### 1. **Task Timeout Unmarshaling Bug** (CRITICAL)
**File**: [teamserver/internal/server/server.go](teamserver/internal/server/server.go#L218)

**Problem**: The JSON struct field was `TimeoutSecond` (singular) but the UI sends `timeout_seconds` (plural). This caused the timeout field to be zero-valued when parsing task creation requests.

```go
// BEFORE (WRONG)
var req struct {
    TimeoutSecond int `json:"timeout_seconds"`
}

// AFTER (FIXED)
var req struct {
    TimeoutSeconds int `json:"timeout_seconds"`
}
```

**Impact**: Tasks were created with 0-second timeout, likely causing immediate failures or unexpected behavior.

---

### 2. **Missing AgentID in Task Results** (HIGH)
**File**: [teamserver/internal/registry/registry.go](teamserver/internal/registry/registry.go#L273)

**Problem**: The `GetResults()` function was not populating the `AgentID` field when retrieving results from the database. The UI needs this field to match results to agents.

```go
// BEFORE (INCOMPLETE)
var r api.TaskResult
if err := rows.Scan(&r.TaskID, &r.Status, &r.Output, &r.Error, &r.Timestamp); err != nil {
    continue
}
out = append(out, r)  // r.AgentID is empty!

// AFTER (FIXED)
var r api.TaskResult
var taskID string
if err := rows.Scan(&taskID, &r.Status, &r.Output, &r.Error, &r.Timestamp); err != nil {
    continue
}
r.TaskID = taskID
r.AgentID = agentID  // Now properly set
out = append(out, r)
```

**Impact**: UI couldn't display which agent completed a task, breaking the tasking workflow.

---

## Testing

### Manual Test Flow

1. **Start teamserver**:
```bash
export REDFORGE_CORS_ORIGINS="http://localhost:5173"
export REDFORGE_ADMIN_PASS="admin123"
cd teamserver
go run ./cmd/teamserver
```

2. **Register an agent**:
```bash
# In another terminal
cd agent
cargo run
```

3. **Create a task via UI or CLI**:
```bash
curl -X POST https://localhost:9080/api/operator/task \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"agent_id":"<agent-id>","command":"pwd","args":[],"timeout_seconds":30}'
```

4. **Check task results**:
```bash
curl https://localhost:9080/api/operator/results?agent_id=<agent-id> \
  -H "Authorization: Bearer <token>"
```

### Expected Behavior

- ✅ Task is created with correct timeout
- ✅ Agent receives task in next heartbeat
- ✅ Agent executes command
- ✅ Result is submitted and stored
- ✅ UI can retrieve results and display them

---

## Files Modified

- [teamserver/internal/server/server.go](teamserver/internal/server/server.go) - Fixed JSON struct field name
- [teamserver/internal/registry/registry.go](teamserver/internal/registry/registry.go) - Fixed missing AgentID in results

---

## Related Issues

- CORS configuration must be set for UI-to-teamserver communication (see [cors-setup.md](cors-setup.md))
- Protocol mismatch between UI and teamserver (see [api.ts fix](ui/src/lib/api.ts))

---

## Database Schema

For reference, the tasks table structure:
```sql
CREATE TABLE tasks (
    task_id TEXT PRIMARY KEY,
    agent_id TEXT NOT NULL,
    command TEXT NOT NULL,
    args JSON,
    timeout_seconds INT,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

Task lifecycle:
1. `pending` - Created, waiting for agent to retrieve
2. `in-progress` - Agent retrieved the task
3. `success` / `error` / `timeout` - Agent completed and submitted result
