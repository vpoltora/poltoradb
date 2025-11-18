// Package parser provides functions to parse SQL-like queries from stdin.
package parser

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	MinSQLQueryWords        = 4 // minimum words before proper token splitting
	MinSQLCreateTableTokens = 8
	MinSQLInsertTokens      = 11
	MinSQLSelectTokens      = 5
)

const (
	SQLCreateTableTokensCreate = "create"
	SQLCreateTableTokensTable  = "table"

	SQLInsertTokensInsert = "insert"
	SQLInsertTokensInto   = "into"

	SQLSelectTokensSelect = "select"
	SQLSelectTokensFrom   = "from"
)

// ParseCreateTableStatement function parses a CREATE TABLE SQL-like statement.
//
// Valid input example: ["create" "table" "users" "(" "id" "INT" "," "username" "TEXT" "," "email" "TEXT" ")" ";"]
func ParseCreateTableStatement(tokens []string) ([]Column, error) {
	splittedTokens, err := splitTokens(tokens)
	if err != nil {
		return nil, fmt.Errorf("error splitting tokens: %v", err)
	}

	// Minimum: create table table_name ( col_name col_type ) ; - that's 8 tokens
	if len(splittedTokens) < MinSQLCreateTableTokens {
		return nil, fmt.Errorf("incomplete create table command")
	}

	if strings.ToLower(splittedTokens[1]) != SQLCreateTableTokensTable {
		return nil, fmt.Errorf("expected `table` keyword after `create`")
	}

	tableName := splittedTokens[2]
	if err := isValidIdentifier(tableName); err != nil {
		return nil, fmt.Errorf("invalid table name: %v", err)
	}

	if splittedTokens[3] != "(" {
		return nil, fmt.Errorf("expected '(' after table name")
	}

	// Check last tokens
	if splittedTokens[len(splittedTokens)-1] != ";" {
		return nil, fmt.Errorf("expected ';' at the end of the statement")
	}

	if splittedTokens[len(splittedTokens)-2] != ")" {
		return nil, fmt.Errorf("expected ')' before ';'")
	}

	columns := make([]Column, 0)
	for i := 4; i < len(splittedTokens)-2; i += 3 {
		// Check last column (no comma)
		if splittedTokens[i+2] == ")" {
			col, err := parseColumn(splittedTokens[i], splittedTokens[i+1])
			if err != nil {
				return nil, err
			}
			columns = append(columns, col)
			break
		}

		if splittedTokens[i+2] != "," {
			return nil, fmt.Errorf("expected ',' after column definition, got '%s'", splittedTokens[i+2])
		}

		col, err := parseColumn(splittedTokens[i], splittedTokens[i+1])
		if err != nil {
			return nil, err
		}
		columns = append(columns, col)
	}

	return columns, nil
}

func parseColumn(name, typeStr string) (Column, error) {
	if err := isValidIdentifier(name); err != nil {
		return Column{}, fmt.Errorf("invalid column name '%s': %v", name, err)
	}
	if !isValidColumnType(typeStr) {
		return Column{}, fmt.Errorf("invalid column type '%s'", typeStr)
	}
	return Column{Name: name, Type: typeStr}, nil
}

func isValidIdentifier(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("identifier cannot be empty")
	}

	for i, ch := range name {
		if i == 0 {
			if !unicode.IsLetter(ch) && ch != '_' {
				return fmt.Errorf("identifier must start with a letter or underscore")
			}
		} else {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
				return fmt.Errorf("identifier can only contain letters, digits, or underscores")
			}
		}
	}

	return nil
}

var validColumnTypes = map[string]bool{
	"int":  true,
	"text": true,
}

func isValidColumnType(typeName string) bool {
	return validColumnTypes[strings.ToLower(typeName)]
}

