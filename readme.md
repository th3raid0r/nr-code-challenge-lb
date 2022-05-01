# nr-code-challenge-lb

This repository contains my code challenge submission for my New Relic SRE Application.

## Milestones

[X] Reference Proxy - Traefik (Go)
[X] Simple Fizz and Buzz containers
[X] Functional docker compose file for reference
[ ] Initial Golang files for Demo
[ ] Accept an http GET
[ ] Forward GET to backend
[ ] Add second backend and enable round robin
[ ] Add test framework
[ ] Add tests (Accept a GET request, Forward GET Request, Round Robin)
[ ] Refactor and add relevant parameters
[ ] Update Readme.md


## Pre-requisites:

The following components are required to build and run the solution.

* Docker
  * docker-compose
* Golang (optional if using docker)

## How to Use:

Run the provided docker-compose file. It currently provides a reference load balancer with two identical backends to provide a functional example of what the solution should look like. 

From here, we will develop our GO application to perform the same function. This will immediately be wrapped in a `Dockerfile` and added as a new service in `docker-compose.yml`, likely before it is fully functional. 

## Licenses

The code examples herin falls under no license and may be duplicated without attribution to me.

However our language of choice and runtime selection fall under different licenses:

Docker - [License](https://github.com/moby/moby/blob/master/LICENSE)
Docker Desktop - [License](https://docs.docker.com/subscription/#docker-desktop-license-agreement)
Golang - [License](https://go.dev/LICENSE)

