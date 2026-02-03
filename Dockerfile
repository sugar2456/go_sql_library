FROM golang:1.25-alpine

WORKDIR /app

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Install dependencies
COPY go.mod go.sum* ./
RUN go mod download 2>/dev/null || go mod tidy

# Copy source code
COPY . .

# Ensure dependencies are resolved
RUN go mod tidy

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
