package controllers

import (
	"api/src/Products/application"
	"api/src/Products/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateProductController struct {
	useCaseUpdate *application.UpdateProduct
	useCaseGet    *application.GetAllProduct // Asegúrate de que esté inicializado
}

func NewUpdateProductController(useCaseUpdate *application.UpdateProduct, useCaseGet *application.GetAllProduct) *UpdateProductController {
	return &UpdateProductController{
		useCaseUpdate: useCaseUpdate,
		useCaseGet:    useCaseGet, // Inicializar el useCaseGet
	}
}

func (updateProduct *UpdateProductController) ShortPolling(ctx *gin.Context) {
	id := ctx.Param("id")

	productId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	ticker := time.NewTicker(5 * time.Second) // Consultar cada 5 segundos
	defer ticker.Stop()

	var lastProduct domain.Product // Usamos esta variable para verificar cambios

	for i := 0; i < 6; i++ { // Intentar máximo 30s (6 ciclos de 5s)
		select {
		case <-ticker.C:
			// Obtener la lista de productos actualizada
			products, err := updateProduct.useCaseGet.Execute()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Buscar el producto específico
			var currentProduct domain.Product
			for _, p := range products {
				if p.Id == int32(productId) {
					currentProduct = p
					break
				}
			}

			// Si el producto cambió, responder inmediatamente
			if currentProduct != lastProduct {
				ctx.JSON(http.StatusOK, gin.H{"message": "Producto actualizado", "product": currentProduct})
				lastProduct = currentProduct // Actualizar el último producto
				return
			}
		}
	}

	// Si no hubo cambios en 30s, devolver el último estado
	ctx.JSON(http.StatusOK, gin.H{"message": "Sin cambios", "product": lastProduct})
}
