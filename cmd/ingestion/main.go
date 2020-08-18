package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-redis/redis/v8"

	"city-suggestions/model"
	"city-suggestions/repository"
)

func main() {
	var (
		repo = repository.NewLocal()
		raws []model.CityRaw
	)

	bs, err := ioutil.ReadFile("cities.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(bs, &raws)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	for key, raw := range raws {
		_, err := repo.GeoAdd(ctx, model.RedisKeyCities.String(), raw.ToGeoLocation()).Result()
		if err != nil {
			log.Panicf("Error ingesting city (%d '%s'): %v", key, raw.Name, err)
		}

		scoreFloat, err := repo.ZScore(ctx, model.RedisKeyCities.String(), raw.Name).Result()
		if err != nil {
			log.Panicf("Error reading zscore (%d '%s'): %v", key, raw.Name, err)
		}
		score := int64(scoreFloat)

		searchAdd := func(ctx context.Context, score int64, name string, country string) *redis.SliceCmd {
			cmd := redis.NewSliceCmd(
				ctx,
				"FT.ADD",
				model.RedisKeyCitiesFT.String(), // Key
				score,                           // Score
				"1.0",                           // Weight
				"FIELDS",
				"name", fmt.Sprintf("%s", name),
				"country", fmt.Sprintf("%s", country),
			)
			return cmd
		}(ctx, score, raw.Name, raw.Country)

		err = repo.Process(ctx, searchAdd)
		if err != nil {
			log.Printf("Error processing search add command (%d '%s'): %v", key, raw.Name, err)
		}

		_, err = searchAdd.Result()
		if err != nil {
			log.Printf("Error in search add command result (%d '%s'): %v", key, raw.Name, err)
		}

		log.Printf("Ingested city %d '%s'\n", key, raw.Name)
	}
}
