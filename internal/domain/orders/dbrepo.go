package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"reviewbot/app"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

// CustomerStore represents a customer entity at the Database.
type CustomerStore struct {
	UUID             string
	FirstName        string
	LastName         string
	Email            string
	PhoneNumber      string
	RegistrationDate time.Time
}

// OrderStore represents an order entity at the Database.
type OrderStore struct {
	UUID         string
	OrderUUID    string
	CustomerUUID string
	Status       string
	PlacedDate   time.Time
}

// OrderProductStore represents an order product entity at the Database.
type OrderProductStore struct {
	UUID         string
	OrderUUID    string
	CustomerUUID string
	Items        int
	ProductUUID  string
}

// ProductStore represents a product entity at the Database.
type ProductStore struct {
	UUID               string
	Name               string
	Description        string
	Image              string
	AvailabilityStatus string
	AvailableItems     int
}

// DatabaseRepository implements the OrdersRepository interface.
type DatabaseRepository struct {
	db *sqlx.DB
}

// NewDatabaseRepository returns a new DatabaseRepository.
func NewDatabaseRepository(db *sqlx.DB) *DatabaseRepository {
	return &DatabaseRepository{
		db: db,
	}
}

// OrderStoreToOrder converts an OrderStore object to an app.Order
func (ds *DatabaseRepository) OrderStoreToOrder(orderStore OrderStore, customerStore CustomerStore) app.Order {
	return app.Order{
		UUID:       orderStore.UUID,
		Customer:   ds.CustomerStoreToCustomer(customerStore),
		Status:     app.OrderStatus(orderStore.Status),
		PlacedDate: orderStore.PlacedDate,
	}
}

// CustomerStoreToCustomer converts a CustomerStore object to an app.Customer
func (ds *DatabaseRepository) CustomerStoreToCustomer(customerStore CustomerStore) app.Customer {
	return app.Customer{
		UUID:             customerStore.UUID,
		FirstName:        customerStore.FirstName,
		LastName:         customerStore.LastName,
		Email:            customerStore.Email,
		PhoneNumber:      customerStore.PhoneNumber,
		RegistrationDate: customerStore.RegistrationDate,
	}
}

// OrderProductStoreToOrderProduct converts an OrderProductStore object to an app.OrderProduct
func (ds *DatabaseRepository) OrderProductStoreToOrderProduct(orderProductStore OrderProductStore,
	productStore ProductStore) app.OrderProduct {
	return app.OrderProduct{
		UUID:        orderProductStore.UUID,
		OrderUUID:   orderProductStore.OrderUUID,
		ProductUUID: orderProductStore.CustomerUUID,
		Items:       orderProductStore.Items,
		Product:     ds.ProductStoreToProduct(productStore),
	}
}

// ProductStoreToProduct converts ProductStore object to an app.Product
func (ds *DatabaseRepository) ProductStoreToProduct(productStore ProductStore) app.Product {
	return app.Product{
		UUID:               productStore.UUID,
		Name:               productStore.Name,
		Description:        productStore.Description,
		Image:              productStore.Image,
		AvailabilityStatus: productStore.AvailabilityStatus,
		AvailableItems:     productStore.AvailableItems,
	}
}

// GetOrderByUUID retrieves from storage an order by its UUID.
func (ds *DatabaseRepository) GetOrderByUUID(ctx context.Context, orderUUID string) (*app.Order, error) {
	dialect := goqu.Dialect("mysql")
	sqlQuery, _, err := dialect.Select("uuid", "customer_uuid", "status", "placed_date").
		From("orders").Where(goqu.C("uuid").Eq(orderUUID)).ToSQL()
	if err != nil {
		return nil, app.NewError("Error while preparing querying for order",
			fmt.Errorf("get by uuid: %w", err))
	}

	orderStore := new(OrderStore)
	err = ds.db.QueryRowContext(ctx, sqlQuery).Scan(&orderStore.UUID, &orderStore.CustomerUUID,
		&orderStore.Status, &orderStore.PlacedDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NewError("Order does not exist", app.ErrNoRecords)
		}
		return nil, app.NewError("Error while getting order", fmt.Errorf("get by uuid: %w", err))
	}

	sqlQuery, _, err = dialect.Select("uuid", "first_name", "last_name", "email", "phone_number", "registration_date").
		From("customers").Where(goqu.C("uuid").Eq(orderStore.CustomerUUID)).ToSQL()
	if err != nil {
		return nil, app.NewError("Error while preparing querying for customer",
			fmt.Errorf("get by uuid: %w", err))
	}

	customerStore := new(CustomerStore)
	err = ds.db.QueryRowContext(ctx, sqlQuery).Scan(&customerStore.UUID, &customerStore.FirstName,
		&customerStore.LastName, &customerStore.Email, &customerStore.PhoneNumber, &customerStore.RegistrationDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NewError("Order's customer does not exist", app.ErrNoRecords)
		}
		return nil, app.NewError("Error while getting order's customer", fmt.Errorf("get by uuid: %w", err))
	}

	order := ds.OrderStoreToOrder(*orderStore, *customerStore)
	return &order, nil
}

