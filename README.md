# kafka-golang

A collection of utilities for creating event sourcing applications
with Kafka and Golang.

This is still under heavy development and nowhere near production ready.
Feel free to take any ideas you like, but please mention the original
author (me!). 

## Core Concepts

- Producers and Consumers talk to each other via Avro-encoded messages.
- Avro schemata are versioned.
- Personal data is encrypted and keys may be forgotten.

## Batteries Included

The repository contains a Dockerfile which may be used to create a
changeling container. It includes binaries for running RESTful aggregates,
carrying out administrative tasks etc.

## Hacking

There are examples provided for how to use the library part of this project
in `cmd/example/`. Take a look!

## Components

### Schema Repository

Store Avro schemata in Kafka and retrieve them via HTTP or by using a 
local Consumer.

#### Local Repository

The local repository can be used by creating a `schema.NewLocalRepo()`.
This allows you the encode and decode data using schemata known to the
repository. You may refer to schemata by name and version.

#### CLI

You may create a new schema (read: a version 0 of a schema) by using the
`newschema` program provided in the changeling container.

Upgrading schemata is still under development.

#### RESTful Aggregate

To use the schemata from a different platform, you may query the RESTful
schema repository aggregate. You can start the aggregate with the
`explorer` program provided in the changeling container.

The API spec will be released in the future. In the meantime, take a look
at `cmd/schema/explorer.go` and `schema/explorer/dto.go`.

#### Schema Explorer

There is a web frontend for browsing the schemata, but it is not publicly
available yet. Sorry!

### GDPR "Forget Me"

_in Progress_