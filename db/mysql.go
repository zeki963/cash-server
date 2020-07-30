package db

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//SqlDB DB控制項
var SqlDB *sql.DB

//Dbcannot DB連線
func Dbcannot(DBInfo string) (err error) {
	SqlDB, err = sql.Open("mysql", DBInfo)
	if err != nil {
		return err
	}
	err = SqlDB.Ping()
	if err != nil {
		return err
	}
	SqlDB.SetConnMaxLifetime(100 * time.Second) //最大連接周期，超時close
	SqlDB.SetMaxOpenConns(100)                  //設置最大連接數
	return nil
}

//GetJSON  通用JSON
func GetJSON(sqlString string, taskID string) (string, error) {
	rows, err := SqlDB.Query(sqlString, taskID)
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