// ParseInsertStatement function parses an INSERT SQL-like statement.
//
// Valid input example: ["insert" "into" "users" "(" "id" "," "username" "," "email" ")" "values" "(" "1" "," "john_doe" "," "test@test.com" ")" ";" ]
func ParseInsertStatement(tokens []string) (InsertStatementData, error) {
	splittedTokens, err := splitTokens(tokens)
	if err != nil {
		return InsertStatementData{}, fmt.Errorf("error splitting tokens: %v", err)
	}

	// Minimum: insert into table ( col ) values ( val ) - that's 11 tokens
	if len(splittedTokens) < MinSQLInsertTokens {
		return InsertStatementData{}, fmt.Errorf("incomplete insert statement")
	}

	// Check "into" keyword
	if strings.ToLower(splittedTokens[1]) != SQLInsertTokensInto {
		return InsertStatementData{}, fmt.Errorf("expected `into` keyword after `insert`")
	}

	// Validate table name
	tableName := splittedTokens[2]
	if err := isValidIdentifier(tableName); err != nil {
		return InsertStatementData{}, fmt.Errorf("invalid table name: %v", err)
	}

	// Check last token - ";"
	if splittedTokens[len(splittedTokens)-1] != ";" {
		return InsertStatementData{}, fmt.Errorf("expected ';' at the end of the statement")
	}

	// Check opening parenthesis for columns
	if splittedTokens[3] != "(" {
		return InsertStatementData{}, fmt.Errorf("expected '(' after table name")
	}

	// Find "values" keyword and its position
	valuesIdx := findKeywordIndex(splittedTokens, "values")
	if valuesIdx == -1 {
		return InsertStatementData{}, fmt.Errorf("expected `values` keyword in insert statement")
	}

	// Before values keyword should be columns definition
	columns := parseColumnsList(splittedTokens, 4, valuesIdx)

	// Check opening parenthesis for values
	if valuesIdx+1 >= len(splittedTokens) || splittedTokens[valuesIdx+1] != "(" {
		return InsertStatementData{}, fmt.Errorf("expected '(' after values keyword")
	}

	// Find closing parenthesis for values
	// Closing paren should be at position len(splittedTokens)-2, because last token should be ";"
	closingParenIdx := len(splittedTokens) - 2
	if splittedTokens[closingParenIdx] != ")" {
		return InsertStatementData{}, fmt.Errorf("expected ')' before ';'")
	}

	// Extract values between ( and )
	values := make([]string, 0)
	for i := valuesIdx + 2; i < closingParenIdx; i++ {
		token := splittedTokens[i]

		// Skip commas
		if token == "," {
			continue
		}

		// Extract value (remove quotes if present)
		value := token
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.TrimPrefix(strings.TrimSuffix(value, "\""), "\"")
		}

		values = append(values, value)
	}

	// Create pairs map from columns and values
	if len(columns) != len(values) {
		return InsertStatementData{}, fmt.Errorf("number of columns (%d) does not match number of values (%d)", len(columns), len(values))
	}

	pairs := make(map[string]string)
	for i, columnName := range columns {
		pairs[columnName] = values[i]
	}

	return InsertStatementData{TableName: tableName, Pairs: pairs}, nil
}

// ParseSelectStatement function parses a SELECT SQL-like statement.
//
// Valid input example: ["select" "*" "from" "users" ";" ]
// Valid input example: ["select" "id" "," "email" "from" "users" ";" ]
func ParseSelectStatement(tokens []string) (SelectStatementData, error) {
	splittedTokens, err := splitTokens(tokens)
	if err != nil {
		return SelectStatementData{}, fmt.Errorf("error splitting tokens: %v", err)
	}

	// Minimum: select * from table ; - that's 5 tokens
	if len(splittedTokens) < MinSQLSelectTokens {
		return SelectStatementData{}, fmt.Errorf("incomplete select statement")
	}

	// Check last token - ";"
	if splittedTokens[len(splittedTokens)-1] != ";" {
		return SelectStatementData{}, fmt.Errorf("expected ';' at the end of the statement")
	}

	// Find "from" keyword and its position
	fromIdx := findKeywordIndex(splittedTokens, "from")
	if fromIdx == -1 {
		return SelectStatementData{}, fmt.Errorf("expected `from` keyword in select statement")
	}

	// Before from keyword should be columns definition
	columns := parseColumnsList(splittedTokens, 1, fromIdx)

	// After from should be table name
	if fromIdx+1 >= len(splittedTokens) {
		return SelectStatementData{}, fmt.Errorf("expected table name after `from` keyword")
	}

	tableName := splittedTokens[fromIdx+1]
	if err := isValidIdentifier(tableName); err != nil {
		return SelectStatementData{}, fmt.Errorf("invalid table name: %v", err)
	}

	return SelectStatementData{TableName: tableName, Columns: columns}, nil
}

