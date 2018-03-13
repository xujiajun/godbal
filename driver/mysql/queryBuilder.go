package mysql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

const (
	SELECT = iota
	DELETE
	UPDATE
	INSERT
)

const (
	ISSORT    = "isSort"
	ISJOIN    = "isJoin"
	ISDEFAULT = "isDefault"
	INNER     = "INNER"
	LEFT      = "LEFT"
	RIGHT     = "RIGHT"
)

type QueryBuilder struct {
	firstResult int
	maxResults  int
	state       *sql.Stmt
	queryType   int
	sqlParts    map[string]interface{}
	database    *Database
	params      []interface{}
	flag        string
}

var sqlParts = map[string]interface{}{
	"select":  "",
	"from":    nil,
	"where":   "",
	"groupBy": "",
	"having":  "",
	"orderBy": map[string]string{},
	"values":  map[string]string{},
	"set":     map[string]string{},
	"join":    map[string]string{},
}

var params = []interface{}{}

func NewQueryBuilder(database *Database) *QueryBuilder {
	return &QueryBuilder{
		firstResult: 0,
		maxResults:  -1,
		queryType:   SELECT,
		sqlParts:    sqlParts,
		database:    database,
		params:      params,
		flag:        ISDEFAULT,
	}
}

func (queryBuilder *QueryBuilder) Select(value interface{}) *QueryBuilder {
	queryBuilder.queryType = SELECT
	queryBuilder.sqlParts["select"] = value

	return queryBuilder
}

func (queryBuilder *QueryBuilder) From(table string, alias string) *QueryBuilder {
	queryBuilder.setFromWrap(table, alias)

	return queryBuilder
}

func (queryBuilder *QueryBuilder) Update(table string, alias string) *QueryBuilder {
	queryBuilder.queryType = UPDATE
	queryBuilder.setFromWrap(table, alias)

	return queryBuilder
}

func (queryBuilder *QueryBuilder) Set(key string, val string) *QueryBuilder {
	queryBuilder.sqlParts["set"].(map[string]string)[key] = val

	return queryBuilder
}

func (queryBuilder *QueryBuilder) OrderBy(sort string, order string) *QueryBuilder {

	queryBuilder.flag = ISSORT
	if order == "" {
		order = "ASC"
	}

	queryBuilder.sqlParts["orderBy"].(map[string]string)[sort] = order

	return queryBuilder
}

func (queryBuilder *QueryBuilder) GroupBy(groupBy string) *QueryBuilder {
	if groupBy == "" {
		return queryBuilder
	}

	queryBuilder.sqlParts["groupBy"] = groupBy

	return queryBuilder
}

func (queryBuilder *QueryBuilder) Having(having string) *QueryBuilder {
	queryBuilder.sqlParts["having"] = having

	return queryBuilder
}

func (queryBuilder *QueryBuilder) SetFirstResult(firstResult int) *QueryBuilder {
	queryBuilder.firstResult = firstResult

	return queryBuilder
}

func (queryBuilder *QueryBuilder) Where(condition string) *QueryBuilder {
	queryBuilder.sqlParts["where"] = condition

	return queryBuilder
}

func (queryBuilder *QueryBuilder) Join(join string, alias string, condition string) *QueryBuilder {
	return queryBuilder.InnerJoin(join, alias, condition)
}

func (queryBuilder *QueryBuilder) wrapJoin(join string, alias string, condition string) *QueryBuilder {
	queryBuilder.flag = ISJOIN
	queryBuilder.sqlParts["join"].(map[string]string)["joinTable"] = join
	queryBuilder.sqlParts["join"].(map[string]string)["joinAlias"] = alias
	queryBuilder.sqlParts["join"].(map[string]string)["joinCondition"] = condition

	return queryBuilder
}

func (queryBuilder *QueryBuilder) InnerJoin(join string, alias string, condition string) *QueryBuilder {
	queryBuilder.sqlParts["join"].(map[string]string)["joinType"] = INNER
	queryBuilder.wrapJoin(join, alias, condition)

	return queryBuilder
}

func (queryBuilder *QueryBuilder) LeftJoin(join string, alias string, condition string) *QueryBuilder {
	queryBuilder.sqlParts["join"].(map[string]string)["joinType"] = LEFT
	queryBuilder.wrapJoin(join, alias, condition)

	return queryBuilder
}

