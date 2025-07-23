package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestCleanArchitecture(t *testing.T) {

	t.Run("domain models independence", func(t *testing.T) {
		archtest.Package(t, "weather-forecast/internal/domain/models").
			ShouldNotDependOn("weather-forecast/internal/domain/usecases/...")
	})

	t.Run("usecase layer isolation", func(t *testing.T) {
		archtest.Package(t, "weather-forecast/internal/domain/usecases").
			Ignoring("weather-forecast/internal/infrastructure/logger/...").
			ShouldNotDependOn(
				"weather-forecast/internal/infrastructure/...",
				"weather-forecast/internal/presentation/...",
			)
	})

	t.Run("infrastructure independence from presentation", func(t *testing.T) {
		archtest.Package(t, "weather-forecast/internal/infrastructure/...").
			ShouldNotDependOn("weather-forecast/internal/presentation/...")
	})

}
