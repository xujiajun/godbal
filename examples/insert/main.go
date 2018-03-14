package main

import (
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()

	queryBuilder := mysql.NewQueryBuilder(database)

	sql := queryBuilder.Insert("userinfo").Value("username", "johnny").Value("departname", "tec").Value("created", "1521010136").GetSQL()

	fmt.Print(sql)
	queryBuilder.PrepareAndExecute()
}
