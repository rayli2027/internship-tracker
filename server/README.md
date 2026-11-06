# Internship Tracker Server

Go (Gin) backend API for the Internship Tracker application.

## Run Locally

```bash
go run .
```

By default, the server uses built-in SQLite (`./data/internship_tracker.db`) with zero external configuration required.
The server automatically creates database tables and seeds demo accounts on first launch.

## Health Check

```bash
curl http://127.0.0.1:8080/api/health
```
