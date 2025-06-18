# Salesforge Golang Engineer Task

## Development

### Prepare environment variables

Copy env example and update variables:

```sh
cp .env.example .env
```

### Prepare database

Create test and development databases:

```sh
make create-db
```

Run migrations:

```sh
make migrate
```

### Run application

To run application:

```sh
make dev
```

## Testing

```sh
make test
```

or verbose:

```sh
go test -v ./...
```

## E2E testing

### Create sequence

#### Request

```sh
curl --request POST \
  --url http://localhost:8080/sequences \
  --header 'content-type: application/json' \
  --data '{
  "name": "Test Sequence",
  "openTrackingEnabled": true,
  "clickTrackingEnabled": true,
  "steps": [
    {
      "subject": "Test Subject",
      "content": "Test Content"
    }
  ]
}'
```

#### Response

```json
{
  "id": 1,
  "name": "Test Sequence",
  "openTrackingEnabled": true,
  "clickTrackingEnabled": true,
  "createdAt": "2025-06-18T23:37:15.639442Z",
  "updatedAt": "2025-06-18T23:37:15.639442Z",
  "steps": [
    {
      "id": 1,
      "subject": "Test Subject",
      "content": "Test Content",
      "createdAt": "2025-06-18T23:37:15.639442Z",
      "updatedAt": "2025-06-18T23:37:15.639442Z"
    }
  ]
}
```

### Get sequence

#### Request

```sh
curl --request GET \
  --url http://localhost:8080/sequences/1
```

#### Response

```json
{
  "id": 1,
  "name": "Test Sequence",
  "openTrackingEnabled": false,
  "clickTrackingEnabled": false,
  "createdAt": "2025-06-18T23:37:15.639442Z",
  "updatedAt": "2025-06-18T23:39:55.435925Z",
  "steps": [
    {
      "id": 1,
      "subject": "Test Subject",
      "content": "Test Content",
      "createdAt": "2025-06-18T23:37:15.639442Z",
      "updatedAt": "2025-06-18T23:37:15.639442Z"
    }
  ]
}
```

### Update sequence

#### Request

```sh
curl --request PUT \
  --url http://localhost:8080/sequences/1 \
  --header 'content-type: application/json' \
  --data '{
  "openTrackingEnabled": false,
  "clickTrackingEnabled": false
}'
```

#### Response

```json
{
  "id": 1,
  "openTrackingEnabled": false,
  "clickTrackingEnabled": false,
  "updatedAt": "2025-06-18T23:39:55.435925Z"
}
```

### Update step

#### Request

```sh
curl --request PUT \
  --url http://localhost:8080/sequences/1/steps/1 \
  --header 'content-type: application/json' \
  --data '{
  "subject": "Updated Subject",
  "content": "Updated Content"
}'
```

#### Response

```json
{
  "id": 1,
  "subject": "Updated Subject",
  "content": "Updated Content",
  "createdAt": "2025-06-18T23:37:15.639442Z",
  "updatedAt": "2025-06-18T23:41:41.086308Z"
}
```

### Delete step

#### Request

```sh
curl --request DELETE \
  --url http://localhost:8080/sequences/1/steps/1
```

#### Response

```Status Code - 204```
