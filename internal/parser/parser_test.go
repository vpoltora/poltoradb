package parser

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseCreateTableStatement(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		want    []Column
		wantErr string
	}{
		{
			name:   "valid create with two columns",
			tokens: sqlTokens("create table users (id INT, email TEXT);"),
			want: []Column{
				{Name: "id", Type: "INT"},
				{Name: "email", Type: "TEXT"},
			},
		},
		{
			name:    "invalid table name",
			tokens:  sqlTokens("create table 7users (id INT);"),
			wantErr: "invalid table name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCreateTableStatement(tt.tokens)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("unexpected columns. want=%v got=%v", tt.want, got)
			}
		})
	}
}

func TestParseInsertStatement(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		want    InsertStatementData
		wantErr string
	}{
		{
			name:   "valid insert",
			tokens: sqlTokens(`insert into users (id, username, email) values (1, "john_doe", "test@test.com");`),
			want: InsertStatementData{
				TableName: "users",
				Pairs: map[string]string{
					"id":       "1",
					"username": "john_doe",
					"email":    "test@test.com",
				},
			},
		},
		{
			name:    "columns and values mismatch",
			tokens:  sqlTokens(`insert into users (id, username, email) values (1, "john_doe");`),
			wantErr: "does not match number of values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInsertStatement(tt.tokens)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.TableName != tt.want.TableName {
				t.Fatalf("unexpected table name. want=%s got=%s", tt.want.TableName, got.TableName)
			}

			if !reflect.DeepEqual(got.Pairs, tt.want.Pairs) {
				t.Fatalf("unexpected pairs. want=%v got=%v", tt.want.Pairs, got.Pairs)
			}
		})
	}
}

func TestParseSelectStatement(t *testing.T) {
	tests := []struct {
		name    string
		tokens  []string
		want    SelectStatementData
		wantErr string
	}{
		{
			name:   "valid select",
			tokens: sqlTokens("select id, email from users;"),
			want: SelectStatementData{
				TableName: "users",
				Columns:   []string{"id", "email"},
			},
		},
		{
			name:    "invalid table name",
			tokens:  sqlTokens("select * from 7users;"),
			wantErr: "invalid table name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSelectStatement(tt.tokens)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got.TableName != tt.want.TableName {
				t.Fatalf("unexpected table name. want=%s got=%s", tt.want.TableName, got.TableName)
			}

			if !reflect.DeepEqual(got.Columns, tt.want.Columns) {
				t.Fatalf("unexpected columns. want=%v got=%v", tt.want.Columns, got.Columns)
			}
		})
	}
}

func sqlTokens(input string) []string {
	return strings.Split(input, " ")
}
