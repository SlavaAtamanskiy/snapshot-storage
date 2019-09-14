package utils

import (
	"github.com/lithammer/shortuuid"
	"strings"
	"time"
)

func GenerateDocLink(collectionName string) string {

	curTime := GetLocalTime().Format("2006-01-02T15:04:05")
	return strings.Join([]string{collectionName, curTime, shortuuid.New()}, "_")

}

func GetLocalTime() time.Time {

	return time.Now().Local()

}
