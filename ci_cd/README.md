# CI/CD

This folder documents the CI/CD layout and conventions for the repo.

Source of truth for GitHub Actions lives in `.github/workflows/`.

Planned pipelines (simulation-only):
- `dev`: unit tests + lint
- `QA`: integration tests (mock agent ↔ server)
- `main`: tag releases + generate docs

