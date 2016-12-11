package definition

import (
	"github.com/jinzhu/gorm"
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

var (
	PreloadShop            = CreatePreload(&Shop{})
	PreloadOrder           = CreatePreload(&Order{})
	PreloadCacheManageHash = CreatePreload(&CacheManageHash{})
)

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
