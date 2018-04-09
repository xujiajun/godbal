package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()
	queryBuilder := mysql.NewQueryBuilder(database)

	rowsAffected, _ := queryBuilder.Update("userinfo", "u").Set("u.username", "joe").Set("u.flag", "0").Where("u.uid=?").
		SetParam(4).PrepareAndExecute()
	fmt.Println(rowsAffected)
}
