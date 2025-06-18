package model

import (
	"context"
	"time"
)

type Sequence struct {
	ID                   uint64    `json:"id"`
	Name                 string    `json:"name"`
	OpenTrackingEnabled  bool      `json:"openTrackingEnabled"`
	ClickTrackingEnabled bool      `json:"clickTrackingEnabled"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	Steps                []*Step   `json:"steps"`
}

type Step struct {
	ID         uint64    `json:"id"`
	SequenceID uint64    `json:"-"`
	Subject    string    `json:"subject"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type SequenceStore interface {
	// Creates a sequence with steps.
	CreateSequence(ctx context.Context, sequence *Sequence) error
	// Fetch a sequence by ID.
	FetchSequence(ctx context.Context, id uint64) (*Sequence, error)
	// Update sequence open or click tracking.
	UpdateSequence(ctx context.Context, id uint64, sequence *Sequence) error
	// Update a sequence step (new subject or content).
	UpdateStep(ctx context.Context, id uint64, step *Step) error
	// Delete a sequence step.
	DeleteStep(ctx context.Context, id uint64, step *Step) error
}
