package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/vpoltora/poltoradb/internal/parser"
)

// ParseQuery reads a SQL-like query from stdin, tokenizes it, and parses it into structured data.
//
// Input examples:
// ./poltoradb create table users (id INT, username TEXT, email TEXT);
// ./poltoradb insert into users (id, username, email) values (1, "john_doe", "test@gmail.com");
// ./poltoradb select * from users;
func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")

	query, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("error reading input:", err)
		return
	}

	tokens := strings.Fields(strings.TrimSpace(query))
	if len(tokens) == 0 {
		fmt.Println("no input provided")

		return
	}

	switch strings.ToLower(tokens[0]) {
	case parser.SQLCreateTableTokensCreate:
		columns, err := parser.ParseCreateTableStatement(tokens)
		if err != nil {
			fmt.Println("error parsing create table statement:", err)
			return
		}

		fmt.Printf("Parsed columns: %+v\n", columns)
	case parser.SQLInsertTokensInsert:
		parsedInsertStatementData, err := parser.ParseInsertStatement(tokens)
		if err != nil {
			fmt.Println("error parsing insert into statement:", err)
			return
		}

		fmt.Printf("Parsed insert statement data: %+v\n", parsedInsertStatementData)
	case parser.SQLSelectTokensSelect:
		parsedSelectStatementData, err := parser.ParseSelectStatement(tokens)
		if err != nil {
			fmt.Println("error parsing select statement:", err)
			return
		}

		fmt.Printf("Parsed select statement data: %+v\n", parsedSelectStatementData)
	default:
		fmt.Println("Unknown command")
	}
}
