# E2E Tests

Run:

```sh
./tests/e2e/run.sh
```

This validates, end-to-end, in a local lab environment:
- Postgres + teamserver boot
- operator login
- agent register + heartbeat
- operator task enqueue + result persistence
- telemetry latest endpoint
