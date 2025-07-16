# Go-API-Transactions

Rest api without auth. Allow to make transactions between two wallets.

## Stack

- cleanenv
- chi (router)
- postgresql (data source)
- pq
- slog
- goose (migrations)
- validator

## Fun facts

- custom colored logger for developing (slogpretty)
- using plpgsql in migrations
- Makefile for developing

## Try it out

On Linux/Unix:
```
CONFIG_PATH=./config/dev.yml go run -C . ./cmd/qual
```

Docker
```
docker build -t qual:latest .
docker run --rm -e CONFIG_PATH=./config/dev.yml -p 8080:8080 qual:latest
```

## Routes

`POST /api/send`
Sends specified amount of money from one to another wallet

**Request body:**
```
{
    "from": "2bc80169",
    "to": "3c229f02",
    "amount": 10
}
```

`GET /api/wallet/{address}/balance`
This should return balance for current wallet adress

`GET /api/transactions?count=N`
This should return last N transactions on server

