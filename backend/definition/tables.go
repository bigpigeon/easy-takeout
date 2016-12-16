package definition

import (
	"reflect"
	"time"

	"github.com/jinzhu/gorm"
)

func nestedFieldNames(v interface{}) map[string][]string {
	result := map[string][]string{}
	scope := &gorm.Scope{Value: v}
	ms := scope.GetModelStruct()
	for _, field := range ms.StructFields {
		if field.IsNormal == false && field.Relationship != nil {
			kind := field.Relationship.Kind
			if result[kind] == nil {
				result[kind] = []string{}
			}
			result[kind] = append(result[kind], field.Name)
			subResult := nestedFieldNames(reflect.New(field.Struct.Type).Interface())
			for subkind, items := range subResult {
				if result[subkind] == nil {
					result[subkind] = []string{}
				}
				for _, subitem := range items {
					result[subkind] = append(result[subkind], field.Name+"."+subitem)
				}
			}
		}

	}
	return result
}

func kindNestSelect(kinds map[string][]string, selects []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, sl := range selects {
			for _, item := range kinds[sl] {
				db = db.Preload(item)
			}
		}
		return db
	}
}

func NestedSelect(table interface{}, selects []string) func(db *gorm.DB) *gorm.DB {
	ret := nestedFieldNames(table)
	return kindNestSelect(ret, selects)
}

//db struct
type (
	User struct {
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `sql:"index"`
		Name      string     `gorm:"type:varchar(256);primary_key"`
		Pass      string
	}

	Discount struct {
		Full   int
		Reduce int
		ShopId uint `gorm:"index"`
	}

	Shop struct {
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time `sql:"index"`

		Address    string `gorm:"type:varchar(1024);primary_key"`
		BeginPrice int
		BeginCost  int
		Discounts  []Discount `gorm:"ForeignKey:ShopId`
	}

	UserItemCell struct {
		ID         uint   `gorm:"primary_key"`
		UserItemId uint   `gorm:"unique_index:idx_item_name"`
		Name       string `gorm:"type:varchar(256);unique_index:idx_item_name"`
		Num        uint
	}

	UserItem struct {
		ID      uint           `gorm:"primary_key"`
		OrderId uint           `gorm:"unique_index:idx_order_user;index"`
		Cell    []UserItemCell `gorm:"ForeignKey:UserItemId"`

		User     User   `gorm:"ForeignKey:UserName"`
		UserName string `gorm:"type:varchar(256);unique_index:idx_order_user"`
	}

	Order struct {
		gorm.Model
		ShopAddr string `gorm:"type:varchar(1024);index"`
		Shop     Shop   `gorm:"ForeignKey:ShopAddr"`

		Tag   string     `gorm:"type:varchar(64);unique_index"`
		Items []UserItem `gorm:"ForeignKey:OrderId"`

		User     *User   `gorm:"ForeignKey:UserName"`
		UserName *string `gorm:"type:varchar(256);`

		EndAt *time.Time
	}

	CacheManageHash struct {
		gorm.Model
		Key    string                 `gorm:"type:varchar(256);unique_index"`
		Fields []CacheManageHashField `gorm:"ForeignKey:HashId"`
	}
	CacheManageHashField struct {
		ID     uint   `gorm:"primary_key"`
		HashId uint   `gorm:"index:idx_hash_field"`
		Field  string `gorm:"type:varchar(100);index:idx_hash_field"`
		Value  string
	}
)

func dropAndCreate(tables ...interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.DropTableIfExists(tables...).CreateTable(tables...)
	}
}

func Migrate(db *gorm.DB) error {
	db.Scopes(dropAndCreate(
		&User{}, &Discount{}, &Shop{}, &UserItemCell{}, &UserItem{}, &Order{},
		&CacheManageHash{}, &CacheManageHashField{},
	))
	return nil
}

func (cell *UserItemCell) Incr(value int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(cell).Update("num", gorm.Expr("num + ?", value))
	}
}
