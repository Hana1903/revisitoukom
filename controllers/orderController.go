package controllers

import (
	"net/http"
	"revisitoukom/config"
	"revisitoukom/models"
	"time"
	"github.com/gin-gonic/gin"
)

// Membuat Order Baru
func CreateOrder(c *gin.Context) {
	var input struct {
		UserID    int    `json:"user_id" binding:"required"`
		PacketID  int    `json:"packet_id" binding:"required"`
		OrderDate string `json:"order_date" binding:"required"`
	}

	// Menyalin Input JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengubah tanggal dalam bentuk string ke time,Time
	orderDate, err := time.Parse("2006-01-02", input.OrderDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Atur waktu ke 00:00:00 untuk memastikan tidak ada waktu yang disertakan
	orderDate = time.Date(orderDate.Year(), orderDate.Month(), orderDate.Day(), 0, 0, 0, 0, time.UTC)

	// Menemukan Packet berdasarkan ID
	var packet models.Packet
	if err := config.DB.Where("id = ?", input.PacketID).First(&packet).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Packet not found"})
		return
	}

	// Create a new order
	order := models.Order{
		UserID:    input.UserID,
		PacketID:  input.PacketID,
		OrderDate: orderDate,
		Amount:    packet.Price,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	config.DB.Create(&order)
	// Respond dengan membuat order
	c.JSON(http.StatusOK, order)
}

// Mengambil Semua Data Oreder
func GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := config.DB.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// Ambil Data Order dengan ID
func GetOrderByID(c *gin.Context) {
	var order models.Order
	if err := config.DB.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// Update an Order by ID
func UpdateOrder(c *gin.Context) {
	var order models.Order
	if err := config.DB.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	var input struct {
		IDUser    int    `json:"id_user" 	binding:"required"`
		IDPacket  int    `json:"id_packet" 	binding:"required"`
		OrderDate string `json:"order_date" 	binding:"required"`
		Status    int    `json:"status" 		binding:"required"`
	}
	// Bind JSON input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Parsing tanggal bertipe string ke time.Time
	orderDate, err := time.Parse("2006-01-02", input.OrderDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}
	// Update Details Order
	order.UserID = input.IDUser
	order.PacketID = input.IDPacket
	order.OrderDate = orderDate
	order.Status = input.Status
	order.UpdatedAt = time.Now()
	config.DB.Save(&order)
	c.JSON(http.StatusOK, order)
}

// Menghapus Order dengan ID
func DeleteOrder(c *gin.Context) {
	var order models.Order
	if err := config.DB.Where("id = ?", c.Param("id")).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	// Delete order
	config.DB.Delete(&order)
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted successfully"})
}
