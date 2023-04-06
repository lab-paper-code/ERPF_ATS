package types

import (
	"encoding/base64"
	"fmt"
)

func GetBase64String(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func CheckAuthKey(id, pwd, authKey string) bool {
	expectedAuthKey := GetBase64String(fmt.Sprintf("%s:%s", id, pwd))
	return expectedAuthKey == authKey
}
