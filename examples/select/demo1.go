package main

import (
	"encoding/json"
	"fmt"
	"github.com/xujiajun/godbal"
	"github.com/xujiajun/godbal/driver/mysql"
)

func main() {
	connection, err := godbal.NewMysql("root:123@tcp(127.0.0.1:3306)/test?charset=utf8").Open()
	if err != nil {
		panic(err)
	}

	queryBuilder := mysql.NewQueryBuilder(connection)
	rows, _ := queryBuilder.Select("uid,username,price,flag").From("userinfo", "").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").Query()

	jsonString, _ := json.Marshal(&rows)
	fmt.Print(string(jsonString))

}
