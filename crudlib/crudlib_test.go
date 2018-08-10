package crudlib_test

import (
	"testing"

	"github.com/pdk/crudgen/crudlib"
)

func TestRebind(t *testing.T) {

	s1 := "delete from foo where id = ?"
	result := crudlib.QuestionMark.Rebind(s1)
	if result != s1 {
		t.Errorf("expected %s, but got %s", s1, result)
	}

	cases := []struct {
		input, edollar, enamed string
	}{
		{
			input:   "select * from t where id = ?",
			edollar: "select * from t where id = $1",
			enamed:  "select * from t where id = :arg",
		},
		{
			input:   "select * from t where id = ? or name like ?",
			edollar: "select * from t where id = $1 or name like $2",
			enamed:  "select * from t where id = :arg or name like :arg",
		},
		{
			input:   "update t set a = ?, b = ? where id = ?",
			edollar: "update t set a = $1, b = $2 where id = $3",
			enamed:  "update t set a = :arg, b = :arg where id = :arg",
		},
		{
			input:   "select * from t",
			edollar: "select * from t",
			enamed:  "select * from t",
		},
		{
			input:   "select * from t where id = ? and alpha = beta",
			edollar: "select * from t where id = $1 and alpha = beta",
			enamed:  "select * from t where id = :arg and alpha = beta",
		},
	}

	for _, c := range cases {

		result := crudlib.DollarNum.Rebind(c.input)
		if result != c.edollar {
			t.Errorf("expected %s, but got %s", c.edollar, result)
		}

		result = crudlib.ColonName.Rebind(c.input)
		if result != c.enamed {
			t.Errorf("expected %s, but got %s", c.enamed, result)
		}
	}
}

func TestInsertStatement(t *testing.T) {
	cases := []struct {
		tableName   string
		columnNames []string
		expected    string
	}{
		{
			tableName:   "foo",
			columnNames: []string{"alpha", "beta", "gamma"},
			expected:    "insert into foo (alpha, beta, gamma) values (?, ?, ?)",
		},
		{
			tableName:   "foo",
			columnNames: []string{"alpha"},
			expected:    "insert into foo (alpha) values (?)",
		},
	}

	for _, c := range cases {

		result := crudlib.InsertStatement(c.tableName, c.columnNames)
		if result != c.expected {
			t.Errorf("expected %s, but got %s", c.expected, result)
		}
	}
}

func TestInsertReturningStatement(t *testing.T) {
	cases := []struct {
		tableName       string
		columnNames     []string
		returningColumn string
		expected        string
	}{
		{
			tableName:       "foo",
			columnNames:     []string{"alpha", "beta", "gamma"},
			returningColumn: "id",
			expected:        "insert into foo (alpha, beta, gamma) values (?, ?, ?) returning id",
		},
		{
			tableName:       "foo",
			columnNames:     []string{"alpha"},
			returningColumn: "foo_dim_id",
			expected:        "insert into foo (alpha) values (?) returning foo_dim_id",
		},
	}

	for _, c := range cases {

		result := crudlib.InsertReturningStatement(c.tableName, c.columnNames, c.returningColumn)
		if result != c.expected {
			t.Errorf("expected %s, but got %s", c.expected, result)
		}
	}
}

func TestUpdateStatement(t *testing.T) {
	cases := []struct {
		tableName    string
		keyColumns   []string
		valueColumns []string
		expected     string
	}{
		{
			tableName:    "foo",
			keyColumns:   []string{"id"},
			valueColumns: []string{"alpha", "beta", "gamma"},
			expected:     "update foo set alpha = ?, beta = ?, gamma = ? where id = ?",
		},
		{
			tableName:    "foo",
			keyColumns:   []string{"id"},
			valueColumns: []string{"alpha"},
			expected:     "update foo set alpha = ? where id = ?",
		},
		{
			tableName:    "foo",
			keyColumns:   []string{"key1", "key2"},
			valueColumns: []string{"alpha"},
			expected:     "update foo set alpha = ? where key1 = ? and key2 = ?",
		},
	}

	for _, c := range cases {

		result := crudlib.UpdateStatement(c.tableName, c.keyColumns, c.valueColumns)
		if result != c.expected {
			t.Errorf("expected %s, but got %s", c.expected, result)
		}
	}
}

func TestDeleteStatement(t *testing.T) {
	cases := []struct {
		tableName  string
		keyColumns []string
		expected   string
	}{
		{
			tableName:  "foo",
			keyColumns: []string{"id"},
			expected:   "delete from foo where id = ?",
		},
		{
			tableName:  "foo",
			keyColumns: []string{"key1", "key2"},
			expected:   "delete from foo where key1 = ? and key2 = ?",
		},
	}

	for _, c := range cases {

		result := crudlib.DeleteStatement(c.tableName, c.keyColumns)
		if result != c.expected {
			t.Errorf("expected %s, but got %s", c.expected, result)
		}
	}
}

func TestSelectStatement(t *testing.T) {
	cases := []struct {
		tableName         string
		columns           []string
		expected          string
		additionalClauses string
	}{
		{
			tableName: "foo",
			columns:   []string{"alpha", "beta", "gamma"},
			expected:  "select foo.alpha, foo.beta, foo.gamma from foo",
		},
		{
			tableName: "foo",
			columns:   []string{"alpha"},
			expected:  "select foo.alpha from foo",
		},
		{
			tableName:         "foo",
			columns:           []string{"alpha", "beta", "gamma"},
			additionalClauses: "where foo.alpha = ?",
			expected:          "select foo.alpha, foo.beta, foo.gamma from foo where foo.alpha = ?",
		},
		{
			tableName:         "foo",
			columns:           []string{"alpha"},
			additionalClauses: "order by alpha",
			expected:          "select foo.alpha from foo order by alpha",
		},
	}

	for _, c := range cases {

		result := crudlib.SelectStatement(c.tableName, c.columns, c.additionalClauses)
		if result != c.expected {
			t.Errorf("expected %s, but got %s", c.expected, result)
		}
	}
}
