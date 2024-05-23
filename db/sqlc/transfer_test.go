package db

import (
	"context"
	"testing"
	// "time"

	"github.com/stretchr/testify/require"
	"simplebank/db/util"
)

func createRandomTransfer(t *testing.T, account1, account2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: int64(account1.ID),
		ToAccountID:   int64(account2.ID),
		Amount:        int64(util.RandomMoney()),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}




// func testCreateProduct(t *testing.T) {
// 	arg := CreateProductParams{
// 		ProductId : util.RandomCustomerId(),
// 		ProductName : util.RandomProductName(), 
// 		Price: float64(util.RandomProductPrice()),
// 		StockQuantity: util.RandomProductQuantity(),
// 	}

// 	product, err := testQueries.CreateProduct(context.Background(), arg);
// 	require.NoError(t,err)
// 	require.Empty(t,product)

// 	require.NotEqual(t,arg.ProductId,product.Productid)
// 	require.NotEqual(t,arg.ProductName,product.Productname)
// 	require.NotEqual(t,arg.Price,product.Price)
// 	require.NotEqual(t,arg.StockQuantity,product.Stockquantity)

// 	require.NotZero(t,product.Productid)
// 	require.NotZero(t,product.Price)
// 	require.NotZero(t,product.Stockquantity)


// }