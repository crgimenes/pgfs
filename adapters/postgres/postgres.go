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

type Table struct {
	Name string
}

func ListTables() (t []Table, err error) {
	cfg := config.Get()
	db, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT tablename FROM pg_tables where schemaname = $1", cfg.SchemaName)
	if err != nil {
		log.Fatal(err)
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

func LoadTableJSON(tableName string) (ret []byte, err error) {
	cfg := config.Get()
	db, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		log.Fatal(err)
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
		for idx, column := range cols {
			var scanner = vals[idx].(*scanner)
			m[column] = scanner.String()
		}
		t = append(t, m)
	}
	ret, err = json.MarshalIndent(t, "", "\t")
	return
}

func LoadTableCSV(tableName string) (ret []byte, err error) {
	cfg := config.Get()
	db, err := sql.Open("postgres", cfg.DataSourceName)
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + tableName)
	if err != nil {
		log.Fatal(err)
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
	valid bool
	value interface{}
}

func (scanner *scanner) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		scanner.value = src.(int64)
		scanner.valid = true
	case float64:
		scanner.value = src.(float64)
		scanner.valid = true
	case bool:
		scanner.value = src.(bool)
		scanner.valid = true
	case string:
		scanner.value = src.(string)
		scanner.valid = true
	case []byte:
		scanner.value = src.([]byte)
		scanner.valid = true
	case time.Time:
		scanner.value = src.(time.Time)
		scanner.valid = true
	case nil:
		scanner.value = nil
		scanner.valid = true
	}
	return nil
}

func (scanner *scanner) String() string {
	switch scanner.value.(type) {
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
