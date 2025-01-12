package models

import "time"

type ExamQuestion struct {
	ID     			uint   			`gorm:"primaryKey"`
	ExamID 			int64 			`json:"exam_id"`
	QuestionID 		int64			`json:"question_id"`
	UserAnswer		string			`json:"user_answer"`
	CreatedAt 		time.Time 		`json:"created_at"`
	UpdatedAt 		time.Time 		`json:"updated_at"`
}