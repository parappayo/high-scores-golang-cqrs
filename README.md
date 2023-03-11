# high-scores-golang-cqrs

Demo of a Go web service broken up using CQRS and workers-queues patterns

The primary goal is learning, with a secondary goal of demonstrating what I can do.

The usual caveat applies: do NOT run this code in prod. It is for demo purposes only. No warranty is provided.

## Setup

No code yet, coming soon.

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

## References

* [CQRS (Martin Fowler)](https://martinfowler.com/bliki/CQRS.html)
* [Web-Queue-Worker (Microsoft)](https://learn.microsoft.com/en-us/azure/architecture/guide/architecture-styles/web-queue-worker)
