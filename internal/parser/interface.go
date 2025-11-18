package parser

type Parser interface {
	ParseCreateTableStatement(splittedTokens []string) ([]Column, error)
	ParseSelectStatement(splittedTokens []string) (SelectStatementData, error)
	ParseInsertStatement(splittedTokens []string) (InsertStatementData, error)
}
