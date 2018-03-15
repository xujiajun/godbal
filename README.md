# godbal  [![GoDoc](https://godoc.org/github.com/xujiajun/godbal/driver/mysql?status.svg)](https://godoc.org/github.com/xujiajun/godbal/driver/mysql) [![Go Report Card](https://goreportcard.com/badge/github.com/xujiajun/godbal)](https://goreportcard.com/report/github.com/xujiajun/godbal)  <a href="https://travis-ci.org/xujiajun/godbal"><img src="https://travis-ci.org/xujiajun/godbal.svg?branch=master" alt="Build Status"></a> [![Coverage Status](https://coveralls.io/repos/github/xujiajun/godbal/badge.svg?branch=master)](https://coveralls.io/github/xujiajun/godbal?branch=master)  [![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/xujiajun/godbal/master/LICENSE)
Database Abstraction Layer (dbal) for go (now only support mysql)

## Motivation

I wanted a DBAL that ***no ORM***、***no Reflect***, support ***SQL builder***  following good practices and well tested code.

## Requirements

Go 1.7 or above.

## Installation

```
go get github.com/xujiajun/godbal
```

## Supported Databases

* mysql

## Getting Started

```
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

	queryBuilder := mysql.NewQueryBuilder(database)
	sql := queryBuilder.Select("uid,username,price,flag").From("userinfo", "").SetFirstResult(0).
		SetMaxResults(3).OrderBy("uid", "DESC").GetSQL()

	fmt.Print(sql) // SELECT uid,username,price,flag FROM userinfo ORDER BY uid DESC LIMIT 0,3
	fmt.Print("\n")

	rows, _ := queryBuilder.Query()

	jsonString, _ := json.Marshal(&rows)
	fmt.Print(string(jsonString)) 
  // result like: {"0":{"flag":"1","price":"111.00","uid":"6","username":"johnny2"},"1":{"flag":"1","price":"111.00","uid":"5","username":"johnny2"},"2":{"flag":"0","price":"123.99","uid":"4","username":"joe"}}
}

```

## Examples

* [select](https://github.com/xujiajun/godbal/blob/master/examples/select/main.go)

* [insert](https://github.com/xujiajun/godbal/blob/master/examples/insert/main.go)

* [delete](https://github.com/xujiajun/godbal/blob/master/examples/delete/main.go)

* [update](https://github.com/xujiajun/godbal/blob/master/examples/update/main.go)

* [join](https://github.com/xujiajun/godbal/blob/master/examples/join/main.go)

* [transaction](https://github.com/xujiajun/godbal/blob/master/examples/transaction/main.go)

## Contributing

If you'd like to help out with the project. You can put up a Pull Request.

## Author

* [xujiajun](https://github.com/xujiajun)

## License

The godbal is open-sourced software licensed under the [MIT Licensed](http://www.opensource.org/licenses/MIT)
