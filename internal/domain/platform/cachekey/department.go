package cachekey

import "fmt"

func DepartmentByID(ID int) string {
	return fmt.Sprintf("department:id:%d", ID)
}
