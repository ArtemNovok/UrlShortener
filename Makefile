start:
	@echo Starting url-shortener
	docker-compose down -v
	docker-compose up --build  -d 
	@echo Url-shorener is up and running
stop:
	@echo Stopping containers...
	docker-compose down -v
	@echo containers are stopped
start_local:
	go run cmd/url-shortener/main.go -config=config/local.yaml
