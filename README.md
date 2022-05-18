## Create DB

`createdb -h 0.0.0.0 -p 5434 -U postgres links`

## Add db entity

`migrate create -dir ./migrations -ext sql entity`

## Run migrations
`migrate -path ./migrations/ -database postgres://postgres:changeme@localhost:5434/links?sslmode=disable up`

## Rollback migrations

`migrate -path ./migrations/ -database postgres://postgres:changeme@localhost:5434/links?sslmode=disable down`

## API

### Register URL

`curl -i -XPOST -H 'Content-Type: application/json' http://localhost:8080/api/v1/links/ -d '{ "url": "https://utcc.utoronto.ca/~cks/space/blog/programming/GoConcurrencyStillNotEasy" }'`

### Retrieve URL by token

`curl -i -XGET -H 'Authorization: lZ2w3leloRO2YP3uyxsV2jO4RcsgzhAg' http://localhost:8080/api/v1/links/7blHqQ`
