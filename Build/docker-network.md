### To check existing docker networks
> docker network ls

### Creating Subnet for setup
> docker network create --subnet 172.16.1.0/24 st-local-net

### Remove network
> docker network rm <Id>

### Run Zookeeper
> docker-compose -f zk-docker-compose.yaml up -d zookeeper

### Run single instance of Kafka container
> docker-compose -f zk-docker-compose.yaml -f kafka-docker-compose.yaml up -d kafka-1

### Exec bash in container
> docker exec -it kafka-1 bash