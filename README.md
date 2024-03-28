# PUBG

This is a small application meant to be run as a cron job to update a redis database with some stats on players in the pubg leaderboard

## Build
To build the container directly into minikube

```shell
eval $(minikube docker-env)
docker build -t pubg . 
```

## Tests
```shell
go test
```
