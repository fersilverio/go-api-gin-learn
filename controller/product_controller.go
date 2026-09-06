package controller

import (
	"go-api/model"
	"go-api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type productController struct {
	// Usecases
	productUsecase usecase.ProductUsecase
}

func NewProductController(usecase usecase.ProductUsecase) productController {
	return productController{
		productUsecase: usecase,
	}
}

func (p *productController) GetProducts(ctx *gin.Context) {
	products, err := p.productUsecase.GetProducts()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (p *productController) CreateProduct(ctx *gin.Context) {
	var product model.Product

	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedProduct, err := p.productUsecase.CreateProduct(product)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedProduct)
}

func (p *productController) GetProductById(ctx *gin.Context) {

	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Id do produto nao pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um numero",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUsecase.GetProductById(productId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := model.Response{
			Message: "Produto nao foi encontrado na base de dados",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *productController) UpdateProduct(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Id do produto nao pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um numero",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var product model.Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = p.productUsecase.UpdateProduct(productId, &product)

	if err != nil {
		if err == usecase.ErrNotFound {
			response := model.Response{
				Message: "Produto nao foi encontrado na base de dados",
			}
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Produto atualizado com sucesso"})

}

func (p *productController) ParcialUpdateProduct(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Id do produto nao pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um numero",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var product model.Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = p.productUsecase.ParcialUpdateProduct(productId, &product)

	if err != nil {
		if err == usecase.ErrNotFound {
			response := model.Response{
				Message: "Produto nao foi encontrado na base de dados",
			}
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		if err == usecase.ErrNoFieldsPatch {
			response := model.Response{
				Message: "Nenhum campo fornecido para atualização",
			}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Produto atualizado com sucesso"})
}

func (p *productController) DeleteProduct(ctx *gin.Context) {
	id := ctx.Param("productId")

	if id == "" {
		response := model.Response{
			Message: "Id do produto nao pode ser nulo",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	productId, err := strconv.Atoi(id)

	if err != nil {
		response := model.Response{
			Message: "Id do produto precisa ser um numero",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = p.productUsecase.DeleteProduct(productId)

	if err != nil {
		if err == usecase.ErrNotFound {
			response := model.Response{
				Message: "Produto nao foi encontrado na base de dados",
			}
			ctx.JSON(http.StatusNotFound, response)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Produto deletado com sucesso"})
}
