package cachekey

import "fmt"

func RoleByIDKey(ID int) string {
	return fmt.Sprintf("role:id:%d", ID)
}