// GetOrderProductsByOrderUUID retrieves from storage an order's products by order's UUID.
func (ds *DatabaseRepository) GetOrderProductsByOrderUUID(ctx context.Context, orderUUID string) ([]app.OrderProduct, error) {
	dialect := goqu.Dialect("mysql")
	sqlQuery, _, err := dialect.Select("uuid", "order_uuid", "product_uuid", "items").
		From("order_products").Where(goqu.C("order_uuid").Eq(orderUUID)).ToSQL()
	if err != nil {
		return nil, app.NewError("Error while preparing querying for order products",
			fmt.Errorf("get by uuid: %w", err))
	}

	rows, err := ds.db.QueryContext(ctx, sqlQuery)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NewError("Order Products does not exist", app.ErrNoRecords)
		}
		return nil, app.NewError("Error while getting order products", fmt.Errorf("get by uuid: %w", err))
	}
	defer rows.Close()
	orderProducts := []app.OrderProduct{}
	for rows.Next() {
		var orderProduct OrderProductStore
		if err := rows.Scan(&orderProduct.UUID, &orderProduct.OrderUUID, &orderProduct.ProductUUID, &orderProduct.Items); err != nil {
			return nil, app.NewError("Error while reading order products", fmt.Errorf("get all: %w", err))
		}

		// Fetch ProductStore item
		sqlQuery, _, err := dialect.Select("uuid", "name", "status", "description", "image", "availability_status", "available_items").
			From("products").Where(goqu.C("uuid").Eq(orderProduct.ProductUUID)).ToSQL()
		if err != nil {
			return nil, app.NewError("Error while preparing querying for product",
				fmt.Errorf("get by uuid: %w", err))
		}

		productStore := new(ProductStore)
		err = ds.db.QueryRowContext(ctx, sqlQuery).Scan(&productStore.UUID, &productStore.Name,
			&productStore.Description, &productStore.Image, &productStore.AvailabilityStatus, &productStore.AvailableItems)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, app.NewError("Product does not exist", app.ErrNoRecords)
			}
			return nil, app.NewError("Error while getting product", fmt.Errorf("get by uuid: %w", err))
		}

		orderProducts = append(orderProducts, ds.OrderProductStoreToOrderProduct(orderProduct, *productStore))
	}

	return orderProducts, nil
}

// GetProductByUUID retrieves from storage a product by its UUID.
func (ds *DatabaseRepository) GetProductByUUID(ctx context.Context, productUUID string) (*app.Product, error) {
	dialect := goqu.Dialect("mysql")
	sqlQuery, _, err := dialect.Select("uuid", "name", "status", "description", "image", "availability_status", "available_items").
		From("products").Where(goqu.C("uuid").Eq(productUUID)).ToSQL()
	if err != nil {
		return nil, app.NewError("Error while preparing querying for product",
			fmt.Errorf("get by uuid: %w", err))
	}

	productStore := new(ProductStore)
	err = ds.db.QueryRowContext(ctx, sqlQuery).Scan(&productStore.UUID, &productStore.Name,
		&productStore.Description, &productStore.Image, &productStore.AvailabilityStatus, &productStore.AvailableItems)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, app.NewError("Product does not exist", app.ErrNoRecords)
		}
		return nil, app.NewError("Error while getting product", fmt.Errorf("get by uuid: %w", err))
	}

	product := ds.ProductStoreToProduct(*productStore)
	return &product, nil
}
