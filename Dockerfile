# Мультистадийная сборка

# Этап сборки фронтенда
FROM node:18-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Этап сборки бэкенда
FROM golang:1.18-alpine AS backend-builder
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /metrics-server

# Финальный этап
FROM alpine:3.18
WORKDIR /app
COPY --from=frontend-builder /app/frontend/dist /app/static
COPY --from=backend-builder /metrics-server /app/
EXPOSE 8080
CMD ["/app/metrics-server"] 