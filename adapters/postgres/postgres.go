package postgres

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/crgimenes/pgfs/config"
	_ "github.com/lib/pq"
)

// Table metadata
type Table struct {
	Name string
}

var (
	cfg config.Config
	db  *sql.DB
)

// Load config and connect do db
func Load() {
	var err error
	cfg = config.Get()
	db, err = sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		log.Fatal(err)
	}
}

// ListTables get all tables from a schema and return in a slice
func ListTables() (t []Table, err error) {
	rows, err := db.Query("SELECT tablename FROM pg_tables where schemaname = $1", cfg.SchemaName)
	if err != nil {
		return
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return
		}
		t = append(t, Table{Name: name})
	}
	return
}

// LoadTableJSON load a table and trasform to JSON
func LoadTableJSON(tableName string) (ret []byte, err error) {
	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	var t []map[string]string
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		for i := range cols {
			vals[i] = new(scanner)
		}
		if e := rows.Scan(vals...); e != nil {
			log.Fatal(e)
		}
		m := make(map[string]string)
		for i, column := range cols {
			var scanner = vals[i].(*scanner)
			m[column] = scanner.String()
		}
		t = append(t, m)
	}
	ret, err = json.MarshalIndent(t, "", "\t")
	return
}

// LoadTableCSV load a table and trasform to CSV in a byte array
func LoadTableCSV(tableName string) (ret []byte, err error) {
	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err = w.Write(cols)
	if err != nil {
		return
	}

	for rows.Next() {
		vals := make([]interface{}, len(cols))
		for i := range cols {
			vals[i] = new(scanner)
		}
		err = rows.Scan(vals...)
		if err != nil {
			return
		}
		var row []string
		for i := 0; i < len(cols); i++ {
			var scanner = vals[i].(*scanner)
			row = append(row, scanner.String())
		}
		err = w.Write(row)
		if err != nil {
			return
		}
	}
	w.Flush()
	ret = b.Bytes()
	return
}

type scanner struct {
	value interface{}
}

func (scanner *scanner) Scan(src interface{}) error {
	switch src.(type) {
	case int:
		scanner.value = src.(int)
	case int64:
		scanner.value = src.(int64)
	case float64:
		scanner.value = src.(float64)
	case bool:
		scanner.value = src.(bool)
	case string:
		scanner.value = src.(string)
	case []byte:
		scanner.value = src.([]byte)
	case time.Time:
		scanner.value = src.(time.Time)
	case nil:
		scanner.value = nil
	}
	return nil
}

func (scanner *scanner) String() string {
	switch scanner.value.(type) {
	case int:
		return fmt.Sprintf("%v", scanner.value.(int))
	case int64:
		return fmt.Sprintf("%v", scanner.value.(int64))
	case float64:
		return fmt.Sprintf("%v", scanner.value.(float64))
	case bool:
		return fmt.Sprintf("%v", scanner.value.(bool))
	case string:
		return scanner.value.(string)
	case []byte:
		return string(scanner.value.([]byte))
	case time.Time:
		return fmt.Sprintf("%v", scanner.value.(time.Time))
	case nil:
		return "NIL"
	default:
		return fmt.Sprintf("%T not implemented", scanner.value)
	}
}
