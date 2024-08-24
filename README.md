# This project is a course project of the 12th stream Route256 Go Middle.

As part of the training project, a system consisting of several services was implemented, which simulates the operation of a simple online store, including such business processes as:
- adding products to the cart and removing them from it
- view the contents of the shopping cart
- placing an order according to the current composition of the basket
- creating an order
- payment of the order
- cancellation of the order by the user or after the expiration of the payment waiting time

Applied technologies in the project:
  • Golang: mutex, graceful shutdown, errgroup, goroutines, channels and others.
  • Unit, e2e tests with minimock - https://github.com/gojuno/minimock
  • Benchmarks
  • gRPC, protubuf, swagger
  • PostgreSQL: master-slave replication, sqlc, transactions, goose for migrations, sharding
  • CI/CD – working with Pipeline in GitLab. The pipeline consists of several stages, during which code linting, Unit run, integration and e2e tests, build and reversal of the entire application take place.
  • Docker, docker-compose
  • Prometheus
  • Grafana
  • Jaeger
  • Apache Kafka