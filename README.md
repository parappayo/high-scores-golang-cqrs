# high-scores-golang-cqrs

Demo of a Go web service broken up using CQRS and workers-queues patterns

The primary goal is learning, with a secondary goal of demonstrating what I can do.

The usual caveat applies: do NOT run this code in prod. It is for demo purposes only. No warranty is provided.

## Setup

### Prereqs

* a RabbitMQ server that you can use
* a PostgreSQL server that you can use

### PostgreSQL Schema Setup

Coming Soon, basically run some SQL scripts

### Build

* `make install`
* `make build`

### Start Services

* start the services
  * TODO: rename these services to say "service"
  * `./cmd/get-scores`
  * `./cmd/post-score`

### Synthesize Traffic

Coming Soon

### Demo Frontend

Coming Soon, basically open a static html page hosted by nginx

## Tech Stack

* [nginx](https://www.nginx.com/) - http proxy
* [Go](https://go.dev/) - http services, worker processes
* [PostgreSQL](https://www.postgresql.org/) - database
* [Redis](https://redis.com/) - document cache
* [RabbitMQ](https://www.rabbitmq.com/) - command queue

## Backend Architecture

Typically I would create a separate repo for each service, but since this is all one demo, instead `go build` will produce a binary for each service (or tool) in the `bin/` dir.

Web requests go to nginx and are routed based on whether they want to change state (issue a command), or read a document.

State change requests go to the `post-score` service where they are validated and emit commands to the worker queue. A response of 202 Accepted with a job ID is typical.

Document read requests go to the `get-scores` service where results are served from the cache.

Job status requests go to the `get-jobs` service where results are served from the cachce. Jobs that have completed succesfully return a document ID where applicable.

The `worker` process spins up periodically to accept batches of work (commands) from the queue. Each commands typically either results in a write to the database, or checking to see if a document needs to be generated and written to the cache.

## Frontend

A `datagen` command is provided for populating the db with dummy data for demo purposes. To make the system more realistic it works by issuing web requests at random intervals.

A static web page is provided with filters for viewing leaderboards. First the given doc is requested, and if not found then a request is issued to populate the cache.

## Local Server Setup

### Debian / Ubuntu

#### RabbitMQ

Install
* `sudo apt-get install rabbitmq-server`
* `systemctl status rabbitmq-server`

Debug
* `rabbitmq-plugins enable rabbitmq_management`
* browse to `http://localhost:15672` default login is `guest` pw `guest`

## Curl Testing

* `./bin/get-score`
* `curl http://localhost:3000`

* `./bin/post-score`
* `curl -d '{}' -H "Content-Type: application/json" -X POST http://localhost:3010/`

## References

* [CQRS (Martin Fowler)](https://martinfowler.com/bliki/CQRS.html)
* [Web-Queue-Worker (Microsoft)](https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/web-queue-worker)
* [Tutorial: Developing a RESTful API with Go and Gin](https://go.dev/doc/tutorial/web-service-gin)
