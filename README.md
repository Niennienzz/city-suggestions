# City Suggestions

- A simple project trying Geolocation & Full-Text Search capabilities of Redis.

## Requirements

- [Golang](https://golang.org/) v1.14+
- [Make](https://www.gnu.org/software/make/) v3.81+
- [Redis](https://redis.io/) v6.0.6+
- [RediSearch](https://oss.redislabs.com/redisearch/) v1.6+
- [autoComplete.js](https://github.com/TarekRaafat/autoComplete.js) v7.2+

## Run Locally

- It is recommended to use the RediSearch Docker image.

```
zsh> docker run -p 6379:6379 redislabs/redisearch:latest
```

- Create a full-text search index, namely `cities_ft_idx`, in Redis.

```bash
redis-cli> FT.CREATE cities_ft_idx SCHEMA name TEXT WEIGHT 1.0
```

- Run ingestion in shell.

```bash
zsh> make ingest
```

- Run server in shell.

```bash
zsh> make run
```

## API Examples

- Search cities by coordinates.

```
http://localhost:8000/city/coord?lng=1.49129&lat=42.46372
```

- Search cities by a query string, and optionally coordinates.

```
http://localhost:8000/city/search?q=Sant%27Orsola
```
