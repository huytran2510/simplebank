package db

import (
	"context"
	"simplebank/db/util"
	// _"simplebank/db/util"
	"testing"

	"github.com/stretchr/testify/require"
)

// var testQueries *Queries


func testCreateProduct(t *testing.T) {
	arg := CreateProductParams{
		ProductId : util.RandomCustomerId(),
		ProductName : util.RandomProductName(), 
		Price: float64(util.RandomProductPrice()),
		StockQuantity: util.RandomProductQuantity(),
	}

	product, err := testQueries.CreateProduct(context.Background(), arg);
	require.NoError(t,err)
	require.Empty(t,product)

	require.NotEqual(t,arg.ProductId,product.Productid)
	require.NotEqual(t,arg.ProductName,product.Productname)
	require.NotEqual(t,arg.Price,product.Price)
	require.NotEqual(t,arg.StockQuantity,product.Stockquantity)

	require.NotZero(t,product.Productid)
	require.NotZero(t,product.Price)
	require.NotZero(t,product.Stockquantity)


}

func testDeleteProduct(t *testing.T) {
	arg := CreateProductParams{
		ProductId:     util.RandomCustomerId(),
		ProductName:   util.RandomProductName(),
		Price:         float64(util.RandomProductPrice()),
		StockQuantity: util.RandomProductQuantity(),
	}

	product, err := testQueries.CreateProduct(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, product)

	deletedProduct, err := testQueries.DeleteProduct(context.Background(), product.Productid)
	require.NoError(t, err)

	require.Empty(t, deletedProduct)

	_, err = testQueries.GetProduct(context.Background(), deletedProduct.Productid)
	require.Error(t, err)
}