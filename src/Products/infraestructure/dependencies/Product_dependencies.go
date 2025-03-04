package dependencies

import (
	"api/src/Products/application"
	"api/src/Products/infraestructure"
	"api/src/Products/infraestructure/controllers"
	"api/src/core"
	"database/sql"
	"fmt"
)

var (
	mySQL infraestructure.MySQL
	db    *sql.DB
)

func Init() {
	db, err := core.ConnectToDB()

	if err != nil {
		fmt.Println("server error")
		return
	}

	mySQL = *infraestructure.NewMySQL(db)

}
func CloseDB() {
	if db != nil {
		db.Close()
		fmt.Println("Conexión a la base de datos cerrada.")
	}
}
func GetCreateProductController() *controllers.CreateProductController {
	caseCreateProduct := application.NewCreateProduct(&mySQL)
	return controllers.NewCreateProductController(caseCreateProduct)
}

func GetGetAllProductController() *controllers.GetAllProductController {
	caseGetAllProduct := application.NewGetAllProduct(&mySQL)
	return controllers.NewGetAllProductController(*caseGetAllProduct)
}

func GetDeleteProductController() *controllers.DeleteProductController {
	caseDeleteProduct := application.NewDeleteProduct(&mySQL)
	return controllers.NewDeleteProductController(caseDeleteProduct)
}

func GetUpdateProductController() *controllers.UpdateProductController {
	caseUpdateProduct := application.NewUpdateProduct(&mySQL)
	caseGetAllProduct := application.NewGetAllProduct(&mySQL) // Obtén el caso de uso GetAllProduct
	return controllers.NewUpdateProductController(caseUpdateProduct, caseGetAllProduct) // Pasa ambos casos de uso al controlador
}
func GetByIdProductController() *controllers.GetByIdProductController {
	caseGetByIdProduct := application.NewGetByIdProduct(&mySQL)
	return controllers.NewGetByIdProductController(caseGetByIdProduct)
}
