## Proxy Server in Go

* To make a https server, run below commands to generate a self-signed certificate and private key. 
```
chmod +x cert.sh
./cert.sh
``` 

* This contains two proxy server files - one is simple http server and another is https server which maintains a TCP between it and the target server.

```
go run [filename]
``` 

* To check if proxy is working corretly or not , open a local python server and check.
