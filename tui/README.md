# RedForgeC2 TUI (MVP)

Minimal operator TUI built with Python `textual`.

## Run

```sh
cd tui
python3 -m venv .venv
.venv/bin/pip install -r requirements.txt

# Optional env overrides
export REDFORGE_TEAMSERVER_URL="http://localhost:9080"
export REDFORGE_USERNAME="admin"
export REDFORGE_PASSWORD="<admin-password>"

.venv/bin/python app.py
```
