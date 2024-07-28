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


In the Make file
	# @go build -o bin/receiver data_receiver/main.go 
	## If you are going to build this with main.go 
	# it's going to build this main.go file singlely. 
	# So basically what is oging to happen is if you
	# have other files that main depends on, it's not
	# going to work. So you need to build the whole folder.   

.PHONY in the Makefile - using it when the Makefile is not updating my main file. If it tells you there is no updates, put it to .Phony to force it to update.


LEFT OFF: Coding the invoice aggregator 20:36
