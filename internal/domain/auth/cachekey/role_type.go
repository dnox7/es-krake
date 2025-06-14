package cachekey

import "fmt"

func AllRoleTypesKey() string {
	return "role_type:all"
}

func RoleTypeByIDKey(ID int) string {
	return fmt.Sprintf("role_type:id:%d", ID)
}
