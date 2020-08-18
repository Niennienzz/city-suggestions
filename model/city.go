package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
)

type RedisKey string

const (
	RedisKeyCities   RedisKey = "cities"
	RedisKeyCitiesFT RedisKey = "cities_ft_idx"
)

func (x RedisKey) String() string {
	return string(x)
}

type CityRaw struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lng"`
}

func (x CityRaw) ToCity() *City {
	latitude, err := strconv.ParseFloat(x.Latitude, 64)
	if err != nil {
		panic(err)
	}

	longitude, err := strconv.ParseFloat(x.Longitude, 64)
	if err != nil {
		panic(err)
	}

	return &City{
		Name:      x.Name,
		Country:   x.Country,
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func (x CityRaw) ToGeoLocation() *redis.GeoLocation {
	return x.ToCity().ToGeoLocation()
}

type City struct {
	Name      string  `json:"name"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"-"`
	Longitude float64 `json:"-"`
	Score     int     `json:"score"`
}

func (x City) ToGeoLocation() *redis.GeoLocation {
	return &redis.GeoLocation{
		Name:      x.Name,
		Latitude:  x.Latitude,
		Longitude: x.Longitude,
	}
}

// CitiesFromRediSearchRaw parses the raw response from RediSearch to a slice of City.
// Raw format:
// [
//    total, // number; Total in Redis, but not total returned. Number of returned is controlled by query.
//    city_score_0, // string convertible to number
//    [
//        "name",
//        city_name_0, // string
//        "country",
//        city_country_0, // string
//    ],
//    city_score_1, // string convertible to number
//    [
//        "name",
//        city_name_1, // string
//        "country",
//        city_country_1, // string
//    ]
//    ...
// ]
func CitiesFromRediSearchRaw(records []interface{}) ([]City, error) {
	if len(records) == 0 {
		return nil, errors.New("empty response")
	}

	_, ok := records[0].(int64)
	if !ok {
		return nil, fmt.Errorf("first element in response is not number: %v", records[0])
	}

	records = records[1:]

	var (
		idx = 0
		res = make([]City, 0)
	)

	for idx < len(records) {
		score, ok := records[idx].(string)
		if !ok {
			e := fmt.Errorf("record %d is not score string\n", idx+1)
			return nil, e
		}
		scoreInt, err := strconv.Atoi(score)
		if err != nil {
			e := fmt.Errorf("record %d is not score string convertible to integer: %w\n", idx+1, err)
			return nil, e
		}
		idx += 1
		record, ok := records[idx].([]interface{})
		if !ok {
			log.Printf("TYPE: %T", records[idx])
			e := fmt.Errorf("record %d is not city string slice\n", idx+1)
			return nil, e
		}
		if len(record) != 4 {
			e := fmt.Errorf("record %d is not city string slice of length 4\n", idx+1)
			return nil, e
		}
		if record[0].(string) != "name" || record[2].(string) != "country" {
			e := fmt.Errorf("record %d is not city string slice with 'name' and 'country'\n", idx+1)
			return nil, e
		}
		res = append(res, City{Name: record[1].(string), Country: record[3].(string), Score: scoreInt})
		idx += 1
	}

	return res, nil
}
