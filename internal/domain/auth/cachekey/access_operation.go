package cachekey

import "fmt"

func OperationsByAccessRequirementCode(code string) string {
	return fmt.Sprintf("access_operations:access_requirement_code:%s", code)
}

func OperationsByPermissionID(permissionID int) string {
	return fmt.Sprintf("access_operations:permission_id:%d", permissionID)
}
