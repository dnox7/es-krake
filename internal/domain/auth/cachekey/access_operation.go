package cachekey

import "fmt"

func OperationsByAccessRequirementCode(code string) string {
	return fmt.Sprintf("access_operations:access_requirement_code:%s", code)
}
