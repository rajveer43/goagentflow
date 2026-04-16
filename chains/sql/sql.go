package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/rajveer43/goagentflow/runtime"
)

// Chain generates SQL queries from natural language descriptions.
// Pattern: Composition - LLM-based SQL generation
// Input: natural language query; Output: SQL string
// NOTE: This is a stub. Real implementation would validate/execute SQL.
type Chain struct {
	llm          runtime.LLM
	schema       string // database schema description
	promptTemplate string
	dialect      string // SQL dialect (postgres, mysql, sqlite, etc.)
}

// Input represents the input to a SQL chain.
type Input struct {
	Question string
	Schema   string // optional schema override
}

// Output represents the output from a SQL chain.
type Output struct {
	SQL        string
	Explanation string
}

// New creates a new SQL generation chain.
// llm: language model for SQL generation
// schema: description of database schema (tables, columns, etc.)
func New(llm runtime.LLM, schema string) *Chain {
	return &Chain{
		llm:            llm,
		schema:         schema,
		promptTemplate: defaultSQLPrompt,
		dialect:        "SQL", // generic
	}
}

// SetDialect sets the SQL dialect (postgres, mysql, sqlite, etc.).
func (c *Chain) SetDialect(dialect string) {
	c.dialect = dialect
}

// SetSchema updates the database schema description.
func (c *Chain) SetSchema(schema string) {
	c.schema = schema
}

// Run implements runtime.Chain interface.
// Input: string (natural language question) or Input struct
// Output: string (SQL) or Output struct (with explanation)
func (c *Chain) Run(ctx context.Context, input any) (any, error) {
	question := ""
	schema := c.schema

	switch v := input.(type) {
	case string:
		question = v
	case Input:
		question = v.Question
		if v.Schema != "" {
			schema = v.Schema
		}
	default:
		return nil, fmt.Errorf("expected string or Input, got %T", input)
	}

	if question == "" {
		return nil, fmt.Errorf("question cannot be empty")
	}

	if schema == "" {
		return nil, fmt.Errorf("database schema is required")
	}

	// Build prompt
	prompt := strings.NewReplacer(
		"{dialect}", c.dialect,
		"{schema}", schema,
		"{question}", question,
	).Replace(c.promptTemplate)

	// Get LLM response
	sqlQuery, err := c.llm.Complete(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("SQL generation failed: %w", err)
	}

	// Extract SQL from response (may include explanation)
	sql, explanation := c.extractSQL(sqlQuery)

	return Output{
		SQL:         sql,
		Explanation: explanation,
	}, nil
}

// extractSQL parses the LLM response to extract SQL and any explanation.
func (c *Chain) extractSQL(response string) (string, string) {
	// Look for SQL wrapped in ```sql ... ``` blocks
	sqlStart := strings.Index(response, "```sql")
	if sqlStart >= 0 {
		sqlStart += 6
		sqlEnd := strings.Index(response[sqlStart:], "```")
		if sqlEnd >= 0 {
			sql := strings.TrimSpace(response[sqlStart : sqlStart+sqlEnd])
			explanation := strings.TrimSpace(response[:sqlStart-6]) +
				           strings.TrimSpace(response[sqlStart+sqlEnd+3:])
			return sql, explanation
		}
	}

	// Look for SELECT/INSERT/UPDATE/DELETE at start
	upper := strings.ToUpper(strings.TrimSpace(response))
	if strings.HasPrefix(upper, "SELECT") ||
	   strings.HasPrefix(upper, "INSERT") ||
	   strings.HasPrefix(upper, "UPDATE") ||
	   strings.HasPrefix(upper, "DELETE") {
		return strings.TrimSpace(response), ""
	}

	// Otherwise treat entire response as SQL
	return strings.TrimSpace(response), ""
}

const defaultSQLPrompt = `You are a SQL query generator. Generate a {dialect} query to answer the user's question.

Database Schema:
{schema}

User Question: {question}

Generate the SQL query. If you need to provide an explanation, put it before the SQL wrapped in triple backtick blocks.

SQL Query:`
