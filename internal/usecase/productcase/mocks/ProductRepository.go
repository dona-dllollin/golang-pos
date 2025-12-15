package mocks

import (
	"context"

	"github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"
	mock "github.com/stretchr/testify/mock"
)

type ProductRepository struct {
	mock.Mock
}

// Create Product Mock
func (_m *ProductRepository) Create(ctx context.Context, p *productModel.Product) (int64, error) {
	args := _m.Called(ctx, p)
	return args.Get(0).(int64), args.Error(1)
}
