FROM golang:1.22.6 AS builder
WORKDIR /src
COPY . .
# CGO_ENABLED - позволяет использовать итегрировать c в golang (0 означает не включен)
# GOOS - ОС под которую будет скомпилирован исполняемый файл
RUN CGO_ENABLED=0 GOOS=linux go build -o url-parser

FROM scratch
WORKDIR /app
COPY --from=builder /src/urls.txt .
COPY --from=builder /src/url-parser .
ENTRYPOINT ["./url-parser"]