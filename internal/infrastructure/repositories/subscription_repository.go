package repositories

import (
	"context"
	"errors"
	"time"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"

	"gorm.io/gorm"
)

const DB_TIMEOUT = 3 * time.Second

type SubscriptionRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewSubscriptionRepository(db *gorm.DB, logger logger.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		res := r.db.WithContext(ctx).Create(&subscription)

		if res.Error != nil {
			r.logger.Warnf("Failed to save subscription to database: %s", res.Error.Error())
			return nil, apperrors.DatabaseError
		}

		return &subscription, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*models.Subscription), nil
}

func (r *SubscriptionRepository) GetByEmail(ctx context.Context, email string) (*models.Subscription, error) {

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		subscription := models.Subscription{}
		res := r.db.WithContext(ctx).Where("email = ?", email).First(&subscription)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return nil, nil

			} else {
				r.logger.Warnf("Failed to get subscription from database: %s", res.Error.Error())
				return nil, apperrors.DatabaseError
			}
		}

		return &subscription, nil
	})

	if err != nil {
		return nil, err
	}

	if res != nil {
		return res.(*models.Subscription), nil
	}

	return nil, nil
}

func (r *SubscriptionRepository) GetByToken(ctx context.Context, token string) (*models.Subscription, error) {

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		subscription := models.Subscription{}
		res := r.db.WithContext(ctx).Where("token = ?", token).First(&subscription)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return nil, nil

			} else {
				r.logger.Warnf("Failed to get subscription from database: %s", res.Error.Error())
				return nil, apperrors.DatabaseError
			}
		}

		return &subscription, nil
	})

	if err != nil {
		return nil, err
	}

	if res != nil {
		return res.(*models.Subscription), nil
	}

	return nil, nil

	return res.(*models.Subscription), nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		res := r.db.WithContext(ctx).Save(&subscription)

		if res.Error != nil {
			r.logger.Warnf("Failed to update subscription: %s", res.Error.Error())
			return nil, apperrors.DatabaseError
		}
		return &subscription, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*models.Subscription), nil

}

func (r *SubscriptionRepository) DeleteByToken(ctx context.Context, token string) error {

	_, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {
		res := r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.Subscription{})

		if res.Error != nil {
			r.logger.Warnf("Failed to delete subscription: %s", res.Error.Error())
			return nil, apperrors.DatabaseError
		}

		return nil, nil
	})
	return err
}

func (r *SubscriptionRepository) ListConfirmedByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error) {

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		var subscriptions []models.Subscription
		res := r.db.WithContext(ctx).Where("confirmed = ? AND frequency = ? AND id > ?", true, frequency, lastID).Order("id").Limit(pageSize).Find(&subscriptions)

		if res.Error != nil {
			r.logger.Warnf("Failed to list subscriptions: %s", res.Error.Error())
			return nil, apperrors.DatabaseError
		}

		return subscriptions, nil
	})

	if err != nil {
		return nil, err
	}

	return res.([]models.Subscription), nil

}

func (r *SubscriptionRepository) runWithDeadline(ctx context.Context, handler func(ctx context.Context) (any, error)) (any, error) {
	ctx, cancel := context.WithTimeout(ctx, DB_TIMEOUT)
	defer cancel()
	return handler(ctx)
}
