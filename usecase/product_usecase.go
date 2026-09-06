package usecase

import (
	"errors"
	"go-api/model"
	"go-api/repository"
)

var (
	ErrNotFound      = errors.New("todo not found")
	ErrNoFieldsPatch = errors.New("no fields provided for update")
)

// letra maiuscula para ser visivel fora do pacote
type ProductUsecase struct {
	// Repository
	repository repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return ProductUsecase{
		repository: repo,
	}
}

func (pu *ProductUsecase) GetProducts() ([]model.Product, error) {
	return pu.repository.GetProducts()
}

func (pu *ProductUsecase) CreateProduct(product model.Product) (model.Product, error) {
	productId, err := pu.repository.CreateProduct(product)

	if err != nil {
		return model.Product{}, err
	}

	product.ID = productId

	return product, nil
}

func (pu *ProductUsecase) GetProductById(id_product int) (*model.Product, error) {

	product, err := pu.repository.GetProductById(id_product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (pu *ProductUsecase) UpdateProduct(id_product int, product *model.Product) error {
	rowsAffected, err := pu.repository.UpdateProduct(id_product, product)
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (pu *ProductUsecase) ParcialUpdateProduct(id_product int, product *model.Product) error {
	rowsAffected, err := pu.repository.ParcialUpdateProduct(id_product, product)
	if err != nil {
		if err.Error() == "no fields to update" {
			return ErrNoFieldsPatch
		}
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
