# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app/server
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server .

# Stage 3: Production Runner
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app

# Copy compiled backend binary and built frontend assets
COPY --from=backend-builder /app/server/server ./server
COPY --from=frontend-builder /app/web/dist ./dist

ENV PORT=8080 \
    BIND_ADDR=0.0.0.0 \
    DB_DRIVER=sqlite \
    SQLITE_PATH=/app/data/internship_tracker.db \
    STATIC_DIR=/app/dist \
    GIN_MODE=release

VOLUME /app/data
EXPOSE 8080

CMD ["./server"]