func (queryBuilder *QueryBuilder) RightJoin(join string, alias string, condition string) *QueryBuilder {
	queryBuilder.sqlParts["join"].(map[string]string)["joinType"] = RIGHT
	queryBuilder.wrapJoin(join, alias, condition)

	return queryBuilder
}

func (queryBuilder *QueryBuilder) GetFirstResult() int {
	return queryBuilder.firstResult
}

func (queryBuilder *QueryBuilder) SetMaxResults(maxResults int) *QueryBuilder {
	queryBuilder.maxResults = maxResults
	return queryBuilder
}

func (queryBuilder *QueryBuilder) getMaxResults() int {
	return queryBuilder.maxResults
}

func (queryBuilder *QueryBuilder) SetParameter(param interface{}) *QueryBuilder {
	queryBuilder.params = append(queryBuilder.params, param)
	return queryBuilder
}

func (queryBuilder *QueryBuilder) GetParameter() []interface{} {
	return queryBuilder.params
}

func (queryBuilder *QueryBuilder) GetSQL() string {
	sql := ""
	queryType := queryBuilder.queryType

	switch queryType {
	case INSERT:
		sql = queryBuilder.getSQLForInsert()
	case DELETE:
		sql = queryBuilder.getSQLForDelete()
		break
	case UPDATE:
		sql = queryBuilder.getSQLForUpdate()
		break
	case SELECT:
		sql = queryBuilder.getSQLForSelect()
		break
	default:
		sql = queryBuilder.getSQLForSelect()
		break
	}

	return sql
}

func (queryBuilder *QueryBuilder) setMapWrap(sql string) string {
	setMap := queryBuilder.sqlParts["set"].(map[string]string)

	for k, v := range setMap {
		sql += k + "=" + v + ","
	}

	sql = sql[:len(sql)-1]
	return sql
}

func (queryBuilder *QueryBuilder) getSQLForUpdate() string {
	sql := "UPDATE "

	fromMap := queryBuilder.sqlParts["from"].(map[string]string)

	table := fromMap["table"] + " " + fromMap["alias"]

	sql += table + " SET "

	sql = queryBuilder.setMapWrap(sql)

	if whereStr := queryBuilder.sqlParts["where"].(string); whereStr != "" {
		sql += " WHERE " + whereStr
	}

	return sql
}

func (queryBuilder *QueryBuilder) getSQLForJoins() string {
	sql := ""

	if queryBuilder.flag != ISJOIN {
		return ""
	}

	joinType := queryBuilder.sqlParts["join"].(map[string]string)["joinType"]
	joinTable := queryBuilder.sqlParts["join"].(map[string]string)["joinTable"]
	joinAlias := queryBuilder.sqlParts["join"].(map[string]string)["joinAlias"]
	joinCondition := queryBuilder.sqlParts["join"].(map[string]string)["joinCondition"]

	sql += " " + joinType + " JOIN " + joinTable + " " + joinAlias + " ON " + joinCondition

	return sql
}

func (queryBuilder *QueryBuilder) getFromClauses() string {
	tableSql := ""

	if fromMap := queryBuilder.sqlParts["from"].(map[string]string); fromMap != nil {
		tableSql = fromMap["table"] + " " + fromMap["alias"]
	}

	return tableSql + queryBuilder.getSQLForJoins()
}

func (queryBuilder *QueryBuilder) getSQLForSelect() string {
	sql := "SELECT "

	if selectStr := queryBuilder.sqlParts["select"].(string); selectStr != "" {
		sql += selectStr
	}

	sql += " FROM " + queryBuilder.getFromClauses()

	if whereStr := queryBuilder.sqlParts["where"].(string); whereStr != "" {
		sql += " WHERE " + whereStr
	}

	if groupByStr := queryBuilder.sqlParts["groupBy"].(string); groupByStr != "" {
		sql += " GROUP BY " + groupByStr
	}

	if havingStr := queryBuilder.sqlParts["having"].(string); havingStr != "" {
		sql += " HAVING " + havingStr
	}

	if queryBuilder.flag == ISSORT {
		sql += " ORDER BY "
		orderByMap := queryBuilder.sqlParts["orderBy"].(map[string]string)
		for sort, order := range orderByMap {
			sql += sort + " " + order + ","
		}

		sql = sql[:len(sql)-1]
	}

	if queryBuilder.isLimitQuery() {
		sql += " LIMIT " + strconv.Itoa(queryBuilder.firstResult) + "," + strconv.Itoa(queryBuilder.maxResults)
	}

	return sql
}

