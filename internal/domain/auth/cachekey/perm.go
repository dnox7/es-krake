package cachekey

import "fmt"

func PermByIDKey(ID int) string {
	return fmt.Sprintf("permission:id:%d", ID)
}

func PermsByRoleID(roleID int) string {
	return fmt.Sprintf("permissions:role_id:%d", roleID)
}
