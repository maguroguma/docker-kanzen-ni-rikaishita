
```sh
# APIコンテナの起動
docker container run --name api --rm --publish 11323:1323 docker-kanzen-ni-rikai:api

# DBコンテナの起動
docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
```
