package db

//LogConnect  連線紀錄
type LogConnect struct {
	Model
	StatusCode  int
	LatencyTime string
	ClientIP    string `gorm:"type:varchar(15)"`
	ReqMethod   string `gorm:"type:varchar(100)"`
	ReqURL      string `gorm:"type:varchar(100)"`
	Reqbody     string `gorm:"type:varchar(255)"`
	Reqheader   string `gorm:"type:text"`
}
