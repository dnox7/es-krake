package repository

import (
	"database/sql/driver"
	"pech/es-krake/pkg/utils"
	"testing"
)

func TestAttributeRepository(t *testing.T) {
	expectedDummyRow := utils.CreateSingleMockRow(map[string]driver.Value{})
	mockAttributeEntity := map[string]driver.Value{
		"id":            1,
		"name":          "size",
		"description":   "dummy description",
		"is_required":   false,
		"display_order": 1,
	}

	inputMap := map[string]interface{}{
		"id": 1,
	}

	t.Run("AttributeRepo TakeByID", func(t *testing.T) {
	})
}
