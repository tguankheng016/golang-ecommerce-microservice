## Project Overview

> This is a practical and imaginary eCommerce system built with `Go` as the backend and `React` for the frontend. The system is designed using `Microservices Architecture`, `Vertical Slice Architecture`, and `Clean Architecture` principles. The goal is to demonstrate a scalable, maintainable, and testable approach for building modern eCommerce applications.

This project is extended from [golang-ecommerce-monolith](https://github.com/tguankheng016/golang-ecommerce-monolith) repository.

Other versions of this project are available in these repositories:

- [https://github.com/tguankheng016/golang-ecommerce-monolith](https://github.com/tguankheng016/golang-ecommerce-monolith)
- [https://github.com/tguankheng016/dotnet-commerce-microservice](https://github.com/tguankheng016/dotnet-commerce-microservice)
- [https://github.com/tguankheng016/dotnet-commerce-monolith](https://github.com/tguankheng016/dotnet-commerce-monolith)

OAuth project used:

- [https://github.com/tguankheng016/openiddict-oauth](https://github.com/tguankheng016/openiddict-oauth)

For more details on frontend projects

- [React Admin Portal](https://github.com/tguankheng016/golang-ecommerce-microservice/blob/main/apps/react/README.md)
- [React EShop](https://github.com/tguankheng016/golang-ecommerce-microservice/blob/main/apps/react-eshop/README.md)

# Table of Contents

- [The Goals of This Project](#the-goals-of-this-project)
- [Plan](#plan)
- [Technologies Used](#technologies-used)
- [The Domain and Bounded Context - Service Boundary](#the-domain-and-bounded-context---service-boundary)
- [Quick Start](#quick-start)

## Goals Of This Project

- ✅ Using `Microservices` and `Vertical Slice Architecture` as a high level architecture
- ✅ Using `Event Driven Architecture` with `Nats JetStream` as Message Broker on top of `Watermill` library
- ✅ Using `gRPC` for `internal communication` between microservices
- ✅ Using `PostgresQL` as the relational database.
- ✅ Using `MongoDB` as the NoSQL database.
- ✅ Using `TestContainers` and `testify` for integration testing.
- ✅ Using `Huma` and `Chi` router for handling RESTFul api requests.
- ✅ Using `Goose` to manage migrations.
- ✅ Using `Dependency Injection` and `Inversion of Control` on top of `uber-go/fx` library.
- ✅ Using `Zap` and structured logging
- ✅ Using `Viper` for configuration management
- ✅ Using `YARP` reverse proxy as API Gateway.
- ✅ Using `OpenTelemetry` and `Jaeger` for distributed tracing.
- ✅ Using `OpenIddict` for authentication based on OpenID-Connect and OAuth2.
- ✅ Using `React` to build admin facing and consumer facing application.
- ✅ Using `Azure Key Vault` to manage the secrets.
- ✅ Using `Docker` and `Docker Compose` for deployments.
- ✅ Using `Github Actions` as CI/CD pipeline.
- ✅ Using `AWS EC2` for hosting.

## Plan

> This project is a work in progress, new features will be added over time.

| Feature          | Status       |
| ---------------- | ------------ |
| API Gateway      | Completed ✔️ |
| Identity Service | Completed ✔️ |
| Product Service  | Completed ✔️ |
| Cart Service     | Completed ✔️ |
| Admin Portal     | Completed ✔️ |
| EShop            | Completed ✔️ |

## Technologies Used

- **[`Go`](https://github.com/golang/go)** - The Go programming language.

- **[`Huma`](https://github.com/danielgtaylor/huma)** - Huma REST/HTTP API Framework for Golang with OpenAPI 3.1.

- **[`Goose`](https://github.com/pressly/goose)** - A database migration tool. Supports SQL migrations and Go functions..

- **[`Go-Chi`](https://github.com/go-chi/chi)** - Lightweight, idiomatic and composable router for building Go HTTP services.

- **[`Watermill`](https://github.com/ThreeDotsLabs/watermill)** - Building event-driven applications the easy way in Go..

- **[`Go validating`](https://github.com/RussellLuo/validating)** - A Go library for validating structs, maps and slices.

- **[`Stoplight Elements`](https://github.com/stoplightio/elements)** - Build beautiful, interactive API Docs with embeddable React or Web Components, powered by OpenAPI and Markdown.

- **[`Uber Go Zap`](https://github.com/uber-go/zap)** - Blazing fast, structured, leveled logging in Go.

- **[`Uber Go Fx`](https://github.com/uber-go/fx)** - A dependency injection based application framework for Go.

- **[`OpenIddict`](https://github.com/openiddict/openiddict-core)** - Flexible and versatile OAuth 2.0/OpenID Connect stack for .NET.

- **[`Opentelemetry-Go`](https://github.com/open-telemetry/opentelemetry-go)** - OpenTelemetry Go API and SDK

- **[`Jaeger`](https://github.com/jaegertracing/jaeger)** - About
  CNCF Jaeger, a Distributed Tracing Platform

- **[`gocache`](https://github.com/eko/gocache)** - A complete Go cache library that brings you multiple ways of managing your caches.

- **[`go-redis`](https://github.com/redis/go-redis)** - Redis Go client.

- **[`Copier`](https://github.com/jinzhu/copier)** - Copier for golang, copy value from struct to struct and more.

- **[`Yarp`](https://github.com/microsoft/reverse-proxy)** - Reverse proxy toolkit for building fast proxy servers in .NET.

- **[`MongoDB.Driver`](https://github.com/mongodb/mongo-go-driver)** - The Official Golang driver for MongoDB.

- **[`pgx`](https://github.com/jackc/pgx)** - PostgreSQL driver and toolkit for Go.

- **[`gRPC-go`](https://github.com/grpc/grpc-go)** - The Go language implementation of gRPC. HTTP/2 based RPC

- **[`Nats`](https://github.com/nats-io/nats.go)** - Golang client for NATS, the cloud native messaging system.

- **[`testify`](https://github.com/stretchr/testify)** - A toolkit with common assertions and mocks that plays nicely with the standard library.

- **[`gofakeit`](https://github.com/brianvoe/gofakeit)** - Random fake data generator written in go.

- **[`viper`](https://github.com/spf13/viper)** - Go configuration with fangs.

- **[`azure-sdk`](https://github.com/Azure/azure-sdk-for-go)** - This repository is for active development of the Azure SDK for Go.

- **[`Testcontainers`](https://github.com/testcontainers/testcontainers-go)** - Testcontainers for Go is a Go package that makes it simple to create and clean up container-based dependencies for automated integration/smoke tests.

## The Domain And Bounded Context - Service Boundary

- `Identity Service` - The Identity Service is a bounded context responsible for user authentication and authorization. It handles user creation along with assigning roles and permissions through JWT-based authentication and authorization.

- `Product Service` - The Product Service is a bounded context responsible for handling CRUD operations related to product management.

- `Cart Service` - The Cart Service is a bounded context responsible for handling CRUD operations related to cart management.

## Quick Start

### Prerequisites

Before you begin, make sure you have the following installed:

- **Go 1.23.4**
- **Docker**

Once you have Go 1.23.4 and Docker installed, you can set up the project by following these steps:

Clone the repository:

```bash
git clone https://github.com/tguankheng016/golang-ecommerce-microservice.git
```

Run the development server:

```bash
cd deployments/docker-compose
docker-compose -f docker-compose.yml up -d
```

Once everything is set up, you should be able to access:

- Gateway: [http://localhost:5144](http://localhost:5144)
- Identity Service: [http://localhost:8000/api/v1/docs](http://localhost:8000/api/v1/docs)
- Product Service: [http://localhost:8001/api/v1/docs](http://localhost:8001/api/v1/docs)
- Cart Service: [http://localhost:8002/api/v1/docs](http://localhost:8002/api/v1/docs)
