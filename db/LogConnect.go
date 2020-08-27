package db

//LogConnect  連線紀錄
type LogConnect struct {
	Model
	StatusCode  int
	LatencyTime string `gorm:"type:varchar(15) COMMENT '運算時間'"`
	ClientIP    string `gorm:"type:varchar(15) COMMENT '來源IP'"`
	ReqMethod   string `gorm:"type:varchar(100) COMMENT '方式'"`
	ReqURL      string `gorm:"type:varchar(255) COMMENT '網址'"`
	Reqbody     string `gorm:"type:text "`
	Reqheader   string `gorm:"type:text"`
}
