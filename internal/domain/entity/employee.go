package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Employee struct {
	ID          uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FullName    string          `gorm:"type:varchar(255);not null"`
	JobTitle    string          `gorm:"type:varchar(100);not null;index"`
	Country     string          `gorm:"type:varchar(100);not null;index"`
	GrossSalary decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt  `gorm:"index"`
}

func (Employee) TableName() string {
	return "employees"
}

func NewEmployee(fullName, jobTitle, country string, grossSalary decimal.Decimal) *Employee {
	return &Employee{
		ID:          uuid.New(),
		FullName:    fullName,
		JobTitle:    jobTitle,
		Country:     country,
		GrossSalary: grossSalary,
	}
}