// SplitTokens function splits tokens that may contain punctuation attached to words.
//
// CREATE: ["create" "table" "users" "(" "id" "," "int" "," "email" "," "text" ")" ";"]
// INSERT: ["insert" "into" "users" "(" "id" "," "email" ")" "values" "(" "1" "," "test@gmail.com" ")" ";"]
// SELECT: ["select" "*" "from" "users" ";"]
// SELECT: ["select" "id" "," "email" "from" "users" ";"]
func splitTokens(tokens []string) ([]string, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens to split")
	}

	output := make([]string, 0, len(tokens))
	for _, token := range tokens {
		token = strings.TrimSuffix(token, "\n")
		tokenLen := len(token)

		// parse tokens like "(id,"
		if strings.HasPrefix(token, "(") && strings.HasSuffix(token, ",") {
			if tokenLen < 3 {
				return nil, fmt.Errorf("invalid token: %s - length must be > 2", token)
			}

			output = append(output, "(", token[1:tokenLen-1], ",")

			continue
		}

		// parse tokens like "email);" or "int);"
		if strings.HasSuffix(token, ");") {
			if tokenLen < 3 {
				return nil, fmt.Errorf("invalid token: %s - length must be >= 3", token)
			}

			if !strings.HasSuffix(token[:tokenLen-1], ")") {
				return nil, fmt.Errorf("invalid token: %s - expected ')' before ';'", token)
			}

			output = append(output, token[:tokenLen-2], ")", ";")

			continue
		}

		// parse tokens like "int,"
		if strings.HasSuffix(token, ",") {
			if tokenLen < 2 {
				return nil, fmt.Errorf("invalid token: %s - length must be > 1", token)
			}

			output = append(output, token[:tokenLen-1], ",")

			continue
		}

		// parse tokens like "(id"
		if strings.HasPrefix(token, "(") {
			if tokenLen == 1 {
				return nil, fmt.Errorf("invalid token: %s - length must be > 1", token)
			}

			output = append(output, "(")
			output = append(output, token[1:])

			continue
		}

		// parse tokens like "email)"
		if strings.HasSuffix(token, ")") {
			if tokenLen < 2 {
				return nil, fmt.Errorf("invalid token: %s - length must be >= 2", token)
			}

			output = append(output, token[:tokenLen-1], ")")

			continue
		}

		// parse tokens like "users;"
		if strings.HasSuffix(token, ";") {
			if tokenLen < 2 {
				return nil, fmt.Errorf("invalid token: %s - length must be >= 2", token)
			}

			output = append(output, token[:tokenLen-1], ";")

			continue
		}

		output = append(output, token)
	}

	return output, nil
}

func findKeywordIndex(tokens []string, keyword string) int {
	for i, token := range tokens {
		if strings.ToLower(token) == keyword {
			return i
		}
	}
	return -1
}

func parseColumnsList(splittedTokens []string, startIdx, endIdx int) []string {
	columns := make([]string, 0)
	for i := startIdx; i < endIdx; i++ {
		token := splittedTokens[i]

		// Skip commas and closing parentheses
		if token == "," || token == ")" {
			continue
		}

		columns = append(columns, token)
	}
	return columns
}
