package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"revisitoukom/config" // Gunakan config untuk DB
	"revisitoukom/models"
	"gorm.io/gorm"
	"strconv" // Import untuk konversi uint ke string
)

// POST: Menyimpan jawaban user
func CreateExamQuestion(c *gin.Context) {
	var input struct {
		ExamID     int64   	`json:"exam_id" 	binding:"required"`
		QuestionID int64    `json:"question_id" binding:"required"`
		UserAnswer string  	`json:"user_answer" binding:"required"`
	}

	// Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah ExamID ada di database
	var exam models.Exam
	if err := config.DB.First(&exam, input.ExamID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exam dengan ID " + strconv.Itoa(int(input.ExamID)) + " tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat mencari data ujian: " + err.Error()})
		}
		return
	}

	// Cek apakah QuestionID ada di database
	var question models.Question
	if err := config.DB.First(&question, input.QuestionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Soal dengan ID " + strconv.Itoa(int(input.QuestionID)) + " tidak ditemukan"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan saat mencari data soal: " + err.Error()})
		}
		return
	}

	// Menyimpan jawaban user ke database
	examQuestion := models.ExamQuestion{
		ExamID:     input.ExamID,
		QuestionID: input.QuestionID,
		UserAnswer: input.UserAnswer,
	}

	if err := config.DB.Create(&examQuestion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan jawaban ujian: " + err.Error()})
		return
	}

	// Respons sukses
	c.JSON(http.StatusOK, gin.H{"data": examQuestion})
}

// GET: Mendapatkan semua jawaban user berdasarkan ExamID
func GetUserAnswers(c *gin.Context) {
	var input struct {
		ExamID int64 `json:"exam_id" binding:"required"`
	}

	// Validasi input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil semua jawaban berdasarkan ExamID
	var examQuestions []models.ExamQuestion
	if err := config.DB.Where("exam_id = ?", input.ExamID).Find(&examQuestions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil jawaban ujian: " + err.Error()})
		return
	}

	// Respons sukses
	c.JSON(http.StatusOK, gin.H{"data": examQuestions})
}
