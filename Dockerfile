FROM golang:1.23 as builder
WORKDIR /app
COPY . .
RUN go install github.com/google/wire/cmd/wire@latest
RUN cd cmd/ordersystem && wire
RUN cd cmd/ordersystem && go build -o app

FROM golang:1.23
WORKDIR /app
COPY --from=builder /app/cmd/ordersystem/app .
CMD ["./app"]
