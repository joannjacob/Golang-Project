package product

import (
	"context"
	"errors"
	"product_management/utils"

	"github.com/jinzhu/gorm"
)

type ProductRepositoryInterface interface {
	GetProducts(ctx context.Context, offset string, limit string) ([]*Product, map[string]interface{}, error)
	GetProductById(ctx context.Context, id int) (*Product, map[string]interface{}, error)
	CreateProduct(ctx context.Context, p *Product) (*Product, map[string]interface{}, error)
	UpdateProduct(ctx context.Context, id int, p *Product) (*Product, map[string]interface{}, error)
	DeleteProduct(ctx context.Context, id int) (map[string]interface{}, error)
}

func NewProductRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetProducts(ctx context.Context, offset string, limit string) ([]*Product, map[string]interface{}, error) {
	products := []*Product{}

	err := r.db.Find(&products).Error
	if offset != "" && limit != "" {
		err = r.db.Offset(offset).Limit(limit).Find(&products).Error
	}

	if err != nil {
		return nil, utils.Message(400, "Error in retrieving products"), err
	}

	return products, utils.Message(200, "Success"), nil

}

func (r *Repository) GetProductById(ctx context.Context, id int) (*Product, map[string]interface{}, error) {

	product := &Product{}
	err := r.db.Where("id = ?", id).First(product).Error
	if err != nil {
		return nil, utils.Message(404, "Product does not exist"), err
	}
	return product, utils.Message(200, "Success"), nil

}

func (product *Product) Validate() (map[string]interface{}, bool) {
	if product.Name == "" {
		return utils.Message(400, "Product name is required"), false
	}
	if product.SKUCode == "" {
		return utils.Message(400, "SKU Code is required"), false
	}
	return utils.Message(200, "success"), true
}

func (r *Repository) CreateProduct(ctx context.Context, product *Product) (*Product, map[string]interface{}, error) {
	if resp, ok := product.Validate(); !ok {
		return nil, resp, errors.New("Validation error")
	}

	existingProduct := &Product{}
	err := r.db.Where("sku_code = ?", &product.SKUCode).First(&existingProduct).Error
	if err == nil {
		return nil, utils.Message(400, "Product with same sku code exists"), errors.New("Product with same sku code exists")
	}

	err = r.db.Create(product).Error
	if err != nil {
		return nil, utils.Message(400, "Error in creating product"), err
	}

	return product, utils.Message(200, "Product created successfully"), nil

}

func (r *Repository) UpdateProduct(ctx context.Context, id int, product *Product) (*Product, map[string]interface{}, error) {
	if resp, ok := product.Validate(); !ok {
		return nil, resp, errors.New("Validation error")
	}

	existingProduct := &Product{}
	err := r.db.Where("id = ?", id).First(existingProduct).Error
	if err != nil {
		return nil, utils.Message(404, "Product does not exist"), err
	}
	existingProduct.SKUCode = product.SKUCode
	existingProduct.Name = product.Name
	existingProduct.Description = product.Description
	existingProduct.Color = product.Color
	existingProduct.Size = product.Size

	err = r.db.Save(&existingProduct).Error

	if err != nil {
		return nil, utils.Message(400, "Error in updating product"), err
	}
	return existingProduct, utils.Message(200, "Product Updated"), nil

}

func (r *Repository) DeleteProduct(ctx context.Context, id int) (map[string]interface{}, error) {
	product := &Product{}
	err := r.db.Where("id = ?", id).First(product).Error
	if err != nil {
		return utils.Message(404, "Product does not exist"), err
	}
	err = r.db.Delete(&product).Error
	if err != nil {
		return utils.Message(400, "Error in deleting product"), err
	}
	return utils.Message(200, "Product deleted"), nil

}
