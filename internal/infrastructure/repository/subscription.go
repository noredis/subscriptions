package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/noredis/subscriptions/internal/domain/entity"
	"github.com/noredis/subscriptions/internal/domain/failure"
	"github.com/noredis/subscriptions/internal/domain/interfaces"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) interfaces.SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (repo *SubscriptionRepository) Insert(
	ctx context.Context,
	sub *entity.Subscription,
) (*entity.Subscription, error) {
	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Insert("subscriptions").
		Columns("service_name", "price", "user_id", "start_date", "end_date").
		Values(sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, err
	}

	var id int
	if err := repo.db.QueryRow(ctx, query, args...).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, failure.ErrUserAlreadyHasThisSubscription
		}

		return nil, err
	}

	sub.ID = id
	return sub, nil
}

func (repo *SubscriptionRepository) Update(
	ctx context.Context,
	sub *entity.Subscription,
) (*entity.Subscription, error) {
	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Update("subscriptions").
		Set("service_name", sub.ServiceName).
		Set("price", sub.Price).
		Set("user_id", sub.UserID).
		Set("start_date", sub.StartDate).
		Set("end_date", sub.EndDate).
		Where(squirrel.Eq{"id": sub.ID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	if _, err := repo.db.Exec(ctx, query, args...); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, failure.ErrUserAlreadyHasThisSubscription
		}

		return nil, err
	}

	return sub, nil
}

func (repo *SubscriptionRepository) Delete(ctx context.Context, id int) error {
	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Delete("subscriptions").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	if _, err := repo.db.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil
}

func (repo *SubscriptionRepository) ExistsByID(
	ctx context.Context,
	id int,
) (bool, error) {
	query, args, err := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("1").
		From("subscriptions").
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return false, err
	}

	var dummy int
	if err := repo.db.QueryRow(ctx, query, args...).Scan(&dummy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (repo *SubscriptionRepository) FindByID(
	ctx context.Context,
	id int,
) (*entity.Subscription, error) {
	query, args, err := repo.getQuery().
		Where(squirrel.Eq{"id": id}).
		Limit(1).
		ToSql()
	if err != nil {
		return nil, err
	}

	var sub entity.Subscription
	err = repo.db.QueryRow(ctx, query, args...).
		Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, failure.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &sub, nil
}

func (repo *SubscriptionRepository) Find(
	ctx context.Context,
	f *entity.SubscriptionFilter,
) ([]*entity.Subscription, error) {
	subscriptions := make([]*entity.Subscription, 0)

	qb := repo.getQuery()
	qb = repo.filterHelper(qb, f)

	offset := (f.Page - 1) * f.Limit
	qb = qb.Limit(uint64(f.Limit)).Offset(uint64(offset))

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub entity.Subscription
		err := rows.
			Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (repo *SubscriptionRepository) FindAll(
	ctx context.Context,
	f *entity.SubscriptionFilter,
) ([]*entity.Subscription, error) {
	subscriptions := make([]*entity.Subscription, 0)

	qb := repo.getQuery()
	qb = repo.filterHelper(qb, f)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var sub entity.Subscription
		err := rows.
			Scan(&sub.ID, &sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (repo *SubscriptionRepository) Total(
	ctx context.Context,
	f *entity.SubscriptionFilter,
) (int, error) {
	cb := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("COUNT(*)").
		From("subscriptions")
	cb = repo.filterHelper(cb, f)

	query, args, err := cb.ToSql()
	if err != nil {
		return 0, err
	}

	var total int
	if err := repo.db.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, err
}

func (repo *SubscriptionRepository) filterHelper(
	sb squirrel.SelectBuilder,
	f *entity.SubscriptionFilter,
) squirrel.SelectBuilder {
	if f.ServiceName != "" {
		sb = sb.Where(squirrel.Eq{"service_name": f.ServiceName})
	}

	if f.UserID != "" {
		sb = sb.Where(squirrel.Eq{"user_id": f.UserID})
	}

	if f.StartDate != nil {
		sb = sb.Where(
			squirrel.Or{
				squirrel.GtOrEq{"end_date": f.StartDate},
				squirrel.Expr("end_date IS NULL"),
			},
		)
	}

	if f.EndDate != nil {
		sb = sb.Where(squirrel.LtOrEq{"start_date": f.EndDate})
	}

	return sb
}

func (repo *SubscriptionRepository) getQuery() squirrel.SelectBuilder {
	return squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("id", "service_name", "price", "user_id", "start_date", "end_date").
		From("subscriptions")
}
