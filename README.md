# DeMedia Poc

## How to run

You have to run a hub and peers on two or more terminals

### Hub

Open a terminal, then

```shell
cd hub
```
```shell
go run main.go
```

### Peer

Open a terminal, then

```shell
cd peer
```
```shell
go run main.go
```

### Endpoints

Get peers mapping
```shell
curl --location --request GET '0.0.0.0:8080/peer'
```

Get all todos of particular peer
```shell
curl --location --request GET '0.0.0.0:8080/todo' \
--header 'Peer: QmcQH5NGNk665u33ZP9zuzYVor6nbF2aopFeCkbsaBSsCs'
```

Create a todo on particular peer
```shell
curl --location --request POST '0.0.0.0:8080/todo' \
--header 'Peer: QmeHyXp4n72gUgCK6je9jcHVzNV9wAxACVPWRrUiwWswVJ' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": "5d48c042-3ae6-4b9d-9c6f-37d58c8922a3",
    "task": "mamama ananan",
    "title": "7fhfhfh"
}'
```

### MinIO Docker

```dockerfile
mkdir -p ~/minio/data &&
docker run \
   -p 9000:9000 \
   -p 9090:9090 \
   --name minio \
   -v ~/minio/data:/data \
   -e "MINIO_ROOT_USER=ROOTNAME" \
   -e "MINIO_ROOT_PASSWORD=CHANGEME123" \
   quay.io/minio/minio server /data --console-address ":9090"
```

Set env like below

```shell
export AWS_ACCESS_KEY_ID=accessKeyID
export AWS_SECRET_ACCESS_KEY=secretAccessKey
```
