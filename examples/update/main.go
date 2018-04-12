package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, err := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(database)

	rowsAffected, _ := queryBuilder.Update("userinfo", "u").SetParam(4).Set("u.username", "joe11").Set("u.flag", "1").Where("u.uid=?").PrepareAndExecute()
	fmt.Println(rowsAffected)
}
