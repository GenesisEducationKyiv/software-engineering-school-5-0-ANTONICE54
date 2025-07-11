package dto

type (
	Weather struct {
		Temperature float64
		Humidity    int
		Description string
	}

	WeatherSuccess struct {
		City    string
		Email   string
		Weather Weather
	}

	WeatherError struct {
		City  string
		Email string
	}
)
