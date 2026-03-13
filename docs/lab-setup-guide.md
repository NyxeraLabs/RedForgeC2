# Lab Setup Guide

This guide describes setting up a local lab environment for RedForgeC2.

## Prerequisites

- Docker & Docker Compose
- A machine capable of running multiple containers

## Running the Lab

From the repository root:

```sh
docker-compose up --build
```

This starts:

- `postgres`: database
- `teamserver`: C2 server on port 9080
- `ui`: operator UI on port 5174

## Custom Configuration

Environment variables can be set in `docker/docker-compose.yml` or via an `.env` file.

## Resetting the Lab

To rebuild and reset state:

```sh
docker-compose down -v
docker-compose up --build
```
