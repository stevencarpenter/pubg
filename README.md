# PUBG

This is a small application meant to be run as a cron job to update a redis database with some stats on players in the pubg leaderboard. In this case we are storing the player's keyed on their accountId with a json payload as the value that contains the player's rank, win total, and games played.

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
