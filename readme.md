# nr-code-challenge-lb

This repository contains my code challenge submission for my New Relic SRE Application.

## Milestones

- [X] Reference Proxy - Traefik (Go)
- [X] Simple Fizz and Buzz containers
- [X] Functional docker compose file for reference
- [X] Initial Golang files for Demo
- [X] Accept an http GET
- [X] Forward GET to backend
- [X] Add second backend and enable round robin
- [ ] Add test framework
- [ ] Add tests (Accept a GET request, Forward GET Request, Round Robin)
- [ ] Refactor and add relevant parameters
- [ ] Update Readme.md

### Requirements:

- [X] Application distributes request between 2 or more backend services.
- [X] Application performs round-robin distribution

#### Extras:

- [X] Dynamic
- [X] Healthchecks (dynamic and passive)
- [ ] Sorry Page
- [X] Portable (Go)

## Pre-requisites:

The following components are required to build and run the solution.

* Docker
  * docker-compose
* Golang (optional if using docker)

## How to Use:

Install the dependencies above (for personal use, I recommend Docker Desktop).

Then open up a terminal and:

1. Clone this repository. `git clone https://github.com/th3raid0r/nr-code-challenge-lb.git`
2. Cd into the directory. `cd nr-code-challenge-lb`
3. Create a file named `app.env` and fill it with the following information:
``` 
TARGET_LIST=http://fizz:80,http://buzz:80
PORT=3030
```        
4. Run the docker compose template with the env file. `docker-compose --env-file ./app.env up -d`
5. Create the follwoing [HOSTS](https://www.howtogeek.com/howto/27350/beginner-geek-how-to-edit-your-hosts-file/) file entries*:
```
127.0.0.1 localhost local example.local fizz.example.local buzz.example.local referencelb.example.local demo.example.local
```
6. Navigate to `http://demo.example.local:3030` and observe the "OS: Host" field chance between fizz and buzz between each refresh.
7. If you wish, compare with the traefik reference loadbalancer.
8. Additionally, you can modify `docker-compose.yml` to include another copy of the "echo" instances (fizz and buzz). Just be sure to add the new host to your app.env file!


* \*Note - A default docker desktop configuration will likely bind these domains for you based on our compose file, but it's here if you're configuration is more custom.

## Next Steps:

From here we need test coverage mostly. 

## Attributions

I've cobbled this together from resources found here: 

[Simplelb](https://github.com/kasvith/simplelb)

[Viper](https://github.com/spf13/viper)

## Licenses

The code examples herin falls under no license and may be duplicated without attribution to me.


However our language of choice and runtime selection fall under different licenses:

Docker - [License](https://github.com/moby/moby/blob/master/LICENSE)

Docker Desktop - [License](https://docs.docker.com/subscription/#docker-desktop-license-agreement)

Golang - [License](https://go.dev/LICENSE)

Traefik - [License](https://github.com/traefik/traefik/blob/master/LICENSE.md)

http-https-echo - [License](https://github.com/mendhak/docker-http-https-echo/blob/master/LICENSE.md)

