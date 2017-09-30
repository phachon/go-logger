package go_logger

import (
	"time"
)

type Misc struct {

}

func NewMisc() *Misc {
	return &Misc{}
}

//格式化 unix 时间戳
func (misc *Misc) FormatUnixTime(unixTime int64) string {
	tm := time.Unix(unixTime, 0)
	return tm.Format("2006-01-02 15:04:05")
}

//map Intersect
func (misc *Misc) MapIntersect(defaultMap map[string]interface{}, inputMap map[string]interface{}) map[string]interface{} {
	for key, _ := range defaultMap {
		inputValue, ok := inputMap[key]
		if !ok {
			continue
		}
		defaultMap[key] = inputValue
	}
	return defaultMap
}
