#!/bin/sh

sudo docker run -d --name movieapi_db -v postgres_data:/var/lib/postgresql/data -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=movie_db -p 5432:5432 postgres:latest
sudo docker run -d --name movieapi_redis -p 6379:6379 redis:latest
sudo docker build . -t oluwatofunmi/movieapi
sudo docker run -it -e DB_SOURCE=:movieapi_db//postgres:password@postgres:5432/movie_db -p 8081:8081 oluwatofunmi/movieapi