package models

import (
	"errors"
	"main_project/utils"

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

// Validates if required parameters sent through request body and returns approprite status code and message
func (product *Product) Validate() (map[string]interface{}, bool) {
	if product.Name == "" {
		return utils.Message(400, "Product name is required"), false
	}
	if product.SKUCode == "" {
		return utils.Message(400, "SKU Code is required"), false
	}
	return utils.Message(200, "success"), true
}

func (product *Product) Create() (*Product, map[string]interface{}, error) {
	if resp, ok := product.Validate(); !ok {
		return nil, resp, errors.New("Validation error")
	}

	existingProduct := &Product{}
	err := GetDB().Table("products").Where("sku_code = ?", &product.SKUCode).First(&existingProduct)
	if err.Error == nil {
		return nil, utils.Message(400, "Product with same sku code exists"), err.Error
	}

	err = GetDB().Table("products").Create(product)
	if err.Error != nil {
		return nil, utils.Message(400, "Error in creating product"), err.Error
	}

	return product, utils.Message(200, "Product created successfully"), nil

}

func GetProductById(id string) (*Product, map[string]interface{}, error) {
	product := &Product{}
	err := GetDB().Table("products").Where("id = ?", id).First(product).Error
	if err != nil {
		return nil, utils.Message(404, "Product does not exist"), err
	}
	return product, utils.Message(200, "Success"), nil
}

func GetProducts(offset, limit string) ([]*Product, map[string]interface{}, error) {
	products := make([]*Product, 0)
	err := GetDB().Table("products").Find(&products).Error
	if offset != "" && limit != "" {
		err = GetDB().Table("products").Offset(offset).Limit(limit).Find(&products).Error
	}

	if err != nil {
		return nil, utils.Message(400, "Error in retrieving products"), err
	}

	return products, utils.Message(200, "Success"), nil
}

func DeleteProduct(id string) (map[string]interface{}, error) {
	product := &Product{}
	err := GetDB().Table("products").Where("id = ?", id).First(product).Error
	if err != nil {
		return utils.Message(404, "Product does not exist"), err
	}
	err = GetDB().Table("products").Delete(&product).Error
	if err != nil {
		return utils.Message(400, "Error in deleting product"), err
	}
	return utils.Message(200, "Product deleted"), nil

}

func UpdateProduct(id string, product *Product) (*Product, map[string]interface{}, error) {
	if resp, ok := product.Validate(); !ok {
		return nil, resp, errors.New("Validation error")
	}

	existingProduct := &Product{}
	err := GetDB().Table("products").Where("id = ?", id).First(existingProduct).Error
	if err != nil {
		return nil, utils.Message(404, "Product does not exist"), err
	}
	existingProduct.SKUCode = product.SKUCode
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Color = product.Color
	existingProduct.Size = product.Size

	err = GetDB().Table("products").Save(&existingProduct).Error

	if err != nil {
		return nil, utils.Message(400, "Error in updating product"), err
	}
	return existingProduct, utils.Message(200, "Product Updated"), nil

}
