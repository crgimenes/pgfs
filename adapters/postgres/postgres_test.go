package postgres

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func mockDB() (mock sqlmock.Sqlmock, err error) {
	db, mock, err = sqlmock.New()
	return
}

func Test_ListTables(t *testing.T) {
	mock, err := mockDB()
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	mock.ExpectQuery("SELECT tablename").WillReturnRows(
		sqlmock.NewRows([]string{"tablename"}).AddRow("test"))

	table, err := ListTables()
	if err != nil {
		t.Errorf("expected no errors, but got %v", err)
	}
	if len(table) != 1 {
		t.Errorf("expected 1, but got %v", len(table))
	}
	if table[0].Name != "test" {
		t.Errorf("expected \"test\", but got %v", table[0].Name)
	}
}
