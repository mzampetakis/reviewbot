package orders

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"regexp"
	"reviewbot/app"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func newTestDatabase(t *testing.T) (*sql.DB, *DatabaseRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error creating mock database: %v", err)
	}
	return db, NewDatabaseRepository(sqlx.NewDb(db, "mysql")), mock
}

// TestGetByUUID tests the GetByUUID function of the DatabaseRepository.
func TestGetOrderByUUID(t *testing.T) {
	// Arrange
	db, repo, mock := newTestDatabase(t)
	defer db.Close()

	orderUUID := uuid.New().String()
	customerUUID := uuid.New().String()
	rows := sqlmock.NewRows([]string{"uuid", "customer_uuid", "status", "placed_date"}).
		AddRow(orderUUID, customerUUID, app.OrderStatusPreparing, time.Now())
	// Add an expected query and its result to the mock database.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `uuid`, `customer_uuid`, `status`, " +
		"`placed_date` FROM `orders` WHERE (`uuid` = '" + orderUUID + "')")).
		WillReturnRows(rows)

	customerRows := sqlmock.NewRows([]string{"uuid", "first_name", "last_name", "email", "phone_number", "registration_date"}).
		AddRow(customerUUID, "first", "last", "e@mail.com", "+1234567890", time.Now())
	// Add an expected query and its result to the mock database.
	mock.ExpectQuery(regexp.QuoteMeta("SELECT `uuid`, `first_name`, `last_name`, `email`, `phone_number`, " +
		"`registration_date` FROM `customers` WHERE (`uuid` = '" + customerUUID + "')")).
		WillReturnRows(customerRows)

	// Act: Get the order by its UUID.
	order, err := repo.GetOrderByUUID(context.Background(), orderUUID)
	// Assert
	if err != nil {
		t.Fatalf("Error getting order by UUID: %v", err)
	}
	if order == nil {
		t.Fatalf("Order should not be nil")
	}
	if order.Status != app.OrderStatusPreparing {
		t.Fatalf("Order status mismatch: got %s, want %s", order.Status, app.OrderStatusPreparing)
	}
	if order.Customer.UUID != customerUUID {
		t.Fatalf("Order customer UUID mismatch: got %s, want %s", customerUUID, order.Customer.UUID)
	}
	if order.Customer.FirstName != "first" {
		t.Fatalf("Order customer first name mismatch: got %s, want %s", "first", order.Customer.FirstName)
	}
	if order.Customer.LastName != "last" {
		t.Fatalf("Order customer last name mismatch: got %s, want %s", "last", order.Customer.LastName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Unfulfilled expectations: %v", err)
	}
}
