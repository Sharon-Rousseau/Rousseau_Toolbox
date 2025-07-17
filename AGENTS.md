# AGENTS.md - Development Guidelines & Standards

This document defines the development standards, architectural patterns, and workflow methodologies for this project.

## 🎯 Core Principles

### Mandatory Standards
- **ALWAYS** maintain a `codemap.md` file for project navigation
- **ALWAYS** follow hexagonal architecture patterns
- **ALWAYS** write comprehensive tests for all features
- **ALWAYS** use templ for type-safe templating
- **ALWAYS** use HTMX for dynamic DOM interactions
- **ALWAYS** use HTMX Server-Sent Events (SSE) for real-time data streams

## 🏗️ Architecture Standards

### Hexagonal Architecture
```
┌─────────────────────────────────────────────────────────┐
│                    Presentation Layer                    │
│                  (HTTP, CLI, WebSocket)                  │
├─────────────────────────────────────────────────────────┤
│                    Application Layer                     │
│                    (Use Cases, DTOs)                     │
├─────────────────────────────────────────────────────────┤
│                      Domain Layer                        │
│                (Entities, Business Rules)                │
├─────────────────────────────────────────────────────────┤
│                  Infrastructure Layer                    │
│              (Database, External Services)               │
└─────────────────────────────────────────────────────────┘
```

### Directory Structure
```
project/
├── cmd/                    # Application entrypoints
├── internal/              # Private application code
│   ├── domain/           # Core business entities
│   ├── usecase/          # Application business logic
│   └── adapters/         # External interfaces
│       ├── web/          # HTTP handlers & templates
│       └── repository/   # Data persistence
├── migrations/           # Database migrations
└── tests/               # Integration tests
```

## 🎨 Frontend Architecture

### templ - Type-Safe Templating

#### Overview
templ provides compile-time type-safe HTML templates that integrate seamlessly with Go's type system.

#### Benefits
- **Type Safety**: Compile-time validation of template variables
- **Performance**: Templates compile to efficient Go code
- **IDE Support**: Full autocompletion and refactoring support
- **Zero Runtime Overhead**: No template parsing at runtime

#### Implementation
```go
// dashboard.templ
templ DashboardPage(user User, items []Item) {
    <!DOCTYPE html>
    <html>
        <head>
            <title>Dashboard - { user.Name }</title>
            @HeadAssets()
        </head>
        <body>
            @Header(user)
            <main>
                @ItemList(items)
            </main>
        </body>
    </html>
}

// Handler usage
func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
    user := getUserFromContext(r.Context())
    items, _ := h.service.GetUserItems(user.ID)
    
    component := DashboardPage(user, items)
    component.Render(r.Context(), w)
}
```

### HTMX - Progressive Enhancement

#### Core Attributes
| Attribute | Purpose | Example |
|-----------|---------|---------|
| `hx-get/post/put/delete` | HTTP requests | `hx-post="/api/items"` |
| `hx-target` | Response destination | `hx-target="#item-list"` |
| `hx-swap` | Swap strategy | `hx-swap="innerHTML"` |
| `hx-trigger` | Event trigger | `hx-trigger="click"` |
| `hx-confirm` | User confirmation | `hx-confirm="Delete item?"` |

#### Common Patterns

**Form Submission**
```html
<form hx-post="/items" 
      hx-target="#items-container" 
      hx-swap="beforeend">
    <input name="title" required>
    <button type="submit">Add Item</button>
</form>
```

**Delete with Confirmation**
```html
<button hx-delete="/items/{{ .ID }}"
        hx-target="closest .item"
        hx-swap="outerHTML"
        hx-confirm="Are you sure?">
    Delete
</button>
```

**Polling Updates**
```html
<div hx-get="/status"
     hx-trigger="every 5s"
     hx-swap="innerHTML">
    Loading status...
</div>
```

### Server-Sent Events (SSE) - Real-Time Updates

#### When to Use SSE
- **Real-time dashboards**: Live metrics and monitoring
- **Progress indicators**: Long-running task updates
- **Live feeds**: Chat messages, notifications
- **Collaborative features**: Shared document editing

#### Implementation Pattern

**Client-Side**
```html
<div hx-ext="sse" 
     sse-connect="/events/stream" 
     sse-swap="message">
    <ul id="messages">
        <!-- Messages appear here -->
    </ul>
</div>
```

**Server-Side**
```go
func (h *Handler) EventStream(w http.ResponseWriter, r *http.Request) {
    // Set SSE headers
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    
    // Create event channel
    events := h.eventBus.Subscribe(r.Context())
    
    for {
        select {
        case event := <-events:
            // Render templ component
            component := MessageItem(event.Message)
            
            // Send SSE event
            fmt.Fprintf(w, "event: message\n")
            fmt.Fprintf(w, "data: ")
            component.Render(r.Context(), w)
            fmt.Fprintf(w, "\n\n")
            
            w.(http.Flusher).Flush()
            
        case <-r.Context().Done():
            return
        }
    }
}
```

## 🧪 Testing Strategy

### Test Pyramid
```
        /\
       /E2E\      <- Full system tests (Go + real DB)
      /------\
     /  Integ  \   <- API & database tests  
    /------------\
   /     Unit     \  <- Business logic tests
  /----------------\
```

