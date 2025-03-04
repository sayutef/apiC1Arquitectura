package controllers

import (
	"api/src/Users/application"
	"api/src/Users/domain"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateUserController struct {
	useCaseUpdate *application.UpdateProduct
	useCaseGet    *application.GetAllUser 
}

func NewUpdateUserController(useCaseUpdate *application.UpdateProduct) *UpdateUserController {
	return &UpdateUserController{useCaseUpdate: useCaseUpdate}
}

func (updateUser *UpdateUserController) Execute(ctx *gin.Context) {
	id := ctx.Param("id")

	var User domain.User

	if err := ctx.ShouldBindJSON(&User); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	UserId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return
	}

	err = updateUser.useCaseUpdate.Execute(int32(UserId), User.Name, User.Lastname)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message:": "Usuario actualizado"})
}

func (updateUser *UpdateUserController) ShortPolling(ctx *gin.Context) {
	id := ctx.Param("id")

	UserId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	ticker := time.NewTicker(5 * time.Second) // Consultar cada 5 segundos
	defer ticker.Stop()

	var lastUser domain.User

	for i := 0; i < 6; i++ { // Intenta máximo 30s (6 ciclos de 5s)
		select {
		case <-ticker.C:
			// Obtener la lista de usuarios actualizada
			users, err := updateUser.useCaseGet.Execute()
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Buscar el usuario específico
			var currentUser domain.User
			for _, u := range users {
				if u.Id == int32(UserId) {
					currentUser = u
					break
				}
			}

			// Si el usuario cambió, responder inmediatamente
			if currentUser != lastUser {
				ctx.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado", "user": currentUser})
				return
			}
		}
	}

	// Si no hubo cambios en 30s, devolver el último estado
	ctx.JSON(http.StatusOK, gin.H{"message": "Sin cambios", "user": lastUser})
}