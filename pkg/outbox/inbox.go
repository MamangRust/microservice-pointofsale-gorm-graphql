package outbox

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrInvalidInboxKey = errors.New("invalid consumer inbox key")

// ConsumerInbox is the durable deduplication contract used by Kafka handlers
// (Phase 3 — Durable Idempotency). It replaces in-memory-only deduplication:
// reservations survive consumer restarts and rebalances, so at-least-once
// redelivery cannot send the same email twice.
type ConsumerInbox interface {
	Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error)
	MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error
	Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error
}

// PostgresInbox adapts the schema-agnostic consumer-inbox SQL to a live
// *gorm.DB. It keeps the email service's durable-idempotency behavior
// (Phase 3) without depending on a generated schema package.
type PostgresInbox struct {
	db *gorm.DB
}

func NewInbox(db *gorm.DB) (*PostgresInbox, error) {
	if db == nil {
		return nil, errors.New("consumer inbox requires a non-nil gorm.DB")
	}
	return &PostgresInbox{db: db}, nil
}

func (p *PostgresInbox) Reserve(ctx context.Context, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	return Reserve(ctx, p.db, consumerName, eventKey, topic, partition, offset)
}

func (p *PostgresInbox) MarkProcessed(ctx context.Context, consumerName, eventKey string, reservationVersion int64) error {
	return MarkProcessed(ctx, p.db, consumerName, eventKey, reservationVersion)
}

func (p *PostgresInbox) Release(ctx context.Context, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	return Release(ctx, p.db, consumerName, eventKey, reservationVersion, processingErr)
}

const reserveConsumerInboxSQL = `
WITH reserved AS (
    INSERT INTO consumer_inbox (
        consumer_name, event_key, topic, partition_id, message_offset,
        status, attempts, reservation_version, lease_until, last_error, processed_at
    )
    VALUES (?, ?, ?, ?, ?, 'processing', 1, 1,
            current_timestamp + interval '1 minute', '', NULL)
    ON CONFLICT (consumer_name, event_key) DO UPDATE
    SET status = 'processing',
        attempts = consumer_inbox.attempts + 1,
        reservation_version = consumer_inbox.reservation_version + 1,
        lease_until = current_timestamp + interval '1 minute',
        last_error = '',
        topic = EXCLUDED.topic,
        partition_id = EXCLUDED.partition_id,
        message_offset = EXCLUDED.message_offset
    WHERE consumer_inbox.status <> 'processed'
      AND consumer_inbox.lease_until <= current_timestamp
    RETURNING reservation_version
)
SELECT
    COALESCE((SELECT reservation_version FROM reserved), 0) AS reservation_version,
    EXISTS (SELECT 1 FROM reserved) AS reserved
`

const consumerInboxStatusSQL = `
SELECT status FROM consumer_inbox WHERE consumer_name = ? AND event_key = ?
`

const markConsumerInboxProcessedSQL = `
UPDATE consumer_inbox
SET status = 'processed', processed_at = current_timestamp,
    lease_until = current_timestamp, last_error = ''
WHERE consumer_name = ? AND event_key = ?
  AND status = 'processing' AND reservation_version = ?
`

const releaseConsumerInboxSQL = `
UPDATE consumer_inbox
SET status = 'pending', lease_until = current_timestamp,
    last_error = ?
WHERE consumer_name = ? AND event_key = ?
  AND status = 'processing' AND reservation_version = ?
`

// Reserve claims an event for a consumer. It returns false when the event was
// already processed. An expired processing lease may be reclaimed after a
// consumer crashes.
func Reserve(ctx context.Context, db *gorm.DB, consumerName, eventKey, topic string, partition int32, offset int64) (bool, bool, int64, error) {
	if db == nil || consumerName == "" || eventKey == "" {
		return false, false, 0, ErrInvalidInboxKey
	}
	var result struct {
		ReservationVersion int64
		Reserved           bool
	}
	err := db.WithContext(ctx).Raw(reserveConsumerInboxSQL, consumerName, eventKey, topic, partition, offset).
		Scan(&result).Error
	if err != nil {
		return false, false, 0, err
	}

	// Determine whether this event was already processed: the upsert only
	// returns a row when it actually (re)reserved the key.
	processed := false
	if !result.Reserved {
		var status string
		if scanErr := db.WithContext(ctx).Raw(consumerInboxStatusSQL, consumerName, eventKey).Scan(&status).Error; scanErr == nil {
			processed = status == "processed"
		}
	}
	return result.Reserved, processed, result.ReservationVersion, nil
}

func MarkProcessed(ctx context.Context, db *gorm.DB, consumerName, eventKey string, reservationVersion int64) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	return db.WithContext(ctx).Exec(markConsumerInboxProcessedSQL, consumerName, eventKey, reservationVersion).Error
}

func Release(ctx context.Context, db *gorm.DB, consumerName, eventKey string, reservationVersion int64, processingErr error) error {
	if db == nil || consumerName == "" || eventKey == "" {
		return ErrInvalidInboxKey
	}
	lastError := "consumer processing failed"
	if processingErr != nil {
		lastError = processingErr.Error()
	}
	return db.WithContext(ctx).Exec(releaseConsumerInboxSQL, lastError, consumerName, eventKey, reservationVersion).Error
}