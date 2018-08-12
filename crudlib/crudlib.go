package crudlib

// A library of helper methods for gencrud-generated code.

import (
	"fmt"
	"strconv"
	"strings"
)

// BindStyle is an enum of styles of bind vars.
type BindStyle int

// Varieties for supported bind-variable markers.
const (
	QuestionMark BindStyle = iota
	DollarNum
	ColonName
)

// String to satisfy flag.Var()
func (bs *BindStyle) String() string {
	switch *bs {
	case DollarNum:
		return "dollar"
	case ColonName:
		return "name"
	default:
		return "questionmark"
	}
}

// ConstName returns the go-name of a BindStyle.
func (bs BindStyle) ConstName() string {
	switch bs {
	case DollarNum:
		return "DollarNum"
	case ColonName:
		return "ColonName"
	default:
		return "QuestionMark"
	}
}

// Set from string, to satisfy flag.Var()
func (bs *BindStyle) Set(v string) error {

	if strings.HasPrefix(v, "question") {
		*bs = QuestionMark
		return nil
	}

	if strings.HasPrefix(v, "dollar") {
		*bs = DollarNum
		return nil
	}

	if strings.HasPrefix(v, "name") || strings.HasPrefix(v, "colon") {
		*bs = ColonName
		return nil
	}

	return fmt.Errorf("unrecognized BindStyle value: %s. (should be question, dollar or name)", v)
}

// Rebind will convert question marks to another bind-var style, e.g. "$1, $2,
// ..." or ":arg, :arg, ..." (This is pretty much stolen directly from
// github.com/jmoiron/sqlx.)
func (bs BindStyle) Rebind(query string) string {

	if bs == QuestionMark {
		return query
	}

	sb := strings.Builder{}

	p := 0
	for i := strings.Index(query, "?"); i != -1; i = strings.Index(query, "?") {
		sb.WriteString(query[:i])
		switch bs {
		case DollarNum:
			p++
			sb.WriteRune('$')
			sb.WriteString(strconv.Itoa(p))
		case ColonName:
			sb.WriteString(":arg")
		}

		query = query[i+1:]
	}

	sb.WriteString(query)

	return sb.String()
}

// InsertReturningStatement returns "insert into tablename (...) values (...)
// returning id"
func InsertReturningStatement(tableName string, columnNames []string, returningColumnName string) string {

	columnNamesString := strings.Join(columnNames, ", ")
	bindMarkers := strings.Trim(strings.Repeat("?, ", len(columnNames)), ", ")

	return fmt.Sprintf("insert into %s (%s) values (%s) returning %s",
		tableName, columnNamesString, bindMarkers, returningColumnName)
}

// InsertStatement returns "insert into tablename (...) values (...)"
func InsertStatement(tableName string, columnNames []string) string {

	columnNamesString := strings.Join(columnNames, ", ")
	bindMarkers := strings.Trim(strings.Repeat("?, ", len(columnNames)), ", ")

	return fmt.Sprintf("insert into %s (%s) values (%s)",
		tableName, columnNamesString, bindMarkers)
}

// UpdateStatement returns "update tablename set ... where ..." Note that bind
// values must be used with values before keys.
func UpdateStatement(tableName string, keyColumns []string, valueColumns []string) string {

	setExpressions := strings.Join(valueColumns, " = ?, ") + " = ?"
	whereExpressions := strings.Join(keyColumns, " = ? and ") + " = ?"

	return fmt.Sprintf("update %s set %s where %s",
		tableName, setExpressions, whereExpressions)
}

// DeleteStatement returns "delete from tablename where ..."
func DeleteStatement(tableName string, keyColumns []string) string {

	whereExpressions := strings.Join(keyColumns, " = ? and ") + " = ?"

	return fmt.Sprintf("delete from %s where %s",
		tableName, whereExpressions)
}

// SelectStatement constructs a select statement for the given table and
// columns, with optional additional clauses (where, order by, etc.)
func SelectStatement(tableName string, columnNames []string, additionalClauses string) string {

	columnsString := tableName + "." + strings.Join(columnNames, ", "+tableName+".")

	selectStatement := fmt.Sprintf("select %s from %s", columnsString, tableName)

	if len(additionalClauses) == 0 {
		return selectStatement
	}

	return selectStatement + " " + additionalClauses
}
