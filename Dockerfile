FROM node:lts-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
COPY --from=frontend-builder /app/static/ ./static/

RUN swag init --parseDependency --parseInternal && CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk add --no-cache curl unzip

WORKDIR /app

COPY --from=builder /app/main .

ARG VOICE_ZIP_URL
ARG SE_ZIP_URL
RUN test -n "$VOICE_ZIP_URL" || { echo "VOICE_ZIP_URL build arg is required"; exit 1; } \
    && curl -fSL -o /tmp/voice.zip "$VOICE_ZIP_URL" \
    && mkdir -p internal/quote/data \
    && unzip -qo /tmp/voice.zip -d /tmp/voice \
    && mv /tmp/voice/voice internal/quote/data/audio \
    && rm -rf /tmp/voice.zip /tmp/voice \
    && if [ -n "$SE_ZIP_URL" ]; then \
        curl -fSL -o /tmp/se.zip "$SE_ZIP_URL" \
        && unzip -qo /tmp/se.zip -d /tmp/se \
        && mv /tmp/se/se internal/quote/data/se \
        && rm -rf /tmp/se.zip /tmp/se; \
    fi \
    && apk del curl unzip

EXPOSE 3000

CMD ["./main"]
