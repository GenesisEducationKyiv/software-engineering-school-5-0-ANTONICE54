package database

import (
	"context"
	"errors"
	"weather-forecast/internal/domain/models"
	"weather-forecast/internal/infrastructure/apperrors"
	"weather-forecast/internal/infrastructure/logger"

	"gorm.io/gorm"
)

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

func (r *SubscriptionRepository) Create(subscription models.Subscription) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	res := r.db.WithContext(ctx).Create(&subscription)

	if res.Error != nil {
		r.logger.Warnf("Failed to save subscription to database: %s", res.Error.Error())
		return nil, apperrors.DatabaseError
	}

	return &subscription, nil

}

func (r *SubscriptionRepository) GetByEmail(email string) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

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
}

func (r *SubscriptionRepository) GetByToken(token string) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

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
}

func (r *SubscriptionRepository) Update(subscription models.Subscription) (*models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	res := r.db.WithContext(ctx).Save(&subscription)

	if res.Error != nil {
		r.logger.Warnf("Failed to update subscription: %s", res.Error.Error())
		return nil, apperrors.DatabaseError
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) DeleteByToken(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	res := r.db.WithContext(ctx).Where("token = ?", token).Delete(&models.Subscription{})

	if res.Error != nil {
		r.logger.Warnf("Failed to delete subscription: %s", res.Error.Error())
		return apperrors.DatabaseError
	}

	return nil
}

func (r *SubscriptionRepository) ListConfirmedByFrequency(frequency models.Frequency) ([]models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancel()

	var subscriptions []models.Subscription
	res := r.db.WithContext(ctx).Where("confirmed = ? AND frequency = ?", true, frequency).Find(&subscriptions)

	if res.Error != nil {
		r.logger.Warnf("Failed to list subscriptions: %s", res.Error.Error())
		return nil, apperrors.DatabaseError
	}

	return subscriptions, nil
}
