package productcase

import (
	"context"
	"testing"

	"github.com/dona-dllollin/belajar-clean-arch/internal/domain/productModel"
	"github.com/dona-dllollin/belajar-clean-arch/internal/usecase/productcase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProductUseCase_CreateProduct_Success(t *testing.T) {
	repo := new(mocks.ProductRepository)

	uc := &ProductUseCase{
		productRepo: repo,
	}

	product := &productModel.Product{
		Name:        "Produk Test",
		Description: "Desc",
		Status:      "active",
	}

	repo.
		On("Create", mock.Anything, product).
		Return(int64(1), nil).
		Once()

	id, err := uc.CreateProduct(context.Background(), product)

	require.NoError(t, err)
	require.NotNil(t, id)
	assert.Equal(t, int64(1), *id)

	repo.AssertExpectations(t)
}

// func TestProductUseCase_CreateProduct_ValidationError(t *testing.T) {
// 	repo := new(mocks.ProductRepository)
// 	validator := validation.New() // atau mock kalau mau

// 	uc := &ProductUseCase{
// 		productRepo: repo,
// 		validator:   validator,
// 	}

// 	product := &productModel.Product{} // invalid

// 	id, err := uc.CreateProduct(context.Background(), product)

// 	assert.Error(t, err)
// 	assert.Nil(t, id)

// 	// repo tidak boleh dipanggil
// 	repo.AssertNotCalled(t, "Create")
// }
