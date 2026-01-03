package model

import "gorm.io/gorm"

type Sample struct {
	ID   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`
	*gorm.DB
}

func (Sample) TableName() string {
	return "samples"
}

func NewSample(db *gorm.DB) *Sample {
	return &Sample{
		DB: db,
	}
}

func (s *Sample) Migrate() error {
	return s.DB.AutoMigrate(&Sample{})
}

func (s *Sample) Create(name string) error {
	sample := &Sample{
		Name: name,
	}
	return s.DB.Create(sample).Error
}
