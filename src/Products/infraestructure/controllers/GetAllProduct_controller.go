package controllers

import (
	"api/src/Products/application"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetAllProductController struct {
	GetAllUseCase *application.GetAllProduct
}

func NewGetAllProductController(GetAllUseCase application.GetAllProduct) *GetAllProductController {
	return &GetAllProductController{GetAllUseCase: &GetAllUseCase}
}


func (getallproduct *GetAllProductController) Execute(ctx *gin.Context) {
	timeout := time.After(30 * time.Second) // Tiempo máximo de espera
	ticker := time.NewTicker(3 * time.Second) // Revisar cambios cada 3s

	defer ticker.Stop()

	// Obtener la lista inicial de productos
	products, err := getallproduct.GetAllUseCase.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lastCount := len(products)

	for {
		select {
		case <-timeout:
		
			ctx.JSON(http.StatusOK, gin.H{"message": "No hay cambios en productos", "Productos": products})
			return

		case <-ticker.C:
			
			updatedProducts, err := getallproduct.GetAllUseCase.Execute()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			
			if len(updatedProducts) != lastCount {
				ctx.JSON(http.StatusOK, gin.H{"message": "Lista de productos actualizada", "Productos": updatedProducts})
				return
			}
		}
	}
}
