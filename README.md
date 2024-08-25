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


What is a GATEWAY? In a Micro Service infastructure you are going to have multiple services,
and there is going to be some kind of client. The client could be some kind of mobile
application or a frontend. It could be anything. And that frontend is going to 
retrieve some data. For example, in our case the frontend wants to fetch an invoice for 
a specific OBU. Give me all the coordinates send during this month for this OBU. SO directly
form the IU we are going to access our gateway. Why? Because all our authentication, all our stuff
is going to be in that gateway. So we don't need to do that for these Micro Services. We don't 
want to hassle with authentication in our micro services because for some companies this effect could have
100, 200 micro services. We don't want to implement all that logic. So the GATEWAY is going to
be responsible for that stuff. It's basically, our API. A GATEWAY is a user facing API.
That is only used by the end user. For example, let's say we build a tolling stack, 
it's going to have a mobile app for some people, a UI for some people, a UI for administration,
a UI for the customers, so they are all going to interact with the GATEWAY. 
They are not going to interact with the services but with the GATEWAY. It's a very simple 
JSON API that sits in front of all our microservices. It's going to handle authentication,
and everything else. Most of the time it is going to recieve JSON.  The GATEWAY is a proxy between
the users, the clients, the mobile apps and the services. The GATEWAY needs to sit in between.

INSTRUMENT - means have the aggregator have perfect logging, make the aggregator have perfect metrics.
As if we are going to production.
Instrumentaions is basically makig sure that we have logging already but we also want ot have some metrics.
For example of metrics we want, maybe we want to have a counter that counts each time we access some kind of ....
You should have metrics for each single service. 

Prometheus -- is a deamon (runs in the background) that runs on your system. It will scrape 
metrics from the Prometheus client (endpoint we need to configure). 
And then you can expose them to Grafana.

Grafana - one of the most used things to do metrics.

grpc is actually tcp. We need to have a port etc..

LEFT OFF: Finshing the aggregator service with go-kit
