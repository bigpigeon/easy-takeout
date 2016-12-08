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
		gorm.Model
		Name string `gorm:"unique_index"`
		Pass string
	}

	Discount struct {
		Full   int
		Reduce int
		ShopId uint `gorm:"index"`
	}

	Shop struct {
		gorm.Model
		Address    string `gorm:"type:varchar(256);unique_index"`
		BeginPrice int
		BeginCost  int
		Discounts  []Discount `gorm:"ForeignKey:ShopId`
		Orders     []Order    `gorm:"ForeignKey:ShopId"`
	}
	OrderItem struct {
		gorm.Model
		OrderId uint   `gorm:"index"`
		Name    string `gorm:"index:idx_name_item"`
		Item    string `gorm:"index:idx_name_item"`
		Num     int
	}

	Order struct {
		gorm.Model
		ShopId uint        `gorm:"index"`
		Tag    string      `gorm:"type:varchar(64);unique_index"`
		Items  []OrderItem `gorm:"ForeignKey:OrderId"`
		User   User
		UserId uint
		EndAt  *time.Time
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

/*
Gorm field relations Kind
*/
const (
	HasMany    = "has_many"
	BelongsTo  = "belongs_to"
	HasOne     = "has_one"
	ManyToMany = "many_to_many"
)

type Preload struct {
	All        func(db *gorm.DB) *gorm.DB
	HasMany    func(db *gorm.DB) *gorm.DB
	HasOne     func(db *gorm.DB) *gorm.DB
	BelongsTo  func(db *gorm.DB) *gorm.DB
	ManyToMany func(db *gorm.DB) *gorm.DB
}

func CreatePreload(table interface{}) *Preload {
	ret := nestedFieldNames(table)
	preload := Preload{}
	preload.All = kindNestSelect(ret, []string{HasMany, BelongsTo, HasOne, ManyToMany})
	preload.HasMany = kindNestSelect(ret, []string{HasMany})
	preload.HasOne = kindNestSelect(ret, []string{HasOne})
	preload.BelongsTo = kindNestSelect(ret, []string{BelongsTo})
	preload.ManyToMany = kindNestSelect(ret, []string{ManyToMany})
	return &preload
}

var (
	PreloadShop            = CreatePreload(&Shop{})
	PreloadOrder           = CreatePreload(&Order{})
	PreloadCacheManageHash = CreatePreload(&CacheManageHash{})
)

func Migrate(db *gorm.DB) error {
	db.CreateTable(
		&User{}, &Discount{}, &Shop{}, &OrderItem{}, &Order{},
		&CacheManageHash{}, &CacheManageHashField{},
	)
	return nil
}
