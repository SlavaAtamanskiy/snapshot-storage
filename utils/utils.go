package utils

import (
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/shortuuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GenerateDocLink(collectionName string) string {

	curTime := GetLocalTime().Format("2006-01-02T15:04:05")
	return strings.Join([]string{collectionName, curTime, shortuuid.New()}, "_")

}

func GetLocalTime() time.Time {

	return time.Now().Local()

}

func GetLimitOffset(limit, offset string) (limitInt, offsetInt int) {

	if len(limit) == 0 {
		limitInt = 0
	} else {
		limitInt, _ = strconv.Atoi(limit)
	}

	if len(offset) == 0 {
		offsetInt = 0
	} else {
		offsetInt, _ = strconv.Atoi(offset)
	}

	return limitInt, offsetInt

}

func HandleError(r http.ResponseWriter, e error) {

	var sts int
	var message string
	jStr := `
		{
      "success": false,
			"msg": "&msg"
    }
	`

	switch status.Code(e) {
	case codes.NotFound:
		sts = http.StatusNotFound
		message = "No data found"
	default:
		sts = http.StatusInternalServerError
		message = e.Error()
	}

	data := []byte(strings.Replace(jStr, "&msg", message, 1))

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(sts)
	r.Write(data)

}

func CustomError(r http.ResponseWriter, msg string, sts int) {

	jStr := `
		{
      "success": false,
			"msg": "&msg"
    }
	`
	data := []byte(strings.Replace(jStr, "&msg", msg, 1))

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(sts)
	r.Write(data)

}

func ResponseOk(r http.ResponseWriter, data []byte) {

	var jsonData []byte
	def := `
		{
      "success": true
    }
	`
	if data == nil {
		jsonData = []byte(def)
	} else {
		jsonData = data
	}

	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(http.StatusOK)
	r.Write(jsonData)

}

func DecapitalizeStruct(d interface{}) map[string]interface{} {

	var out = make(map[string]interface{})

	fields := reflect.TypeOf(d)
	values := reflect.ValueOf(d)

	num := fields.NumField()

	for i := 0; i < num; i++ {

		field := fields.Field(i)
		value := values.Field(i)

		key := toSnakeCase(field.Name)

		switch value.Kind() {
		case reflect.String:
			out[key] = value.String()
		case reflect.Int:
			out[key] = strconv.FormatInt(value.Int(), 10)
		case reflect.Int32:
			out[key] = strconv.FormatInt(value.Int(), 10)
		case reflect.Int64:
			out[key] = strconv.FormatInt(value.Int(), 10)
		case reflect.Bool:
			out[key] = value.Bool()
		default:
			unknownValue := value.Interface()
			switch unknownValue.(type) {
			case time.Time:
				out[key] = unknownValue
			}
		}
	}

	return out

}

func toSnakeCase(str string) string {

	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

	return strings.ToLower(snake)

}
