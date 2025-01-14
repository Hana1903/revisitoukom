package controllers

import (
	"revisitoukom/config"
	"revisitoukom/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateExam(c *gin.Context) {
	var input struct {
		OrderID  int64   `json:"order_id"`
		PacketID int64   `json:"packet_id"`
		UserID   int64   `json:"user_id"`
		Score    float64 `json:"score"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil paket untuk mendapatkan durasi_ujian
	var packet models.Packet
	if err := config.DB.First(&packet, input.PacketID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found"})
		return
	}

	// Ambil durasi dari packet
	durationSeconds, err := strconv.Atoi(packet.DurationExam)
	if err != nil || durationSeconds <= 0 { // Jika ada error atau durasi ujian kurang dari atau sama dengan nol
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid duration format in packet", "details": err.Error()})
		return
	}

	// Mengatur start_at ke waktu saat ini dan hitung end_at
	startedAt := time.Now().In(time.Local)
	endedAt := startedAt.Add(time.Duration(durationSeconds) * time.Second)

	// Log untuk debugging
	log.Printf("StartedAt: %v, DurationSeconds: %v, EndedAt: %v", startedAt, durationSeconds, endedAt)

	// Create new Exam with calculated times
	exam := models.Exam{
		OrderID:   input.OrderID,
		PacketID:  input.PacketID,
		UserID:    input.UserID,
		Score:     input.Score,
		StartedAt: startedAt,
		EndedAt:   endedAt,
		CreatedAt: time.Now(),
	}

	// Menyimpan exam ke database
	if err := config.DB.Create(&exam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exam"})
		return
	}

	type Response struct {
		ID        uint     `json:"id"`
		OrderID   int64   `json:"order_id"`
		PacketID  int64   `json:"packet_id"`
		UserID    int64   `json:"user_id"`
		Score     float64 `json:"score"`
		StartedAt string  `json:"started_at"`
		EndedAt   string  `json:"ended_at"`
		CreatedAt string  `json:"created_at"`
		UpdatedAt string  `json:"updated_at"`
	}
	
	// Membuat response struct
	response := Response{
		ID:        exam.ID,
		OrderID:   exam.OrderID,
		PacketID:  exam.PacketID,
		UserID:    exam.UserID,
		Score:     exam.Score,
		StartedAt: exam.StartedAt.Format("2006-01-02 15:04:05"),
		EndedAt:   exam.EndedAt.Format("2006-01-02 15:04:05"),
		CreatedAt: exam.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: exam.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
	
	// Kirim response
	c.JSON(http.StatusCreated, response)

}

func GetRemainingTime(c *gin.Context) {
	var exam models.Exam
	if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	remainingTime := time.Until(exam.EndedAt)
	if remainingTime < 0 { // Jika remainingTime kurang dari nol
		c.JSON(http.StatusOK, gin.H{"remaining_time": "0s"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"remaining_time": remainingTime.String()})
}

// Mendapatkan semua Exam
func GetExams(c *gin.Context) {
	var exams []models.Exam
	if err := config.DB.Find(&exams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exams", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exams)
}

func GetExamByID(c *gin.Context) {
	type ExamResponse struct {
		ID         uint   `json:"id"`
		OrderID    int64  `json:"order_id"`
		PacketID   int64   `json:"packet_id"`
		UserID     int64   `json:"user_id"`
		Score      float64    `json:"score"`
		StartedAt  string `json:"started_at"`
		EndedAt    string `json:"ended_at"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	var exam models.Exam
	if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	response := ExamResponse{
		ID:         exam.ID,
		OrderID:    exam.OrderID,
		PacketID:   exam.PacketID,
		UserID:     exam.UserID,
		Score:      exam.Score,
		StartedAt:  exam.StartedAt.Format("2006-01-02 15:04:05"),
		EndedAt:    exam.EndedAt.Format("2006-01-02 15:04:05"),
		CreatedAt:  exam.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:  exam.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

func GetExamWithQuestions(c *gin.Context) {
	// Ambil packet_id dari query string
	packetID := c.Query("packet_id")
	if packetID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "packet_id is required"})
		return
	}

	// Cari semua pertanyaan berdasarkan packet_id
	var questions []models.Question
	if err := config.DB.Where("packet_id = ?", packetID).Find(&questions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve questions"})
		return
	}

	// Format data output
	var questionResponses []models.QuestionResponse
	for _, question := range questions {
		questionResponses = append(questionResponses, question.ToResponse())
	}

	// Response final
	c.JSON(http.StatusOK, gin.H{
		"packet_id": packetID,
		"questions": questionResponses,
	})
}

func UpdateExamScore(c *gin.Context) {
    var input struct {
        ExamID uint `json:"exam_id" binding:"required"`
    }

    // Validasi input JSON
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Ambil data ujian
    var exam models.Exam
    if err := config.DB.First(&exam, input.ExamID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ujian tidak ditemukan"})
        return
    }

    // Cek apakah ujian telah selesai
    if exam.EndedAt.After(time.Now()) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Ujian belum selesai, skor belum bisa dihitung"})
        return
    }

    // Ambil semua soal ujian terkait
    var examQuestions []models.ExamQuestion
    if err := config.DB.Where("exam_id = ?", input.ExamID).Find(&examQuestions).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil soal ujian"})
        return
    }

    var correctAnswersCount int64
    var answeredQuestionsCount int64

    // Hitung jumlah jawaban benar dan soal yang telah dijawab
    for _, eq := range examQuestions {
        if eq.UserAnswer != "" { // Hitung soal yang telah dijawab
            answeredQuestionsCount++
        }

        var question models.Question
        if err := config.DB.First(&question, eq.QuestionID).Error; err == nil {
            if eq.UserAnswer == question.CorrectAnswer {
                correctAnswersCount++
            }
        }
    }

    // Hitung skor akhir
    score := float64(correctAnswersCount) / float64(len(examQuestions)) * 100

    // Kirim respons
    c.JSON(http.StatusOK, gin.H{"data": gin.H{"score": score}})
}

// Delete an exam by ID
func DeleteExam(c *gin.Context) {
    var exam models.Exam
    if err := config.DB.First(&exam, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ujian tidak ditemukan"})
        return
    }

    // Cek waktu ujian
    if exam.EndedAt.Before(time.Now()) {
        c.JSON(http.StatusForbidden, gin.H{"error": "Ujian sudah selesai, tidak bisa dihapus"})
        return
    }

    if err := config.DB.Delete(&exam).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus ujian"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Ujian berhasil dihapus",
    })
}