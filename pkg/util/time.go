package util

import (
	"fmt"
	"strings"
	"time"
)

var ()

//GetUTCTime 獲取時間
func GetUTCTime() time.Time {
	t := time.Now()
	locallocation, err := time.LoadLocation("UTC")
	if err != nil {
		fmt.Println(err)
	}
	timeUTC := t.In(locallocation)
	return timeUTC
}

//ParseCustomTimeString ParseCustomTimeString
func ParseCustomTimeString(timestring string) time.Time {
	//"2015-03-10 23:58:23 UTC"
	strarray := strings.Split(timestring, "UTC")
	var year, mon, day, hh, mm, ss int
	fmt.Sscanf(strarray[0], "%d-%d-%d %d:%d:%d ", &year, &mon, &day, &hh, &mm, &ss)
	timestringtoparse := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", year, mon, day, hh, mm, ss)
	urlcreatetime, _ := time.Parse(time.RFC3339, timestringtoparse)
	return urlcreatetime
}

//GETNowsqltime 他媽的給我SQL 專用的Now現在時間戳記
func GETNowsqltime() string {
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")
	return ts
}

//ShortDur 型態轉換 time.Duration 2 string
func ShortDur(d time.Duration) string {
	s := d.String()
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return s
}
