package definition

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestDataSubItem struct {
	Id             uint
	TestDataItemId uint
}

type TestDataItem struct {
	Id         uint
	TestDataId uint
	Unique     TestDataSubItem
	SubItems   []TestDataSubItem
}

type TestData struct {
	Id     uint
	Unique TestDataItem
	Items  []TestDataItem
}

func TestNestFieldName(t *testing.T) {
	order := Order{}
	result := nestedFieldNames(&order)
	assert.Equal(t, result, map[string][]string{
		BelongsTo: []string{"Shop", "Items.User", "User"},
		HasMany:   []string{"Shop.Discounts", "Items", "Items.Cell"},
	})
	data := TestData{}
	result = nestedFieldNames(&data)
	assert.Equal(t, result, map[string][]string{
		HasOne:  []string{"Unique", "Unique.Unique", "Items.Unique"},
		HasMany: []string{"Unique.SubItems", "Items", "Items.SubItems"},
	})
}
