package product

import (
	"context"

	"github.com/jinzhu/gorm"
)

type ProductUsecaseInterface interface {
	GetProducts(ctx context.Context, offset string, limit string) ([]*Product, map[string]interface{}, error)
	GetProductById(ctx context.Context, id int) (*Product, map[string]interface{}, error)
	CreateProduct(ctx context.Context, p *Product) (*Product, map[string]interface{}, error)
	UpdateProduct(ctx context.Context, id int, p *Product) (*Product, map[string]interface{}, error)
	DeleteProduct(ctx context.Context, id int) (map[string]interface{}, error)
}

func NewProductUsecase(db *gorm.DB, r ProductRepositoryInterface) *Usecase {
	return &Usecase{
		db:   db,
		repo: r,
	}
}

type Usecase struct {
	db   *gorm.DB
	repo ProductRepositoryInterface
}

func (u *Usecase) GetProducts(ctx context.Context, offset string, limit string) ([]*Product, map[string]interface{}, error) {

	products, msg, err := u.repo.GetProducts(ctx, offset, limit)
	if err != nil {
		return nil, msg, err
	}
	return products, msg, nil
}

func (u *Usecase) GetProductById(ctx context.Context, id int) (*Product, map[string]interface{}, error) {

	product, msg, err := u.repo.GetProductById(ctx, id)
	if err != nil {
		return nil, msg, err
	}
	return product, msg, nil
}

func (u *Usecase) CreateProduct(ctx context.Context, product *Product) (*Product, map[string]interface{}, error) {

	product, msg, err := u.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, msg, err
	}
	return product, msg, nil
}

func (u *Usecase) UpdateProduct(ctx context.Context, id int, product *Product) (*Product, map[string]interface{}, error) {

	product, msg, err := u.repo.UpdateProduct(ctx, id, product)
	if err != nil {
		return nil, msg, err
	}
	return product, msg, nil
}

func (u *Usecase) DeleteProduct(ctx context.Context, id int) (map[string]interface{}, error) {

	msg, err := u.repo.DeleteProduct(ctx, id)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
