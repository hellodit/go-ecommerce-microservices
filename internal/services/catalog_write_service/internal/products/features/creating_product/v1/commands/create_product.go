package createProductCommand

import (
	"time"

	"github.com/mehdihadeli/go-ecommerce-microservices/internal/pkg/utils/validator"
	uuid "github.com/satori/go.uuid"
)

// https://echo.labstack.com/guide/request/
// https://github.com/go-playground/validator

type CreateProduct struct {
	ProductID   uuid.UUID `validate:"required"`
	Name        string    `validate:"required,gte=0,lte=255"`
	Description string    `validate:"required,gte=0,lte=5000"`
	Price       float64   `validate:"required,gte=0"`
	CreatedAt   time.Time `validate:"required"`
}

func NewCreateProduct(name string, description string, price float64) (*CreateProduct, error) {
	command := &CreateProduct{
		ProductID:   uuid.NewV4(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   time.Now(),
	}
	err := validator.Validate(command)
	if err != nil {
		return nil, err
	}

	return command, nil
}
