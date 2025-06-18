package pg_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/danikarik/salesforge/internal/model"
	"github.com/danikarik/salesforge/internal/model/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	connString, ok := os.LookupEnv("TEST_DATABASE_URL")
	if !ok {
		os.Exit(0)
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal(err)
	}
	testPool = pool

	code := m.Run()
	testPool.Close()
	os.Exit(code)
}

func cleanDB(ctx context.Context) {
	testPool.Exec(ctx, "DELETE FROM steps")
	testPool.Exec(ctx, "DELETE FROM sequences")
}

func assertDifference(t *testing.T, tableName string, diff int64, fn func()) {
	t.Helper()

	var before int64
	err := testPool.QueryRow(t.Context(), "SELECT COUNT(*) FROM "+tableName).Scan(&before)
	require.NoError(t, err)

	fn()

	var after int64
	err = testPool.QueryRow(t.Context(), "SELECT COUNT(*) FROM "+tableName).Scan(&after)
	require.NoError(t, err)

	require.Equal(t, before+diff, after, "Expected %d rows in %s table, got %d", before+diff, tableName, after)
}

var testSequence = &model.Sequence{
	Name:                 "Test Sequence",
	OpenTrackingEnabled:  true,
	ClickTrackingEnabled: true,
	Steps: []*model.Step{
		{
			Subject: "Step 1 Subject",
			Content: "Step 1 Content",
		},
		{
			Subject: "Step 2 Subject",
			Content: "Step 2 Content",
		},
	},
}

func TestCreateSequence(t *testing.T) {
	ctx := t.Context()
	cleanDB(ctx)

	store, err := pg.NewStore(testPool)
	require.NoError(t, err)

	assertDifference(t, "sequences", 1, func() {
		err := store.CreateSequence(ctx, testSequence)
		require.NoError(t, err)
	})
}

func TestFetchSequence(t *testing.T) {
	ctx := t.Context()
	cleanDB(ctx)

	store, err := pg.NewStore(testPool)
	require.NoError(t, err)

	err = store.CreateSequence(ctx, testSequence)
	require.NoError(t, err)

	fetchedSequence, err := store.FetchSequence(ctx, testSequence.ID)
	require.NoError(t, err)
	require.Equal(t, testSequence, fetchedSequence)
}

func TestUpdateSequence(t *testing.T) {
	ctx := t.Context()
	cleanDB(ctx)

	store, err := pg.NewStore(testPool)
	require.NoError(t, err)

	err = store.CreateSequence(ctx, testSequence)
	require.NoError(t, err)

	err = store.UpdateSequence(ctx, testSequence.ID, &model.Sequence{
		OpenTrackingEnabled:  false,
		ClickTrackingEnabled: false,
	})
	require.NoError(t, err)

	fetchedSequence, err := store.FetchSequence(ctx, testSequence.ID)
	require.NoError(t, err)
	require.False(t, fetchedSequence.OpenTrackingEnabled)
	require.False(t, fetchedSequence.ClickTrackingEnabled)
}

func TestUpdateStep(t *testing.T) {
	ctx := t.Context()
	cleanDB(ctx)

	store, err := pg.NewStore(testPool)
	require.NoError(t, err)

	err = store.CreateSequence(ctx, testSequence)
	require.NoError(t, err)

	testStep := testSequence.Steps[0]
	err = store.UpdateStep(ctx, testStep.ID, &model.Step{
		SequenceID: testSequence.ID,
		Subject:    "Updated Step 1 Subject",
		Content:    "Updated Step 1 Content",
	})
	require.NoError(t, err)

	fetchedSequence, err := store.FetchSequence(ctx, testSequence.ID)
	require.NoError(t, err)
	require.Equal(t, "Updated Step 1 Subject", fetchedSequence.Steps[0].Subject)
	require.Equal(t, "Updated Step 1 Content", fetchedSequence.Steps[0].Content)
}

func TestDeleteStep(t *testing.T) {
	ctx := t.Context()
	cleanDB(ctx)

	store, err := pg.NewStore(testPool)
	require.NoError(t, err)

	err = store.CreateSequence(ctx, testSequence)
	require.NoError(t, err)

	testStep := testSequence.Steps[0]
	err = store.DeleteStep(ctx, testStep.ID, &model.Step{
		SequenceID: testSequence.ID,
	})
	require.NoError(t, err)

	fetchedSequence, err := store.FetchSequence(ctx, testSequence.ID)
	require.NoError(t, err)
	require.Len(t, fetchedSequence.Steps, 1, "Expected 1 step after deletion")
}
