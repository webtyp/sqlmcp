package sqlmcp

import (
	"webtyp.com/ddl"
	"webtyp.com/fmt"
	"webtyp.com/mcp"
	"webtyp.com/orm"
	"webtyp.com/storage"
)

// ExportFunc produces the DDL SQL for the currently synced schema. Supplied by the caller
// (e.g. ormc.Generator.ExportSQL bound to a dialect compiler) — sqlmcp has no model registry
// of its own to derive this from since orm dropped Open/Register (see orm docs/PLAN.md §2).
type ExportFunc func() (string, error)

// Provider implements mcp.ToolProvider for a live *orm.DB connection.
type Provider struct {
	db       *orm.DB
	exportFn ExportFunc
}

// NewProvider creates a new MCP tool provider wrapping the given DB and DDL export function.
// exportFn may be nil — db_export_schema then reports it isn't configured.
func NewProvider(db *orm.DB, exportFn ExportFunc) *Provider {
	return &Provider{db: db, exportFn: exportFn}
}

// Tools returns the MCP tools available for this DB connection.
// db_schema is only included if the underlying connection implements ddl.SchemaInspector.
func (p *Provider) Tools() []mcp.Tool {
	tools := []mcp.Tool{
		queryTool(p.db),
		execTool(p.db),
		exportTool(p.exportFn),
	}
	if _, ok := p.db.RawConn().(ddl.SchemaInspector); ok {
		tools = append([]mcp.Tool{schemaTool(p.db)}, tools...)
	}
	return tools
}

func scanRowsToText(rows storage.Rows) string {
	cols, _ := rows.Columns()
	var out fmt.Conv
	out.Write(fmt.Convert(cols).Join(" | ").String() + "\n")
	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	for rows.Next() {
		rows.Scan(ptrs...)
		parts := make([]string, len(vals))
		for i, v := range vals {
			parts[i] = fmt.Sprint(v)
		}
		out.Write(fmt.Convert(parts).Join(" | ").String() + "\n")
	}
	if err := rows.Err(); err != nil {
		out.Write("error: " + err.Error())
	}
	return out.String()
}
