package main

import (
	"encoding/json"
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

	//queryBuilder := mysql.NewQueryBuilder(database)
	//rows, _ := queryBuilder.Select("u.uid,u.username,p.address").From("userinfo", "u").SetFirstResult(0).
	//	SetMaxResults(3).InnerJoin("profile", "p", "u.uid = p.uid").Query()
	//
	//fmt.Print(ToJson(rows))

	queryBuilder2 := mysql.NewQueryBuilder(database)
	queryBuilder2 = queryBuilder2.Select("u.uid,u.username,p.address").From("userinfo", "u").SetFirstResult(0).
		SetMaxResults(3).RightJoin("profile", "p", "u.uid = p.uid")

	fmt.Println(queryBuilder2.GetSQL())

	rows, _ := queryBuilder2.QueryAndGetMap()

	jsonString, _ := json.Marshal(&rows)

	fmt.Println(string(jsonString))
}
