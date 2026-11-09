package sqlmcp

import "webtyp.com/model"

var sqlPermitted = model.Permitted{
	Letters:   true,
	Numbers:   true,
	Spaces:    true,
	BreakLine: true,
	Tab:       true,
	Extra: []rune{
		'*', '=', '>', '<', '"', '\'', '(', ')', ',', '.', ';', '_', '-', '+', '/', '%', '?', '!', '@', '#', '$', '^', '&', '|', '~', '`', '[', ']', '{', '}', ':',
	},
}

var QueryArgsModel = model.Definition{
	Name: "query_args",
	Fields: model.Fields{
		{
			Name:      "SQL",
			Type:      model.Text(),
			NotNull:   true,
			Permitted: sqlPermitted,
		},
	},
}

var ExecArgsModel = model.Definition{
	Name: "exec_args",
	Fields: model.Fields{
		{
			Name:      "SQL",
			Type:      model.Text(),
			NotNull:   true,
			Permitted: sqlPermitted,
		},
	},
}
