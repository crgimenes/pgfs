package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func ListTables() (err error) {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=cesar sslmode=disable")
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM pg_catalog.pg_tables")
	if err != nil {
		log.Fatal(err)
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", cols)

	for rows.Next() {
		vals := make([]interface{}, len(cols))
		for i := range cols {
			vals[i] = new(scanner)
		}
		if e := rows.Scan(vals...); e != nil {
			log.Fatal(e)
		}
		for idx, column := range cols {
			var scanner = vals[idx].(*scanner)
			fmt.Printf("%v: %v\n", column, scanner.String())
		}
	}

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
