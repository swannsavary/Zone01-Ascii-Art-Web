FROM golang:1.22.2-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.20

LABEL fr.zone01.titre="docker for ascii-art-web" \
      fr.zone01.description="Docker image for ascii-art-web" \
      fr.zone01.author="swann savary" \
      fr.zone01.version="2.0"

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

# On récupère uniquement ce dont le serveur a besoin à l'exécution
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/*.txt ./
COPY --from=builder /app/ascii-art ./ascii-art

USER app
EXPOSE 8080
CMD ["./main"]