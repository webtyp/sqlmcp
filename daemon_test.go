package sqlmcp

import (
	"strings"
	"testing"

	"webtyp.com/context"
	"webtyp.com/ddl"
	"webtyp.com/mcp"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/storage"
)

func TestDaemonProvider_Tools(t *testing.T) {
	p := NewDaemonProvider()
	tools := p.Tools()
	if len(tools) != 4 {
		t.Errorf("expected 4 tools, got %d", len(tools))
	}

	names := map[string]bool{}
	for _, tool := range tools {
		names[tool.Name] = true
	}

	for _, name := range []string{"db_schema", "db_query", "db_exec", "db_export_schema"} {
		if !names[name] {
			t.Errorf("missing tool %s", name)
		}
	}
}

func TestDaemonProvider_NoDB(t *testing.T) {
	p := NewDaemonProvider()
	ctx := context.Background()

	for _, tool := range p.Tools() {
		_, err := tool.Execute(ctx, mcp.Request{})
		if err == nil || !strings.Contains(err.Error(), "no database configured") {
			t.Errorf("tool %s: expected 'no database configured' error, got %v", tool.Name, err)
		}
	}
}

// mockConn implements storage.Conn + ddl.SchemaInspector for the daemon provider tests.
type mockConn struct {
	queryCalled bool
	execCalled  bool
}

func (m *mockConn) Exec(query string, args ...any) error {
	m.execCalled = true
	return nil
}

func (m *mockConn) QueryRow(query string, args ...any) storage.Scanner {
	return &mockScanner{}
}

func (m *mockConn) Query(query string, args ...any) (storage.Rows, error) {
	m.queryCalled = true
	return &mockRows{}, nil
}

func (m *mockConn) Close() error { return nil }

func (m *mockConn) Compile(q storage.Query, mdl model.Model) (storage.Plan, error) {
	return storage.Plan{}, nil
}

func (m *mockConn) Tables() ([]string, error) { return []string{"users"}, nil }

func (m *mockConn) Columns(table string) ([]ddl.ColumnInfo, error) {
	return []ddl.ColumnInfo{{Name: "id", Type: "INTEGER", PK: true}}, nil
}

type mockScanner struct{}

func (m *mockScanner) Scan(dest ...any) error { return nil }

type mockRows struct{}

func (m *mockRows) Columns() ([]string, error) { return []string{"col1"}, nil }
func (m *mockRows) Next() bool                 { return false }
func (m *mockRows) Scan(dest ...any) error     { return nil }
func (m *mockRows) Close() error               { return nil }
func (m *mockRows) Err() error                 { return nil }

func TestDaemonProvider_WithDB(t *testing.T) {
	p := NewDaemonProvider()
	conn := &mockConn{}
	db := orm.New(conn)
	p.SetDB(db)
	p.SetExportFunc(func() (string, error) { return "-- exported schema", nil })

	ctx := context.Background()

	// Test db_query
	reqQuery := mcp.Request{Params: mcp.CallToolParams{Arguments: `{"SQL": "SELECT 1"}`}}
	_, err := p.Tools()[1].Execute(ctx, reqQuery) // db_query is at index 1
	if err != nil {
		t.Errorf("db_query failed: %v", err)
	}
	if !conn.queryCalled {
		t.Errorf("db_query did not call conn.Query")
	}

	// Test db_exec
	reqExec := mcp.Request{Params: mcp.CallToolParams{Arguments: `{"SQL": "UPDATE x SET y=1"}`}}
	_, err = p.Tools()[2].Execute(ctx, reqExec) // db_exec is at index 2
	if err != nil {
		t.Errorf("db_exec failed: %v", err)
	}
	if !conn.execCalled {
		t.Errorf("db_exec did not call conn.Exec")
	}

	// Test db_schema
	_, err = p.Tools()[0].Execute(ctx, mcp.Request{}) // db_schema is at index 0
	if err != nil {
		t.Errorf("db_schema failed: %v", err)
	}

	// Test db_export_schema
	_, err = p.Tools()[3].Execute(ctx, mcp.Request{}) // db_export_schema is at index 3
	if err != nil {
		t.Errorf("db_export_schema failed: %v", err)
	}

	// Test SetDB(nil)
	p.SetDB(nil)
	_, err = p.Tools()[1].Execute(ctx, reqQuery)
	if err == nil || !strings.Contains(err.Error(), "no database configured") {
		t.Errorf("expected error after SetDB(nil), got %v", err)
	}
}
