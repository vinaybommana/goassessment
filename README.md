# Go Assessment solution
## run the server
```bash
cd goassessment;
go run cmd/serve.go -port 8080

# output
running web server on port: 8080
```

## validation
### Milestone 1
```bash
curl localhost:8080

# output
Hello world!
```

### Milestone 2
```bash
curl localhost:8080/time

# output
2024-03-11@770.5
```

### Milestone 3
```bash
curl \
--include \
--no-buffer \
--header "Connection: Upgrade" \
--header "Upgrade: websocket" \
--header "Host: localhost:8080" \
--header "Sec-WebSocket-Key: 08kp54j1E3z4IfuM1m75tQ==" \
--header "Sec-WebSocket-Version: 13" \
http://localhost:8080/ws

# output
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: Dhg2s/gXkbuTiYk690dns7LmEhY=

2024-03-11@763.5
```

### Milestone 4
restart server and open `localhost:8080/timeupdating` in local web browser


