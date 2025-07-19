package dto

type (
	Weather struct {
		Temperature float64
		Humidity    int
		Description string
	}

	WeatherMailSuccessInfo struct {
		Email   string
		City    string
		Weather Weather
	}
	WeatherMailErrorInfo struct {
		Email string
		City  string
	}
)
