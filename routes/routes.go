package routes

import (
	"revisitoukom/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

	//User
	r.GET("/users", controllers.GetUsers)
	r.GET("/users/:id", controllers.GetUserByID)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.PUT("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	//Packet
	r.GET("/packets", controllers.GetPackets)
	r.GET("/packets-detail/:id", controllers.GetPacketByID)
	r.GET("/packets/:packet_id/questions", controllers.GetQuestionsByPacketID)
	r.GET("/packets-purchased/:id", controllers.GetPacketsByUser)
	r.POST("/packets", controllers.CreatePacket)
	r.PUT("/packets/:id", controllers.UpdatePacket)
	r.DELETE("/packets/:id", controllers.DeletePacket)

	//Question
	r.GET("/questions", controllers.GetQuestions)
	r.GET("/questions/:id", controllers.GetQuestionByID)
	r.POST("/questions", controllers.CreateQuestion)
	r.PUT("/questions/:id", controllers.UpdateQuestion)
	r.DELETE("/questions/:id", controllers.DeleteQuestion)

	//Order
	r.GET("/orders", controllers.GetOrders)
	r.GET("/orders/:id", controllers.GetOrderByID)
	r.POST("/orders", controllers.CreateOrder)
	r.PUT("/orders/:id", controllers.UpdateOrder)
	r.DELETE("/orders/:id", controllers.DeleteOrder)

	//Exam
	r.GET("/exams", controllers.GetExams)
	r.GET("/exams/:id", controllers.GetExamByID)
	r.GET("/exams/:id/remaining-time", controllers.GetRemainingTime)
	r.GET("/exams/packet", controllers.GetExamWithQuestions)
	r.POST("/exams", controllers.CreateExam)
	r.PUT("/exams/update-score", controllers.UpdateExamScore)
	r.DELETE("/exams/:id", controllers.DeleteExam)

	//ExamQuestion
	r.POST("/exam_questions", controllers.CreateExamQuestion)
	r.POST("/exam_questions/user-answers", controllers.GetUserAnswers)
	return r
}