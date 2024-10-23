package storer

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, f func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")
	f(db, mock)
}
func TestCreateProduct(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.00,
		CountInStock: 10,
	}
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {

				// we tell the fake database to expect an INSERT action. This means we’re telling the database:
				// "You should be expecting us to add this product into the store’s database."
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "error occured creating product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnError(fmt.Errorf("error inserting product"))
				_, err := st.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "error occured getting last insert ID",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("error getting last insert ID")))
				_, err := st.CreateProduct(context.Background(), p)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestGetProduct(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.00,
		CountInStock: 10,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)
				// in the down line of we are telling it that we are expecting a SELECT action with the ID of 1. when that happened just return the rows we created above.
				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnRows(rows)
				// after the above line we are calling the GetProduct method from the storer and we are expecting it to return the product we created above.
				gp, err := st.GetProduct(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), gp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)

			},
		},
		{
			name: "error getting the product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error getting product"))
				_, err := st.GetProduct(context.Background(), 1)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}

}

func TestListProducts(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.00,
		CountInStock: 10,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)
				mock.ExpectQuery("SELECT * FROM products").WillReturnRows(rows)
				lp, err := st.ListProducts(context.Background())
				require.NoError(t, err)
				require.Len(t, lp, 1)
				require.Equal(t, int64(1), lp[0].ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "error listing products",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products").WillReturnError(fmt.Errorf("error listing products"))
				_, err := st.ListProducts(context.Background())
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}
func TestUpdateProduct(t *testing.T) {
	p := &Product{
		ID:           1,
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.00,
		CountInStock: 10,
	}
	np := &Product{
		ID:           1,
		Name:         "new test product",
		Image:        "new test.jpg",
		Category:     "new test category",
		Description:  "new test description",
		Rating:       5,
		NumReviews:   200,
		Price:        200.00,
		CountInStock: 20,
	}
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WithArgs(p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock).WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)

				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").WithArgs(np.Name, np.Image, np.Category, np.Description, np.Rating, np.NumReviews, np.Price, np.CountInStock, np.UpdatedAt, np.ID).WillReturnResult(sqlmock.NewResult(1, 1))
				up, err := st.UpdateProduct(context.Background(), np)
				require.NoError(t, err)
				require.Equal(t, int64(1), up.ID)
				require.Equal(t, np.Name, up.Name)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}},
		{
			name: "error updating product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE products SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=? WHERE id=?").WithArgs(np.Name, np.Image, np.Category, np.Description, np.Rating, np.NumReviews, np.Price, np.CountInStock, np.UpdatedAt, np.ID).WillReturnError(fmt.Errorf("error updating product"))
				_, err := st.UpdateProduct(context.Background(), np)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	p := &Product{
		ID:           1,
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       4,
		NumReviews:   100,
		Price:        100.00,
		CountInStock: 10,
	}
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?)").WithArgs(p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock).WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)

				mock.ExpectExec("DELETE FROM products WHERE id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				err = st.DeleteProduct(context.Background(), 1)
				require.NoError(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}},
		{
			name: "error deleting product",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error deleting product"))
				err := st.DeleteProduct(context.Background(), 1)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestCreateOrder(t *testing.T) {
	ois := []OrderItem{
		{
			Name:      "product 1",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
		},
		{
			Name:      "product 2",
			Quantity:  2,
			Image:     "test2.jpg",
			Price:     99.99,
			ProductID: 2,
		},
	}
	o := &Order{
		PaymentMethod: "test payment method",
		TaxPrice:      34,
		ShippingPrice: 123,
		TotalPrice:    1235,
		Items:         ois,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				// Start the transaction
				mock.ExpectBegin()

				// Mock order insertion
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WithArgs(o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice).WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock first order item insertion (order_id = 1)
				mock.ExpectExec("INSERT INTO order_items (name,quantity,image,price,product_id,order_id) VALUES (?,?,?,?,?,?)").WithArgs(ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, 1).WillReturnResult(sqlmock.NewResult(1, 1))

				// Mock second order item insertion (order_id = 1)
				mock.ExpectExec("INSERT INTO order_items (name,quantity,image,price,product_id,order_id) VALUES (?,?,?,?,?,?)").WithArgs(ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, 1).WillReturnResult(sqlmock.NewResult(2, 1))

				// Commit the transaction
				mock.ExpectCommit()

				// Execute the function
				cp, err := st.CreateOrder(context.Background(), o)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)

				// Verify all expectations were met
				er := mock.ExpectationsWereMet()
				require.NoError(t, er)
			},
		},
		{
			name: "failure_rollback",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				// Start the transaction
				mock.ExpectBegin()

				// Mock order insertion failure
				mock.ExpectExec("INSERT INTO orders (payment_method, tax_price, shipping_price, total_price) VALUES (?, ?, ?, ?)").WithArgs(
					o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice,
				).WillReturnError(fmt.Errorf("db error"))

				// Expect rollback
				mock.ExpectRollback()

				// Execute the function
				_, err := st.CreateOrder(context.Background(), o)

				// Expect an error and a rollback
				require.Error(t, err)
				require.Contains(t, err.Error(), "db error")

				// Verify all expectations were met
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestGetOrder(t *testing.T) {
	ois := []OrderItem{
		{
			ID:        1,
			Name:      "product 1",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
			OrderID:   1,
		},
		{
			ID:        2,
			Name:      "product 2",
			Quantity:  2,
			Image:     "test2.jpg",
			Price:     99.99,
			ProductID: 2,
			OrderID:   1,
		},
	}
	o := &Order{
		ID:            1,
		PaymentMethod: "test payment method",
		TaxPrice:      34,
		ShippingPrice: 123,
		TotalPrice:    1235,
		Items:         ois,
	}
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				// Mock the order query
				rows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(o.ID, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)
				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").WithArgs(1).WillReturnRows(rows)

				// Mock the order items query
				rows = sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(ois[0].ID, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, ois[0].OrderID)
				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").WithArgs(1).WillReturnRows(rows)
				gp, err := st.GetOrder(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), gp.ID)
			}},
		{
			name: "error getting order",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM orders WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error getting order"))
				_, err := st.GetOrder(context.Background(), 1)
				require.Error(t, err)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestListOrders(t *testing.T) {
	ois := []OrderItem{
		{
			ID:        1,
			Name:      "product 1",
			Quantity:  1,
			Image:     "test.jpg",
			Price:     99.99,
			ProductID: 1,
			OrderID:   1,
		},
		{
			ID:        2,
			Name:      "product 2",
			Quantity:  2,
			Image:     "test2.jpg",
			Price:     99.99,
			ProductID: 2,
			OrderID:   1,
		},
	}
	o := &Order{
		ID:            1,
		PaymentMethod: "test payment method",
		TaxPrice:      34,
		ShippingPrice: 123,
		TotalPrice:    1235,
		Items:         ois,
	}
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				// Mock the orders query
				rows := sqlmock.NewRows([]string{"id", "payment_method", "tax_price", "shipping_price", "total_price", "created_at", "updated_at"}).
					AddRow(o.ID, o.PaymentMethod, o.TaxPrice, o.ShippingPrice, o.TotalPrice, o.CreatedAt, o.UpdatedAt)
				mock.ExpectQuery("SELECT * FROM orders").WillReturnRows(rows)

				// Mock the order items query
				rows = sqlmock.NewRows([]string{"id", "name", "quantity", "image", "price", "product_id", "order_id"}).
					AddRow(ois[0].ID, ois[0].Name, ois[0].Quantity, ois[0].Image, ois[0].Price, ois[0].ProductID, ois[0].OrderID).
					AddRow(ois[1].ID, ois[1].Name, ois[1].Quantity, ois[1].Image, ois[1].Price, ois[1].ProductID, ois[1].OrderID)

				mock.ExpectQuery("SELECT * FROM order_items WHERE order_id=?").WithArgs(1).WillReturnRows(rows)
				lo, err := st.ListOrders(context.Background())
				require.NoError(t, err)
				require.Len(t, lo, 1)
				require.Equal(t, int64(1), lo[0].ID)
				require.Len(t, lo[0].Items, 2)
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestDeleteOrder(t *testing.T) {

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStorer, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				// Start the transaction
				mock.ExpectBegin()

				// Mock the order items deletion
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
				// Mock the order deletion
				mock.ExpectExec("DELETE FROM order WHERE id=?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

				// Commit the transaction
				mock.ExpectCommit()

				// Execute the function
				err := st.DeleteOrder(context.Background(), 1)
				require.NoError(t, err)

				// Verify all expectations were met
				er := mock.ExpectationsWereMet()
				require.NoError(t, er)
			}},
		{
			name: "failure_rollback",
			test: func(t *testing.T, st *MySQLStorer, mock sqlmock.Sqlmock) {
				// Start the transaction
				mock.ExpectBegin()

				// Mock order deletion failure
				mock.ExpectExec("DELETE FROM order_items WHERE order_id=?").WithArgs(1).WillReturnError(fmt.Errorf("error deleting order"))

				// Expect rollback
				mock.ExpectRollback()

				// Execute the function
				err := st.DeleteOrder(context.Background(), 1)

				// Expect an error and a rollback
				require.Error(t, err)

				// Verify all expectations were met
				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			}},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStorer(db)
				tc.test(t, st, mock)
			})
		})
	}
}
