package main

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()

	c, err := db.QueryContext(ctx, "select * from one")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	for c.Next() {
		x := 0
		y := ""
		err := c.Scan(&x, &y)
		if err != nil {
			panic(err)
		}
		println(x, y)
	}

	println("OK")
}
