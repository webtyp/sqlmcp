# Agent Instructions

## Coding Standards
- **No repeated string literals:** Do not repeat long string literals or character blocks. Use constants or package-level variables instead to improve maintainability.
- **No standard library for core functions:** Always use `webtyp` equivalents instead of the Go standard library where available:
    - `webtyp.com/fmt` instead of `fmt`
    - `webtyp.com/json` instead of `encoding/json`
    - `webtyp.com/model`
    - `webtyp.com/orm`
    - `webtyp.com/mcp`
- **JSON in tests:** For JSON operations in tests, use **only** `webtyp/json`.

## MCP and sqlmcp specifics
- **JSON Schema:** Generating JSON Schema is NOT the responsibility of `sqlmcp`. This is handled by `webtyp/mcp` via `Tool.Args` (which must implement `model.Fielder`).
- **No JSON Schema logic:** Do not include `encodeSchema` functions or manual `InputSchema` JSON strings in tool definitions.
- **SQL Validation:** When defining models for SQL input, use `model.Text()` and provide an explicit `Permitted` whitelist that includes necessary SQL symbols (quotes, operators, etc.), as the default `Text` kind blocks them for XSS protection.
