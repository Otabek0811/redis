package postgres

import (
	"app/api/models"
	"context"
	"testing"
)

func TestCreateProduct(t *testing.T) {

	tests := []struct {
		Name    string
		Input   *models.CreateProduct
		Output  string
		WantErr bool
	}{
		{
			Name: "Product Case 1",
			Input: &models.CreateProduct{
				Name:  "Iphone 10",
				Price: 3_000_000,
			},
			WantErr: false,
		},
		{
			Name: "Product Case 2",
			Input: &models.CreateProduct{
				Name:       "Iphone 14",
				Price:      12_000_000,
				CategoryId: "ab48d6e0-a705-414d-910b-ff4bfda27ccd",
			},
			WantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			id, err := productTestRepo.Create(context.Background(), test.Input)

			if test.WantErr || err != nil {
				t.Errorf("%s: got: %v", test.Name, err)
				return
			}

			if id == "" {
				t.Errorf("%s: got: %v", test.Name, err)
				return
			}
		})
	}
}

func TestGetByIDProduct(t *testing.T) {

	tests := []struct {
		Name    string
		Input   *models.ProductPrimaryKey
		Output  *models.Product
		WantErr bool
	}{
		{
			Name:  "Product Case 1",
			Input: &models.ProductPrimaryKey{Id: "e90a7772-1e3d-4094-8e6b-720b62966556"},
			Output: &models.Product{
				Id:         "e90a7772-1e3d-4094-8e6b-720b62966556",
				Name:       "Iphone 14",
				Price:      12_000_000,
				CategoryId: "ab48d6e0-a705-414d-910b-ff4bfda27ccd",
			},
			WantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			product, err := productTestRepo.GetByID(context.Background(), test.Input)

			if test.WantErr || err != nil {
				t.Errorf("%s: got: %v", test.Name, err)
				return
			}

			if product.Name != test.Output.Name ||
				product.Price != test.Output.Price ||
				product.Id != test.Output.Id ||
				product.CategoryId != test.Output.CategoryId {

				t.Errorf("%s: got: %+v, expected: %+v\n", test.Name, *product, *&test.Output)
				return
			}
		})
	}
}
