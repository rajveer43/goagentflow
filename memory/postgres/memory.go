package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/rajveer43/goagentflow/runtime"
)

// Memory implements runtime.Memory using PostgreSQL as persistent backend.
// Uses two tables: messages (append-only log) and kv_store (key-value pairs).
// Pattern: Repository + Factory
// DSA: B-tree indexes on (session_id, created_at) for efficient message range queries
type Memory struct {
	conn      *pgx.Conn
	sessionID string
	mu        sync.RWMutex
}

// New creates a new Postgres-backed memory instance.
// dsn: Postgres connection string (e.g., "postgres://user:pass@localhost:5432/db")
// sessionID: unique identifier for this agent's memory scope
func New(ctx context.Context, dsn, sessionID string) (*Memory, error) {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres connection failed: %w", err)
	}

	// Ensure tables exist
	if err := createTables(ctx, conn); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("table creation failed: %w", err)
	}

	return &Memory{
		conn:      conn,
		sessionID: sessionID,
	}, nil
}

// createTables idempotently creates the messages and kv_store tables.
func createTables(ctx context.Context, conn *pgx.Conn) error {
	schema := `
	-- Messages table (append-only log)
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		session_id TEXT NOT NULL,
		role TEXT NOT NULL,
		content TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_messages_session_created
		ON messages(session_id, created_at);

	-- Key-value store table
	CREATE TABLE IF NOT EXISTS kv_store (
		session_id TEXT NOT NULL,
		key TEXT NOT NULL,
		value JSONB NOT NULL,
		updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (session_id, key)
	);
	CREATE INDEX IF NOT EXISTS idx_kv_session
		ON kv_store(session_id);
	`

	_, err := conn.Exec(ctx, schema)
	return err
}

// AddMessage appends a message to the message log.
func (m *Memory) AddMessage(ctx context.Context, msg runtime.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	query := `
	INSERT INTO messages (session_id, role, content)
	VALUES ($1, $2, $3)
	`
	return m.conn.QueryRow(ctx, query, m.sessionID, msg.Role, msg.Content).Scan()
}

// GetMessages retrieves all messages for this session in chronological order.
// Uses the B-tree index on (session_id, created_at).
func (m *Memory) GetMessages(ctx context.Context) ([]runtime.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	query := `
	SELECT role, content
	FROM messages
	WHERE session_id = $1
	ORDER BY created_at ASC
	`

	rows, err := m.conn.Query(ctx, query, m.sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]runtime.Message, 0)
	for rows.Next() {
		var msg runtime.Message
		if err := rows.Scan(&msg.Role, &msg.Content); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// Set stores or updates a key-value pair in JSONB format.
func (m *Memory) Set(ctx context.Context, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO kv_store (session_id, key, value)
	VALUES ($1, $2, $3::JSONB)
	ON CONFLICT (session_id, key) DO UPDATE
		SET value = EXCLUDED.value, updated_at = CURRENT_TIMESTAMP
	`

	return m.conn.QueryRow(ctx, query, m.sessionID, key, string(data)).Scan()
}

// Get retrieves a value from the key-value store.
func (m *Memory) Get(ctx context.Context, key string) (any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	query := `
	SELECT value FROM kv_store
	WHERE session_id = $1 AND key = $2
	`

	var jsonVal string
	err := m.conn.QueryRow(ctx, query, m.sessionID, key).Scan(&jsonVal)
	if err == pgx.ErrNoRows {
		return nil, nil // key not found (matches inmemory behavior)
	}
	if err != nil {
		return nil, err
	}

	var result any
	if err := json.Unmarshal([]byte(jsonVal), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Close closes the Postgres connection.
func (m *Memory) Close(ctx context.Context) error {
	return m.conn.Close(ctx)
}
