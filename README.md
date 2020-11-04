# Middleware for any API

## TL;DR

In this repository you will find the code for the service
that will help you hide the original domain of any api.
This service only acts as a middleware and simply forwards the request,
which it receives to a third-party API
and transmits the response received from the API.

## Building

To build this project, you must have Docker installed.
If you don't have it installed, just follow the instructions [on this site](
https://www.docker.com/get-started).
You can use the [Dockerfile](Dockerfile) for build this service,
just write following command.
```bash
docker build --tag 'api-middleware'
```

## Running

```bash
docker run --publish 8080:8080 api-middleware \
	-listen_port 8080 \
	-api_domain domain.com
```
