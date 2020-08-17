package model

import (
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
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (x City) ToGeoLocation() *redis.GeoLocation {
	return &redis.GeoLocation{
		Name:      x.Name,
		Latitude:  x.Latitude,
		Longitude: x.Longitude,
	}
}
