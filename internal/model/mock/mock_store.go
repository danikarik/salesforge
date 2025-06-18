package mock

import (
	"context"

	"github.com/danikarik/salesforge/internal/model"
	mocky "github.com/stretchr/testify/mock"
)

var _ model.SequenceStore = (*MockStore)(nil)

type MockStore struct {
	mocky.Mock
}

func (m *MockStore) CreateSequence(ctx context.Context, sequence *model.Sequence) error {
	args := m.Called(ctx, sequence)
	return args.Error(0)
}

func (m *MockStore) FetchSequence(ctx context.Context, id uint64) (*model.Sequence, error) {
	args := m.Called(ctx, id)

	var sequence *model.Sequence
	if args.Get(0) != nil {
		sequence = args.Get(0).(*model.Sequence)
	}

	return sequence, args.Error(1)
}

func (m *MockStore) UpdateSequence(ctx context.Context, id uint64, sequence *model.Sequence) error {
	args := m.Called(ctx, id, sequence)
	return args.Error(0)
}

func (m *MockStore) UpdateStep(ctx context.Context, id uint64, step *model.Step) error {
	args := m.Called(ctx, id, step)
	return args.Error(0)
}

func (m *MockStore) DeleteStep(ctx context.Context, id uint64, step *model.Step) error {
	args := m.Called(ctx, id, step)
	return args.Error(0)
}
