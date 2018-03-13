package main

import (
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	database, _ := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()
	queryBuilder := mysql.NewQueryBuilder(database)

	res, _ := queryBuilder.Update("userinfo", "u").Set("u.username", "?").Set("u.departname", "?").Where("u.uid=?").
		SetParameter("joe").SetParameter("tecxx").SetParameter(5).PrepareAndExecute()
	fmt.Print(res)
}
