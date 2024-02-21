# Pismo

### Config:
This application takes configuration through enviroment variables which are:
[server.env](cmd%2Fserver.env)

### Prerequisites:
- Docker

### Build and Run 
This will build and spin up both the app and mongodb docker

```shell
chmod +x run.sh
./run.sh
```

To edit environment variables edit them in [docker-compose.yaml](docker-compose.yaml)

### APIS:

#### POST localhost:2020/v1/accounts

Request:

```json
{
  "document_number": "1023929321"
}
```

Response:

```json
{
  "account_id": "65d6787f20e3ade6411f8c09"
}
```

#### Get localhost:2020/v1/accounts/:accID

Response:

```json
{
  "account_id": "65d6786b3343145256f7dfc5",
  "document_number": "1023929321"
}
```

#### POST localhost:2020/v1/transactions

Request:

```json
{
  "account_id": "65d6786b3343145256f7dfc5",
  "operation_type_id": 4,
  "amount": 20
}
```

Response:

```json
{
  "transaction_id": "65d6795304b10134d9a1be7f"
}
```

#### GET http://localhost:2020/v1/health/
Response:
```json
{
    "message": "service is up and running"
}
```


#### Note In Case of error responses will have status code and error message like this:

```json
{
    "message": "account don't exists"
}
```