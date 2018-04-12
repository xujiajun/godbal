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

	rowsAffected, _ := queryBuilder.Update("userinfo", "u").SetParam(4).Set("u.username", "joe11").Set("u.flag", "1").Where("u.uid=?").PrepareAndExecute()
	fmt.Println(rowsAffected)
}
