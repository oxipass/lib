package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"github.com/oxipass/oxilib/consts"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func GenerateTempFilename() string {
	fullPathDBFile := os.TempDir()
	tempFileName := GenerateRandomString(4) + consts.CTempDBFile
	if strings.HasSuffix(fullPathDBFile, "/") {
		fullPathDBFile += tempFileName
	} else {
		fullPathDBFile += "/" + tempFileName
	}
	return fullPathDBFile
}

func GenerateRandomString(length int) string {
	rb := make([]byte, length)
	_, err := rand.Read(rb)

	if err != nil {
		return ""
	}
	b64 := base64.URLEncoding.EncodeToString(rb)
	finalLen := length
	if utf8.RuneCountInString(b64) < finalLen {
		finalLen = utf8.RuneCountInString(b64)
	}
	return b64[0:finalLen]
}

func PrepareTimeForDb(timeIn time.Time) string {
	return timeIn.Format(consts.CDbFormat)
}

//func timeFromDb(dtStr string) (time.Time, error) {
//	dbTime, errDbTime := time.Parse(constDbFormat, dtStr)
//	if errDbTime != nil {
//		return dbTime, formError(BSERR00005ParseTimeFailed, errDbTime.Error())
//	}
//	return dbTime, nil
//}

// EncodeJSON - transforms any structure with JSON param into JSON string
func EncodeJSON(preparedStruct interface{}) (string, error) {
	b, err := json.Marshal(preparedStruct)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

// DecodeJSON - transforms json strin into the appropriate structure
func DecodeJSON(jsonStr string, outStruct interface{}) error {
	data := []byte(jsonStr)
	return json.Unmarshal(data, &outStruct)

}
