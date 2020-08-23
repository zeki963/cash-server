package model

import "cash-server/db"

//   ---------------------------log_connect  表單相關-------------------------------------

//LogConnectAdd 寫入紀錄
func LogConnectAdd(p db.LogConnect) bool {
	dbrut := db.SQLDBX.Create(&p)
	return dbErrBool(dbrut)
}