### Coverage Requirements
- **Unit Tests**: >90% coverage for business logic
- **Integration Tests**: All API endpoints and database operations
- **E2E Tests**: Critical user workflows with real database

### Testing Philosophy
**NEVER mock the SQLite database**. Always use temporary database files for authentic testing.

### Test Organization

#### Unit Tests (with mocked dependencies)
```go
// budget_test.go - Unit test example
func TestBudgetService(t *testing.T) {
    t.Run("validates budget name", func(t *testing.T) {
        // Only mock the repository interface, never the DB
        repo := &mockRepository{}
        service := NewBudgetService(repo)
        
        err := service.CreateBudget("")
        assert.Error(t, err, "empty name should fail")
    })
}
```

#### Integration Tests (with real SQLite)
```go
// budget_integration_test.go
func TestBudgetRepository(t *testing.T) {
    // Create temporary SQLite database
    db, cleanup := createTempDB(t)
    defer cleanup()
    
    repo := sqlite.NewRepo(db)
    
    t.Run("creates and retrieves budget", func(t *testing.T) {
        // Test with real database
        err := repo.Create("Test Budget")
        assert.NoError(t, err)
        
        budgets, err := repo.List()
        assert.NoError(t, err)
        assert.Len(t, budgets, 1)
    })
}

// Helper function for temp DB
func createTempDB(t *testing.T) (*sql.DB, func()) {
    tmpfile, err := os.CreateTemp("", "test-*.db")
    if err != nil {
        t.Fatal(err)
    }
    
    db, err := sql.Open("sqlite", tmpfile.Name())
    if err != nil {
        t.Fatal(err)
    }
    
    // Run migrations
    if err := database.RunMigrations(db, "../../migrations"); err != nil {
        t.Fatal(err)
    }
    
    cleanup := func() {
        db.Close()
        os.Remove(tmpfile.Name())
    }
    
    return db, cleanup
}
```

#### E2E Tests (full system in Go)
```go
// e2e_test.go - Written in Go, not external tools
func TestFullUserJourney(t *testing.T) {
    // Setup test server with temp database
    db, cleanup := createTempDB(t)
    defer cleanup()
    
    // Initialize full application stack
    repo := sqlite.NewRepo(db)
    service := usecase.NewBudgetService(repo)
    server := web.NewServer(service)
    
    // Start test server
    ts := httptest.NewServer(server)
    defer ts.Close()
    
    t.Run("complete budget workflow", func(t *testing.T) {
        // 1. Get initial page
        resp, err := http.Get(ts.URL)
        assert.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
        
        // 2. Create budget via form
        form := url.Values{}
        form.Add("name", "Monthly Budget")
        
        resp, err = http.PostForm(ts.URL+"/add", form)
        assert.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
        
        // 3. Verify budget appears in list
        body, _ := io.ReadAll(resp.Body)
        assert.Contains(t, string(body), "Monthly Budget")
        
        // 4. Delete budget
        req, _ := http.NewRequest("DELETE", ts.URL+"/budget/1", nil)
        resp, err = http.DefaultClient.Do(req)
        assert.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
    })
}
```

## 🚀 Development Workflow

### 1. Problem Analysis
- Understand requirements completely
- Identify affected components
- Review existing patterns
- Plan incremental changes

### 2. Implementation Steps
1. Update domain models if needed
2. Implement use case logic
3. Create/update templ components
4. Add HTMX interactions
5. Write comprehensive tests
6. Update documentation

### 3. Code Review Checklist
- [ ] Tests pass with >90% coverage
- [ ] No security vulnerabilities
- [ ] Follows hexagonal architecture
- [ ] templ components are type-safe
- [ ] HTMX interactions are efficient
- [ ] Error handling is comprehensive
- [ ] Documentation is updated

### 4. Deployment Process
1. Run full test suite
2. Build and validate binaries
3. Apply database migrations
4. Deploy with zero-downtime strategy
5. Monitor metrics and logs

## 🔒 Security Standards

### Input Validation
- Validate all user input at boundaries
- Use parameterized queries for database
- Sanitize output in templates
- Implement CSRF protection

### Authentication & Authorization
- Use secure session management
- Implement proper access controls
- Log security events
- Rate limit sensitive endpoints

## 📊 Performance Guidelines

### Optimization Priorities
1. **Database Queries**: Use indexes, avoid N+1
2. **Template Rendering**: Cache common components
3. **HTMX Requests**: Minimize payload size
4. **SSE Streams**: Implement backpressure

### Monitoring
- Response time percentiles (p50, p95, p99)
- Error rates by endpoint
- Database query performance
- Memory and CPU usage

## 🛠️ Tooling

### Required Tools
```bash
# Install templ CLI
go install github.com/a-h/templ/cmd/templ@latest

# Generate templ code
templ generate

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

### Development Commands
```bash
# Watch mode for templ files
templ generate --watch

# Run with live reload
air

# Database migrations
migrate up
```

## 📚 Additional Resources

### Documentation
- [templ Documentation](https://templ.guide)
- [HTMX Documentation](https://htmx.org)
- [Hexagonal Architecture Guide](https://alistair.cockburn.us/hexagonal-architecture/)

### Best Practices
- Keep components small and focused
- Use semantic HTML elements
- Minimize JavaScript usage
- Optimize for progressive enhancement
- Design for mobile-first

---

*This document serves as the authoritative guide for development standards and practices in this project.*