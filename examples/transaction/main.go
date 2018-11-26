package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func foo() {
	panic("foo func error")
}

func main() {
	database, err := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	defer database.Close()

	tx, err := database.Begin()

	if err != nil {
		log.Fatalln(err)
	}

	defer tx.Rollback()

	queryBuilder := mysql.NewQueryBuilder(database)

	queryBuilder.Update("userinfo", "u").Set("u.username", "joe").Set("u.departname", "tecxx").Where("u.uid=?").
		SetParam(4)

	res, err := tx.PrepareAndExecute(queryBuilder)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(res.RowsAffected())

	foo()

	if err := tx.Commit(); err != nil {
		log.Fatalln(err)
	}
}
