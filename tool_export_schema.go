package sqlmcp

import (
	"webtyp.com/context"
	"webtyp.com/fmt"
	"webtyp.com/mcp"
)

func exportTool(exportFn ExportFunc) mcp.Tool {
	return mcp.Tool{
		Name:        "db_export_schema",
		Description: "Export the full CREATE TABLE DDL for all synced tables as SQL text.",
		Resource:    "database",
		Action:      'r',
		Execute: func(ctx *context.Context, req mcp.Request) (*mcp.Result, error) {
			return executeExportSchema(exportFn, req)
		},
	}
}

func executeExportSchema(exportFn ExportFunc, _ mcp.Request) (*mcp.Result, error) {
	if exportFn == nil {
		return nil, fmt.Err("no database configured: call start_development first")
	}
	sql, err := exportFn()
	if err != nil {
		return nil, err
	}
	return mcp.Text(sql), nil
}
