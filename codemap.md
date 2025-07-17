# Code Map

- `cmd/server` - entry point for the HTTP server.
- `internal/domain` - domain entities.
- `internal/usecase` - business logic use cases.
- `internal/adapters/web` - web/http adapter using Go templates.
- `internal/adapters/repository/sqlite` - sqlite repository implementation.
- `internal/database` - database setup and migrations.
- `migrations` - SQL migration files.
- `internal/adapters/web/templates` - HTML templates for the web adapter.
- Uses the upstream modernc.org/sqlite pure Go SQLite driver.
