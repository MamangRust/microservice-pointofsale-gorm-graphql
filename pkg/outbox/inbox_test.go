package outbox

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// newMockDB builds a *gorm.DB backed by sqlmock so the consumer-inbox SQL can
// be exercised against a real GORM executor without a live PostgreSQL server.
func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm with sqlmock: %v", err)
	}
	return gormDB, mock
}

func TestReserveValidatesKeys(t *testing.T) {
	for name, tc := range map[string]struct {
		tx           *gorm.DB
		consumerName string
		eventKey     string
	}{
		"nil executor": {tx: nil, consumerName: "email-service-group", eventKey: "topic:evt-1"},
		"empty consumer": {tx: &gorm.DB{}, consumerName: "", eventKey: "topic:evt-1"},
		"empty event key": {tx: &gorm.DB{}, consumerName: "email-service-group", eventKey: ""},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, _, err := Reserve(context.Background(), tc.tx, tc.consumerName, tc.eventKey, "topic", 0, 1); !errors.Is(err, ErrInvalidInboxKey) {
				t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
			}
		})
	}
}

func TestMarkProcessedValidatesKeys(t *testing.T) {
	if err := MarkProcessed(context.Background(), nil, "email-service-group", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := MarkProcessed(context.Background(), &gorm.DB{}, "", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestReleaseValidatesKeysAndRecordsError(t *testing.T) {
	if err := Release(context.Background(), nil, "email-service-group", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}

	gormDB, mock := newMockDB(t)
	mock.ExpectExec(`(?s)UPDATE consumer_inbox.*last_error = \$1 WHERE consumer_name = \$2 AND event_key = \$3.*reservation_version = \$4`).
		WithArgs("smtp down", "email-service-group", "topic:evt-1", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := Release(context.Background(), gormDB, "email-service-group", "topic:evt-1", 1, errors.New("smtp down")); err != nil {
		t.Fatalf("Release returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestReserveUsesLeaseFencedUpsert(t *testing.T) {
	gormDB, mock := newMockDB(t)
	mock.ExpectQuery(`(?s)consumer_inbox.*ON CONFLICT.*lease_until <= current_timestamp`).
		WillReturnRows(sqlmock.NewRows([]string{"reservation_version", "reserved"}).AddRow(1, true))

	reserved, processed, version, err := Reserve(context.Background(), gormDB, "email-service-group", "topic:evt-1", "topic", 0, 1)
	if err != nil {
		t.Fatalf("Reserve returned error: %v", err)
	}
	if !reserved || processed || version != 1 {
		t.Fatalf("expected reserved=true processed=false version=1, got reserved=%v processed=%v version=%d", reserved, processed, version)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}

func TestReserveReportsAlreadyProcessed(t *testing.T) {
	gormDB, mock := newMockDB(t)
	mock.ExpectQuery(`(?s)consumer_inbox.*ON CONFLICT.*lease_until <= current_timestamp`).
		WillReturnRows(sqlmock.NewRows([]string{"reservation_version", "reserved"}).AddRow(0, false))
	mock.ExpectQuery(`(?s)SELECT status FROM consumer_inbox`).
		WillReturnRows(sqlmock.NewRows([]string{"status"}).AddRow("processed"))

	reserved, processed, _, err := Reserve(context.Background(), gormDB, "email-service-group", "topic:evt-1", "topic", 0, 1)
	if err != nil {
		t.Fatalf("Reserve returned error: %v", err)
	}
	if reserved || !processed {
		t.Fatalf("expected reserved=false processed=true, got reserved=%v processed=%v", reserved, processed)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sqlmock expectations: %v", err)
	}
}