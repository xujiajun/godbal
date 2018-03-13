package main

import (
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
	"log"
)

func foo() {
	panic("foo func error")
}

func main() {
	database, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	defer database.Close()

	_, err := database.Begin()

	if err != nil {
		log.Fatalln(err)
	}

	defer database.Rollback()

	queryBuilder := mysql.NewQueryBuilder(database)

	rowsAffected, err := queryBuilder.Update("userinfo", "u").Set("u.username", "?").Set("u.departname", "?").Where("u.uid=?").
		SetParameter("joe").SetParameter("tecxx").SetParameter(4).PrepareAndExecute()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Print(rowsAffected)

	foo()

	if err := database.Commit(); err != nil {
		log.Fatalln(err)
	}
}
