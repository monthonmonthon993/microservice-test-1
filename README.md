## RUN Server
```sh
$ cd server
$ go run cmd/main.go
```

## RUN Kafka
```sh
$ cd kafka
$ docker-compose up -d
```

## RUN MongoDB
```sh
$ cd mongodb
$ docker-compose up -d
```

## RUN Workers
```sh
$ cd workers/worker1
$ go run main.go
```
repeat step for worker 2 and 3

## Test Client
```sh
$ cd client
$ go run cmd/main.go
```