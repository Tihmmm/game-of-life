FROM golang:1.22.0-alpine3.19 AS dependencies
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM dependencies AS build
COPY backend ./
RUN CGO_ENABLED=0 go build -o /gol-backend -ldflags="-w -s" ./cmd/app/

FROM tihmmm/golang-alpine-rootless:go-1.22.0-alp-3.19
COPY --chown=user:user --chmod=550 --from=build /gol-backend /home/user/gol-backend
WORKDIR /home/user/
CMD ["./gol-backend"]