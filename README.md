Micro Service - A single small piece of unit (computatinal unit or storage unit) and it's going to do one thing.
It's going to do that thing very well. Very ochestrated, measured, and instrumentated. It's going to have different 
transports and maybe there will be a storage attached to it.

OBU (Onboard Unit) - sits on the inside of a truck (vehicle on the road) and sends GPS coordinates at intervals. And we are going to send them (we are going to replicate that).
We are going to use a web socket connection that will send these messages over the web socket and we are going to recieve that in our first micro service. And
then we are going to pit it on KAFKA (Queue).

make receiver must be run first, then make obu to initialize the connection.

In KAFKA we have the concept of the producer and the consumer.
The producer is going to produce. the receiver receives messages and it is
going to produce that on a topic. And then we are going to have another 
service later on and that is going to consume that data. 
Kafka depends on the zookeeper library otherwise it will not work.

DOCKER COMMAND TO START IMAGE / CONTAINER
$ docker-compose up -d --force-recreate


LEFT OFF: Producing to Kafka with logging middleware: 13:20
