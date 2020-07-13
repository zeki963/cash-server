package main

var (
	mysqlip         = "34.92.36.107"
	mysqlport       = "32001"
	mysqldb         = "zor"
	mysqluser       = "root"
	mysqlpassword   = "cqig7777"
	mysqlparameters string
	dblink          = mysqluser + ":" + mysqlpassword + "@tcp(" + mysqlip + ":" + mysqlport + ")/" + mysqldb + "?" + mysqlparameters
	dblinkTEST      = "root:cqig7777@tcp(34.92.36.107:32001)/zor"
)
