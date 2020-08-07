package db

import (
	"database/sql"
	"encoding/json"
	"time"

	//註解
	_ "github.com/go-sql-driver/mysql"
)

//SQLDB DB控制項
var SQLDB *sql.DB

//Dbcannot DB連線
func Dbcannot(DBInfo string) (err error) {
	SQLDB, err = sql.Open("mysql", DBInfo)
	if err != nil {
		return err
	}
	err = SQLDB.Ping()
	if err != nil {
		return err
	}
	SQLDB.SetConnMaxLifetime(100 * time.Second) //最大連接周期，超時close
	SQLDB.SetMaxOpenConns(100)                  //設置最大連接數
	return nil
}

//GetJSON  SQL通用JSON
func GetJSON(sqlString string, taskID string) (string, error) {
	rows, err := SQLDB.Query(sqlString, taskID)
	if err != nil {
		return "", err
	}
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