func (queryBuilder *QueryBuilder) getSQLForDelete() string {
	sql := "DELETE "

	if fromMap := queryBuilder.sqlParts["from"].(map[string]string); fromMap != nil {
		tableSql := fromMap["table"]

		sql += " FROM " + tableSql

		if whereStr := queryBuilder.sqlParts["where"].(string); whereStr != "" {
			sql += " WHERE " + whereStr
		}

		return sql
	}

	return sql
}

func (queryBuilder *QueryBuilder) getSQLForInsert() string {
	sql := "INSERT INTO "
	if fromMap := queryBuilder.sqlParts["from"].(map[string]string); fromMap != nil {
		tableSql := fromMap["table"]
		sql += tableSql + " SET "
		sql = queryBuilder.setMapWrap(sql)

		return sql
	}

	return sql
}

func (queryBuilder *QueryBuilder) isLimitQuery() bool {
	return queryBuilder.maxResults >= -1 || queryBuilder.firstResult >= 0
}

func (queryBuilder *QueryBuilder) executeQuery(query string) (map[int]map[string]string, error) {
	if queryBuilder.params != nil {
		rows, err := queryBuilder.database.Query(query, queryBuilder.params...)

		return getRowsMap(rows), err
	}

	rows, err := queryBuilder.database.Query(query, nil)

	return getRowsMap(rows), err
}

func getRowsMap(rows *sql.Rows) map[int]map[string]string {
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	columnPointers := make([]interface{}, count)

	result := map[int]map[string]string{}
	resultId := 0

	for rows.Next() {
		for i, _ := range columns {
			columnPointers[i] = &values[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			panic(err)
		}

		record := map[string]string{}

		for i, col := range columns {
			var v interface{}
			val := values[i]

			if str, ok := val.(string); ok {
				v = str
			} else {
				v = val

				switch v.(type) {
				case int64:
					res := strings.Split(fmt.Sprintf("%s", v), "=")
					resTmp := res[1]

					v = resTmp[:len(resTmp)-1]
				}
			}

			record[col] = fmt.Sprintf("%s", v)
		}

		result[resultId] = record
		resultId++
	}

	return result
}

func (queryBuilder *QueryBuilder) Query() (map[int]map[string]string, error) {
	if queryBuilder.queryType == SELECT {
		return queryBuilder.executeQuery(queryBuilder.GetSQL())
	}
	return nil, nil
}

func (queryBuilder *QueryBuilder) prepareAndExecute() sql.Result {
	stmt, err := queryBuilder.database.Prepare(queryBuilder.GetSQL())
	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec(queryBuilder.params...)
	if err != nil {
		panic(err)
	}

	return res
}

func (queryBuilder *QueryBuilder) PrepareAndExecute() (int64, error) {
	if queryBuilder.queryType == INSERT {
		res := queryBuilder.prepareAndExecute()
		return res.LastInsertId()
	}

	if queryBuilder.queryType == DELETE {
		res := queryBuilder.prepareAndExecute()
		return res.RowsAffected()
	}

	if queryBuilder.queryType == UPDATE {
		res := queryBuilder.prepareAndExecute()
		return res.RowsAffected()
	}

	return -1, nil
}

func (queryBuilder *QueryBuilder) Insert(table string) *QueryBuilder {
	queryBuilder.queryType = INSERT
	queryBuilder.setFromWrap(table, "")

	return queryBuilder
}

func (queryBuilder *QueryBuilder) Delete(table string) *QueryBuilder {
	queryBuilder.queryType = DELETE
	queryBuilder.setFromWrap(table, "")

	return queryBuilder
}

func (queryBuilder *QueryBuilder) setFromWrap(table string, alias string) {
	queryBuilder.sqlParts["from"] = map[string]string{
		"table": table,
		"alias": alias,
	}
}
