package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"city-suggestions/model"
	"city-suggestions/repository"

	"github.com/go-redis/redis/v8"
)

var (
	repo                     = repository.NewLocal()
	errLongitudeMissing      = errors.New("'lng' is required but missing")
	errLongitudeInvalidRange = errors.New("'lng' is not in valid range [-180, 180]")
	errLatitudeMissing       = errors.New("'lat' is required but missing")
	errLatitudeInvalidRange  = errors.New("'lat' is not in valid range [-90, 90]")
)

func main() {
	http.HandleFunc("/city/coord", cityByCoord)
	http.HandleFunc("/city/search", cityBySearch)
	log.Println(http.ListenAndServe(":8000", nil))
}

func cityByCoord(w http.ResponseWriter, r *http.Request) {
	lng, lat, err := parseLngLat(r)
	if err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, err)
		return
	}

	query := new(redis.GeoRadiusQuery)
	query.Radius = 100
	query.Unit = "km"
	query.WithCoord = true
	query.WithDist = true
	query.WithGeoHash = true
	query.Count = 25
	query.Sort = "ASC"

	locations, err := repo.GeoRadius(r.Context(), model.RedisKeyCities.String(), lng, lat, query).Result()
	if err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	logAndWriteResponse(w, r, locations)
}

func cityBySearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		err := errors.New("'q' is required a query parameter")
		logAndWriteError(w, r, http.StatusBadRequest, err)
		return
	}

	var (
		ctx = context.Background()
		lng = new(float64)
		lat = new(float64)
		err error
	)

	*lng, *lat, err = parseLngLat(r)
	if errors.Is(err, errLongitudeMissing) || errors.Is(err, errLatitudeMissing) {
		log.Println("optional 'lng' or 'lat' not set")
		lng, lat, err = nil, nil, nil
	}
	if err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, err)
		return
	}

	search := func(ctx context.Context, query string) *redis.SliceCmd {
		cmd := redis.NewSliceCmd(
			ctx,
			"FT.SEARCH",
			model.RedisKeyCitiesFT.String(), // Key
			query,                           // Query
			"LIMIT", "0", "100",             // Limit
		)
		return cmd
	}(ctx, query)

	err = repo.Process(ctx, search)
	if err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, err)
		return
	}

	// TODO:
	//  - Parse results properly. (Consider using the RediSearch client as well.)
	//  - Find a way to combine the query & the coords search. (Currently lng & lat are not in use.)
	result, err := search.Result()
	if err != nil {
		logAndWriteError(w, r, http.StatusBadRequest, err)
		return
	}

	logAndWriteResponse(w, r, result)
}

func parseLngLat(r *http.Request) (float64, float64, error) {
	const (
		maxLng      = 180.0
		maxLat      = 90.0
		maxLatRedis = 85.05112878
	)

	lngStr := r.URL.Query().Get("lng")
	if lngStr == "" {
		return 0.0, 0.0, errLongitudeMissing
	}

	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		return 0.0, 0.0, err
	}

	// [-180, 180]
	if lng < -maxLng || lng > maxLng {
		return 0.0, 0.0, errLongitudeInvalidRange
	}

	latStr := r.URL.Query().Get("lat")
	if latStr == "" {
		return 0.0, 0.0, errLatitudeMissing
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return 0.0, 0.0, err
	}

	// [-90, 90]
	if lat < -maxLat || lat > maxLat {
		return 0.0, 0.0, errLatitudeInvalidRange
	}

	// [-90, -85.05112878]
	if lat <= -maxLatRedis {
		lat = -maxLatRedis
	}

	// [85.05112878, 90]
	if lat >= maxLatRedis {
		lat = maxLatRedis
	}

	return lng, lat, nil
}

type Error struct {
	Status     int    `json:"status"`
	StatusText string `json:"status_text"`
	Message    string `json:"message"`
}

func logAndWriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	errResp := Error{
		Status:     status,
		StatusText: http.StatusText(status),
		Message:    err.Error(),
	}

	bs, err := json.Marshal(errResp)
	if err != nil {
		log.Printf("%s %s %s\n", r.URL.Path, http.StatusText(http.StatusInternalServerError), err.Error())
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}

	log.Printf("%s %s %s\n", r.URL.Path, http.StatusText(status), errResp.Message)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(bs)
}

func logAndWriteResponse(w http.ResponseWriter, r *http.Request, v interface{}) {
	bs, err := json.Marshal(v)
	if err != nil {
		logAndWriteError(w, r, http.StatusInternalServerError, err)
		return
	}

	log.Printf("%s %s\n", r.URL.Path, http.StatusText(http.StatusOK))
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bs)
}
