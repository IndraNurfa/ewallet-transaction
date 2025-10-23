FROM golang:1.23.4

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod tidy

COPY . .

COPY .env .

RUN go build -o ewallet-transaction

RUN chmod +x ewallet-transaction

EXPOSE 8082

CMD [ "./ewallet-transaction" ]