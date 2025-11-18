package parser

type Column struct {
	Name string
	Type string
}

type InsertStatementData struct {
	TableName string
	Pairs     map[string]string // column name -> value
}

type SelectStatementData struct {
	TableName string
	Columns   []string
}
