package mappers

import (
	"subscription-service/internal/domain/models"
	"subscription-service/internal/infrastructure/database"
)

func DomainToDatabase(domain models.Subscription) database.Subscription {
	return database.Subscription{
		ID:        domain.ID,
		Email:     domain.Email,
		City:      domain.City,
		Token:     domain.Token,
		Frequency: database.Frequency(domain.Frequency),
		Confirmed: domain.Confirmed,
	}
}

func DatabaseToDomain(db database.Subscription) models.Subscription {
	return models.Subscription{
		ID:        db.ID,
		Email:     db.Email,
		City:      db.City,
		Token:     db.Token,
		Frequency: models.Frequency(db.Frequency),
		Confirmed: db.Confirmed,
	}
}

func DatabaseSliceToDomain(dbSubscriptions []database.Subscription) []models.Subscription {
	domainSubscriptions := make([]models.Subscription, len(dbSubscriptions))
	for i, dbSub := range dbSubscriptions {
		domainSubscriptions[i] = DatabaseToDomain(dbSub)
	}
	return domainSubscriptions
}
