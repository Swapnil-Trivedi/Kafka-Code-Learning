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