package cfg

import ("database/sql"
	_ "github.com/lib/pq")

var DB *sql.DB

func init() {

	var err error

	DB, err = sql.Open("postgres", "postgres://postgres:sql@localhost/coach?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err)
	}
}
