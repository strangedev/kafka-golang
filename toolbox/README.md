# Kafka Toolbox

The Kafka Toolbox is a Docker container which has a golang environment,
and `librdkafka`installed.
It may be used to run, test and build event sourcing applications with 
the `github.com/strangedev/kafka-golang` go module.
It can also be used as a build stage for building application containers
based on `github.com/strangedev/kafka-golang`.

## Building

(from project's root directory)

```shell script
docker build -t kafka-golang-toolbox toolbox
```

## Development

The `docker-compose.yml` file provided may be used to set up a toolbox
with a volume containing the `github.com/strangedev/kafka-golang` source.

```shell script
docker-compose -f toolbox/docker-compose.yml up -d
docker exec -it toolbox_toolbox_1 ash
```
