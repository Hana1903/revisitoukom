package main

import (
	"revisitoukom/config"
	"revisitoukom/models"
	"revisitoukom/routes"
)
func main()  {
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{}, &models.Packet{}, &models.Question{}, &models.Order{}, &models.Exam{}, &models.ExamQuestion{})
	router := routes.SetupRoutes()
	router.Run("0.0.0.0:8080")
}