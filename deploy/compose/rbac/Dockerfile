# ============================
# Stage 1: Build stage
# ============================

FROM golang:1.24.6-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum ./

COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY iam/go.mod iam/go.sum ./iam/
COPY rbac/go.mod rbac/go.sum ./rbac/

RUN go mod download

COPY platform ./platform
COPY shared ./shared
COPY iam ./iam
COPY rbac ./rbac

ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

RUN chmod +x grpc-health-probe

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-rbac ./rbac/cmd/main.go

# ============================
# Stage 2: Final image
# ============================

FROM alpine:3.21.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app
USER appuser

COPY --from=builder /app/app-rbac .
COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe
COPY rbac/config ./config
COPY rbac/db ./db

ENV CONFIG_PATH=/app/config/development.yaml

EXPOSE 50052

ENTRYPOINT ["./app-rbac"]
