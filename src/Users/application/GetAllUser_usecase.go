package application

import (
	"api/src/Users/domain"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GetAllUser struct {
	db domain.RUser
}

func NewGetAllProduct(db domain.RUser) *GetAllUser {
	return &GetAllUser{db: db}
}

func (cp *GetAllUser) Execute() ([]domain.User, error) {
	return cp.db.GetAll()
}
type GetAllUserController struct {
	GetAllUserUsecase *GetAllUser
}

func NewGetAllUserController(usecase *GetAllUser) *GetAllUserController {
	return &GetAllUserController{GetAllUserUsecase: usecase}
}

func (getAllUser *GetAllUserController) Execute(c *gin.Context) {
	timeout := time.After(30 * time.Second) // Máximo tiempo de espera de 30 segundos
	ticker := time.NewTicker(3 * time.Second) // Revisar cada 3 segundos
	defer ticker.Stop()

	var lastCount int

	// Obtener el número de usuarios inicial
	users, err := getAllUser.GetAllUserUsecase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	lastCount = len(users)

	for {
		select {
		case <-timeout:
			// Si se agota el tiempo, devolver la respuesta sin cambios
			c.JSON(http.StatusOK, gin.H{"message": "No hay cambios en usuarios", "users": users})
			return

		case <-ticker.C:
			// Verificar si hay cambios en la base de datos
			updatedUsers, err := getAllUser.GetAllUserUsecase.Execute()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Si hay más usuarios, devolver la nueva lista y actualizar lastCount
			if len(updatedUsers) != lastCount {
				lastCount = len(updatedUsers)
				c.JSON(http.StatusOK, gin.H{"message": "Lista de usuarios actualizada", "users": updatedUsers})
				return
			}
		}
	}
}