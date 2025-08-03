package repositories

import (
	"context"
	"errors"
	"subscription-service/internal/domain/models"
	"subscription-service/internal/infrastructure/database"
	infraerror "subscription-service/internal/infrastructure/errors"
	"subscription-service/internal/infrastructure/mappers"
	"time"
	"weather-forecast/pkg/logger"

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
	log := r.logger.WithContext(ctx)

	log.Debugf("Creating subscription for email: %s", subscription.Email)

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		dbSubscription := mappers.DomainToDatabase(subscription)

		res := r.db.WithContext(ctx).Create(&dbSubscription)

		if res.Error != nil {
			log.Errorf("Failed to save subscription to database: %s", res.Error.Error())
			return nil, infraerror.ErrDatabase
		}
		domainSubscription := mappers.DatabaseToDomain(dbSubscription)

		log.Debugf("Subscription created successfully for email: %s", subscription.Email)

		return &domainSubscription, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*models.Subscription), nil
}

func (r *SubscriptionRepository) GetByEmail(ctx context.Context, email string) (*models.Subscription, error) {
	log := r.logger.WithContext(ctx)

	log.Debugf("Looking up subscription by email: %s", email)

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		dbSubscription := database.Subscription{}
		res := r.db.WithContext(ctx).Where("email = ?", email).First(&dbSubscription)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				log.Debugf("No subscription found for email: %s", email)
				return nil, nil

			} else {
				log.Errorf("Failed to get subscription from database: %s", res.Error.Error())
				return nil, infraerror.ErrDatabase
			}
		}
		domainSubscription := mappers.DatabaseToDomain(dbSubscription)

		log.Debugf("Subscription found for email: %s", email)
		return &domainSubscription, nil
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
	log := r.logger.WithContext(ctx)

	log.Debugf("Looking up subscription by token: %s", token)

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		dbSubscription := database.Subscription{}
		res := r.db.WithContext(ctx).Where("token = ?", token).First(&dbSubscription)

		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				log.Debugf("No subscription found for token: %s", token)
				return nil, nil

			} else {
				log.Errorf("Failed to get subscription from database: %s", res.Error.Error())
				return nil, infraerror.ErrDatabase
			}
		}

		domainSubscription := mappers.DatabaseToDomain(dbSubscription)

		log.Debugf("Subscription found for token: %s", token)
		return &domainSubscription, nil
	})

	if err != nil {
		return nil, err
	}

	if res != nil {
		return res.(*models.Subscription), nil
	}

	return nil, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, subscription models.Subscription) (*models.Subscription, error) {
	log := r.logger.WithContext(ctx)

	log.Debugf("Updating subscription with id: %s", subscription.ID)

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		dbSubscription := mappers.DomainToDatabase(subscription)

		res := r.db.WithContext(ctx).Save(&dbSubscription)

		if res.Error != nil {
			log.Errorf("Failed to update subscription: %s", res.Error.Error())
			return nil, infraerror.ErrDatabase
		}

		domainSubscription := mappers.DatabaseToDomain(dbSubscription)

		log.Debugf("Subscription with id %s successfuly updated", subscription.ID)

		return &domainSubscription, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(*models.Subscription), nil

}

func (r *SubscriptionRepository) DeleteByToken(ctx context.Context, token string) error {
	log := r.logger.WithContext(ctx)

	log.Debugf("Deleting subscription by token: %s", token)

	_, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {
		res := r.db.WithContext(ctx).Where("token = ?", token).Delete(&database.Subscription{})

		if res.Error != nil {
			log.Errorf("Failed to delete subscription: %s", res.Error.Error())
			return nil, infraerror.ErrDatabase
		}

		log.Debugf("Subscription deleted successfully for token: %s", token)
		return nil, nil
	})
	return err
}

func (r *SubscriptionRepository) ListConfirmedByFrequency(ctx context.Context, frequency models.Frequency, lastID, pageSize int) ([]models.Subscription, error) {
	log := r.logger.WithContext(ctx)

	log.Debugf("Listing subscription by frequency: %s", frequency)

	res, err := r.runWithDeadline(ctx, func(ctx context.Context) (any, error) {

		var dbSubscriptions []database.Subscription
		res := r.db.WithContext(ctx).Where("confirmed = ? AND frequency = ? AND id > ?", true, frequency, lastID).Order("id").Limit(pageSize).Find(&dbSubscriptions)

		if res.Error != nil {
			log.Errorf("Failed to list subscriptions: %s", res.Error.Error())
			return nil, infraerror.ErrDatabase
		}
		domainSubscriptions := mappers.DatabaseSliceToDomain(dbSubscriptions)

		log.Debugf("Subscriptions successfuly listed for frequency: %s", frequency)
		return domainSubscriptions, nil
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
