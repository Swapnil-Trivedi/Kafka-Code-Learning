## To creat a single topic in kafka
 kafka-topics --create --bootstrap-server kafka-1:9091 \
 --replication-factor 1 --partitions 1 --topic hello-topic

## List all topics
kafka-topics --list --bootstrap-server kafka-1:9091

## Console producer
kafka-console-producer --bootstrap-server kafka-1:9091 \
 --topic hello-topic 

## Console consumer
 kafka-console-consumer --bootstrap-server kafka-1:9091 \
 --topic hello-topic --from-beginning

## Creating topic with replicaation factor for multi-broker setup
kafka-topics --create --bootstrap-server kafka-1:9091 \
--replication-factor 3 --partitions 1 --topic Multibroker-App
> Multi broker setup, partitions are replicated to secondary nodes along with the data and offset.

## Describe topic
 kafka-topics --describe --bootstrap-server kafka-1:9091 \
--topic Multibroker-App

## Alter topic increase partition
 kafka-topics --bootstrap-server kafka-1:9091 --alter \
 --topic Multibroker-App --partitions 2
