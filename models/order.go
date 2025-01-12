package models

import (
	"time"
)

type Order struct {
	ID        	uint       		`gorm:"primaryKey"`
	UserID    	int       		`gorm:"not null"`
	PacketID  	int       		`gorm:"not null"`
	Status    	int       		`gorm:"not null"`
	OrderDate 	time.Time 		`gorm:"type:date;not null"` 
	Amount    	float64   		`gorm:"type:decimal(10,2)"`
	CreatedAt 	time.Time 		`gorm:"autoCreateTime"`
	UpdatedAt 	time.Time 		`gorm:"autoUpdateTime"`
}
