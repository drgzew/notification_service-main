# Сервис рассылки уведомлений

## Требования
- Go (совместимая версия указана в `go.mod`)
- Docker + Docker Compose (для локальных контейнеров Postgres и Kafka)


docker compose up -d zookeeper kafka
docker compose up -d postgres 
docker compose up app