package main

import (
	"fmt"

	"github.com/azuar4e/microservices-tfg/internal/controllers"
	"github.com/azuar4e/microservices-tfg/internal/handlers"
	"github.com/azuar4e/microservices-tfg/internal/initializers"
	"github.com/azuar4e/microservices-tfg/internal/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToPostgres()
	initializers.SyncDB()
	initializers.ConnectToDynamo()
	initializers.ConnectToSQS()
	initializers.ConnectToSNS()
}

func main() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	//registro y login de users
	v1.POST("/singin", controllers.SinginHandler)
	v1.POST("/singup", controllers.RegisterHandler)
	v1.Use(middleware.AuthMiddleware())
	//handlers de operaciones que requieren validacion
	v1.POST("/jobs", handlers.CreateJobHandler)
	v1.GET("/jobs", handlers.GetJobsHandler)
	v1.GET("/jobs/:id", handlers.GetJobByIdHandler)
	v1.GET("/validate", controllers.Validate)

	r.Run(":9090")
	fmt.Println("Escuchando en el puerto 9090")
}
