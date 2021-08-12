package product

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	SKUCode     string `gorm:"column:sku_code;not null;unique"`
	Name        string `gorm:"column:name;not null"`
	Description string `gorm:"column:description"`
	Color       string `gorm:"column:color"`
	Size        string `gorm:"column:size"`
}
