# Golang Web Socket

This repo will help you for a basic websocket connection in Golang
## Cassandra Setup
```
CREATE KEYSPACE notifications WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
```
```
CREATE TABLE notifications.users (id text PRIMARY KEY, name text);
```
```
INSERT INTO notifications.users (id, name) VALUES ('user1', 'Alice');

INSERT INTO notifications.users (id, name) VALUES ('user2', 'Bob');
````
## Run Go file

```
go run main.go
```
Will run go on port 8081.

websocket will be on the link 

- ws://localhost:8081/ws
- ws://localhost:8081/ws?userID=user1
- ws://localhost:8081/ws?userID=user2

based on DB table entries

## Install http-server

We need a http server to run client code.

```
npm install -g http-server
```
run the http server on same directory
```
http-server -p 8080
```

You can now listen ws connection through

- http://localhost:8080/?userID=user1
- http://localhost:8080/?userID=user2


based on DB table entries

## Send Notification

```
curl --location --request POST 'http://localhost:8081/send?userID=user1' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-raw 'message=Hello%20world!'

```

you can send notification to same connection from front-end through function(via console)

```
sendMessage("hiiiiiiiiiii")
```


