package util

import (
	"time"
	"wails_study/project/logger"
)

func UnixToTimeString(timestamp int64) string {
	unix := time.Unix(int64(timestamp), 0)
	return unix.Format("2006-01-02 15:04:05")
}

func UnixToTimestampString(timestamp int64) string {
	unix := time.UnixMilli(timestamp)
	return unix.Format("2006-01-02 15:04:05.000")
}

func DiffMilliseconds(t1Str, t2Str string) int64 {
	t1, err := time.Parse("2006-01-02 15:04:05.000", t1Str)
	if err != nil {
		logger.Warn("DiffMilliseconds err: ", err.Error())
		return 0
	}

	t2, err := time.Parse("2006-01-02 15:04:05.000", t2Str)
	if err != nil {
		logger.Warn("DiffMilliseconds err: ", err.Error())
		return 0
	}

	diff := t2.Sub(t1)

	return diff.Milliseconds()
}
