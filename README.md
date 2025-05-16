# 🌤️ Weather Forecast 

[🔗 Live Demo](http://weather-update.pp.ua)  - link to the deployed service

**Weather Forecast API** is a web service that allows users to subscribe to weather forecast updates for a selected city. The service supports two types of subscriptions: **hourly** and **daily**. After subscribing, users receive a confirmation email to activate their subscription.

---

## 🚀 Quick Start

To run the service locally:

```bash
git clone https://github.com/ANTONICE54/weather-forecast-api.git
cd weather-forecast-api
docker compose up
```

> ⚠️ For local email testing, [Mailtrap.io](https://mailtrap.io) is used. In the deployed version, the service is configured with `smtp.gmail.com`.

---

## 🖼️ Web Interface

The project includes a simple HTML page with a form for subscriptions:
- Email address
- City name
- Frequency: hourly or daily

![image](https://github.com/user-attachments/assets/55772f07-8c17-453e-aedb-94567e938449)


---

## 📬 API Endpoints


### GET /weather

Get the current weather forecast

##### Example Input: 
```
{
	"city": "Kyiv"
} 
```

### POST /subscribe

Subscribe to weather updates (a confirmation email will be sent)

##### Example Input: 
```
{
	"email": "youremail@mail.com",
	"city": "Kyiv",
	"frequency": "hourly"
} 
```



### GET /confirm/{token}

Confirm the email subscription

##### URL Parameters:
- `token` – confirmation token (UUID format) sent via email

##### Example:
`GET /confirm/3fa85f64-5717-4562-b3fc-2c963f66afa6`

### GET /unsubscribe/{token}

Cancel the email subscription.

##### URL Parameters:
- `token` – unsubscribe token (UUID format)

##### Example:
`GET /unsubscribe/3fa85f64-5717-4562-b3fc-2c963f66afa6`

---

## 🛠️ Technologies Used

- **Go** with [Gin](https://github.com/gin-gonic/gin)
- **Docker** & Docker Compose
- **SMTP**: Gmail (production) / Mailtrap (testing)
- **HTML** frontend for subscriptions

---


## ⚙️ Environment Variables

The service uses the following environment variables defined in a `.env` file:

| Variable             | Description |
|----------------------|-------------|
| `TIMEZONE`           | Timezone for scheduling email delivery (e.g., `Europe/Kyiv`). |
| `SERVER_HOST`        | Public host URL of the API. |
| `SERVER_PORT`        | Port on which the server runs locally. |
| `DB_USER`            | Username for the PostgreSQL database. |
| `DB_PASSWORD`        | Password for the PostgreSQL database. |
| `DB_NAME`            | Name of the PostgreSQL database. |
| `DB_HOST`            | Hostname of the PostgreSQL server (e.g., `postgres`). |
| `DB_PORT`            | Port for the PostgreSQL server (default: `5432`). |
| `WEATHER_API_URL`    | URL of the weather API endpoint used to fetch current weather data. |
| `WEATHER_API_KEY`    | API key to access the weather service. |
| `MAILER_HOST`        | SMTP host used for sending emails (e.g., Gmail or Mailtrap). |
| `MAILER_PORT`        | Port used for the SMTP server (e.g., `587` for Gmail). |
| `MAILER_USERNAME`    | Username/email used for SMTP authentication. |
| `MAILER_PASSWORD`    | Password or app-specific password for the SMTP server. |
| `MAILER_FROM`        | Email address that will appear in the "From" field of sent emails. |

