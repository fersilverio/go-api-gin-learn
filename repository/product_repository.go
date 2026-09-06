package repository

import (
	"database/sql"
	"fmt"
	"go-api/model"
	"strings"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {

	query := "SELECT id, product_name, price FROM product"

	rows, err := pr.connection.Query(query)

	if err != nil {
		fmt.Println(err)
		return []model.Product{}, err
	}

	var productList []model.Product
	var productObj model.Product

	for rows.Next() {
		err = rows.Scan(
			&productObj.ID,
			&productObj.Name,
			&productObj.Price,
		)

		if err != nil {
			fmt.Println(err)
			return []model.Product{}, err
		}

		productList = append(productList, productObj)
	}

	rows.Close()

	return productList, nil

}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	var id int

	query, err := pr.connection.Prepare("INSERT INTO product (product_name, price) VALUES ($1, $2) RETURNING id")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()

	return id, nil
}

func (pr *ProductRepository) GetProductById(id_product int) (*model.Product, error) {

	query, err := pr.connection.Prepare("SELECT * FROM product WHERE id = $1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var produto model.Product

	err = query.QueryRow(id_product).Scan(
		&produto.ID,
		&produto.Name,
		&produto.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	query.Close()
	return &produto, nil
}

func (pr *ProductRepository) UpdateProduct(id_product int, product *model.Product) (int64, error) {
	query, err := pr.connection.Prepare("UPDATE product SET product_name = $1, price = $2 WHERE id = $3 RETURNING *")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	result, err := query.Exec(product.Name, product.Price, id_product)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()
	return rowsAffected, nil
}

func (pr *ProductRepository) ParcialUpdateProduct(id_product int, product *model.Product) (int64, error) {
	var setClauses []string
	var args []interface{}

	paramCount := 1

	if product.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("product_name = $%d", paramCount))
		args = append(args, product.Name)
		paramCount++
	}

	if product.Price != 0 {
		setClauses = append(setClauses, fmt.Sprintf("price = $%d", paramCount))
		args = append(args, product.Price)
		paramCount++
	}

	if len(setClauses) == 0 {
		return 0, fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE product SET %s WHERE id = $%d", strings.Join(setClauses, ", "), paramCount)
	args = append(args, id_product)

	stmt, err := pr.connection.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return rowsAffected, nil
}
