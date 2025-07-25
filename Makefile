.PHONY: stop-shared start-shared up-full down-full dev-weather dev-gateway dev-subscription dev-email dev-weather-broadcast stop-weather stop-gateway stop-subscription stop-email stop-weather-broadcast

start-shared:
	@echo "Starting shared services..."
	docker compose up -d
	@echo "Shared services started"



stop-shared:
	docker compose down

up-full: start-shared dev-weather dev-subscription dev-gateway dev-email dev-weather-broadcast
	@echo "Full system is up"

down-full: stop-weather stop-gateway stop-subscription stop-email stop-weather-broadcast stop-shared
	@echo "Everything stopped"

dev-weather:
	@echo "Starting Weather service..."
	cd services/weather && docker compose up -d
	@echo "Weather service started"

dev-gateway:
	@echo "Starting Gateway service..."
	cd services/gateway && docker compose up -d
	@echo "Gateway service started"

dev-subscription:
	@echo "Starting Subscription service..."
	cd services/subscription && docker compose up -d
	@echo "Subscription service started"

dev-email:
	@echo "Starting Email service..."
	cd services/email && docker compose up -d
	@echo "Email service started"

dev-weather-broadcast:
	@echo "Starting Weather Broadcast service..."
	cd services/weather-broadcast && docker compose up -d
	@echo "Weather Broadcast service started"

stop-weather:
	@echo "Stopping Weather service..."
	cd services/weather && docker compose down

stop-gateway:
	@echo "Stopping Gateway service..."
	cd services/gateway && docker compose down

stop-subscription:
	@echo "Stopping Subscription service..."
	cd services/subscription && docker compose down

stop-email:
	@echo "Stopping Email service..."
	cd services/email && docker compose down

stop-weather-broadcast:
	@echo "Stopping Weather Broadcast service..."
	cd services/weather-broadcast && docker compose down

logs-shared:
	docker compose logs -f

logs-weather:
	cd services/weather && docker compose logs -f

logs-gateway:
	cd services/gateway && docker compose logs -f

logs-subscription:
	cd services/subscription && docker compose logs -f

logs-email:
	cd services/email && docker compose logs -f

logs-weather-broadcast:
	cd services/weather-broadcast && docker compose logs -f