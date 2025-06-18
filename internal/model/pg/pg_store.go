package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/salesforge/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ model.SequenceStore = (*PGStore)(nil)

type PGStore struct {
	pool    *pgxpool.Pool
	builder sq.StatementBuilderType
}

func NewStore(pool *pgxpool.Pool) (*PGStore, error) {
	// Check if the pool is valid and can connect to the database
	// This is a simple ping to ensure the connection is alive.
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PGStore{
		pool:    pool,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}, nil
}

func (s *PGStore) CreateSequence(ctx context.Context, sequence *model.Sequence) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sql, args, err := s.builder.
		Insert("sequences").
		Columns("name", "open_tracking_enabled", "click_tracking_enabled").
		Values(sequence.Name, sequence.OpenTrackingEnabled, sequence.ClickTrackingEnabled).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	if err := tx.QueryRow(ctx, sql, args...).Scan(
		&sequence.ID,
		&sequence.CreatedAt,
		&sequence.UpdatedAt,
	); err != nil {
		return err
	}

	if len(sequence.Steps) > 0 {
		builder := s.builder.
			Insert("steps").
			Columns("sequence_id", "subject", "content").
			Suffix("RETURNING id, subject, content, created_at, updated_at")
		for _, step := range sequence.Steps {
			builder = builder.Values(sequence.ID, step.Subject, step.Content)
		}
		sql, args, err = builder.ToSql()
		if err != nil {
			return err
		}

		rows, err := tx.Query(ctx, sql, args...)
		if err != nil {
			return err
		}
		defer rows.Close()

		var inserted []*model.Step
		for rows.Next() {
			step := &model.Step{SequenceID: sequence.ID}
			if err := rows.Scan(
				&step.ID,
				&step.Subject,
				&step.Content,
				&step.CreatedAt,
				&step.UpdatedAt,
			); err != nil {
				return err
			}
			inserted = append(inserted, step)
		}
		sequence.Steps = inserted
	}

	return tx.Commit(ctx)
}

func (s *PGStore) FetchSequence(ctx context.Context, id uint64) (*model.Sequence, error) {
	sql, args, err := s.builder.
		Select(
			"id",
			"name",
			"open_tracking_enabled",
			"click_tracking_enabled",
			"created_at",
			"updated_at",
		).
		Where(sq.Eq{"id": id}).
		From("sequences").
		ToSql()
	if err != nil {
		return nil, err
	}

	var sequence model.Sequence
	if err := s.pool.QueryRow(ctx, sql, args...).Scan(
		&sequence.ID,
		&sequence.Name,
		&sequence.OpenTrackingEnabled,
		&sequence.ClickTrackingEnabled,
		&sequence.CreatedAt,
		&sequence.UpdatedAt,
	); err != nil {
		return nil, err
	}

	sql, args, err = s.builder.
		Select(
			"id",
			"sequence_id",
			"subject",
			"content",
			"created_at",
			"updated_at",
		).
		Where(sq.Eq{"sequence_id": id}).
		From("steps").
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		step := &model.Step{SequenceID: sequence.ID}
		if err := rows.Scan(
			&step.ID,
			&step.SequenceID,
			&step.Subject,
			&step.Content,
			&step.CreatedAt,
			&step.UpdatedAt,
		); err != nil {
			return nil, err
		}
		sequence.Steps = append(sequence.Steps, step)
	}

	return &sequence, nil
}

func (s *PGStore) UpdateSequence(ctx context.Context, id uint64, sequence *model.Sequence) error {
	sql, args, err := s.builder.
		Update("sequences").
		Set("open_tracking_enabled", sequence.OpenTrackingEnabled).
		Set("click_tracking_enabled", sequence.ClickTrackingEnabled).
		Set("updated_at", "NOW()").
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	if err := s.pool.QueryRow(ctx, sql, args...).Scan(
		&sequence.ID,
		&sequence.CreatedAt,
		&sequence.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *PGStore) UpdateStep(ctx context.Context, id uint64, step *model.Step) error {
	sql, args, err := s.builder.
		Update("steps").
		Set("subject", step.Subject).
		Set("content", step.Content).
		Set("updated_at", "NOW()").
		Where(sq.Eq{"id": id, "sequence_id": step.SequenceID}).
		Suffix("RETURNING id, created_at, updated_at").
		ToSql()
	if err != nil {
		return err
	}

	if err := s.pool.QueryRow(ctx, sql, args...).Scan(
		&step.ID,
		&step.CreatedAt,
		&step.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *PGStore) DeleteStep(ctx context.Context, id uint64, step *model.Step) error {
	sql, args, err := s.builder.
		Delete("steps").
		Where(sq.Eq{"id": id, "sequence_id": step.SequenceID}).
		ToSql()
	if err != nil {
		return err
	}

	cmd, err := s.pool.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
