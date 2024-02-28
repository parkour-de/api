package graph

import (
	"context"
	"math"
	"pkv/api/src/domain"
	"reflect"
	"testing"
	"time"
)

func TestGetLocations(t *testing.T) {
	cities := map[string]domain.Location{
		"Hamburg": {
			Key: "Hamburg",
			Lat: 53.55,
			Lng: 9.99,
		},
		"Berlin": {
			Key: "Berlin",
			Lat: 52.52,
			Lng: 13.4,
		},
		"Munich": {
			Key: "Munich",
			Lat: 48.14,
			Lng: 11.58,
		},
		"Bremen": {
			Key: "Bremen",
			Lat: 53.07,
			Lng: 8.81,
		},
	}
	db, _, err := Init("../../../config.yml", true)
	if err != nil {
		t.Fatalf("db initialisation failed: %s", err)
	}
	locations := []domain.Location{
		cities["Hamburg"],
		cities["Berlin"],
		cities["Munich"],
		cities["Bremen"],
	}
	for i := range locations {
		err := db.Locations.Create(&locations[i], context.Background())
		if err != nil {
			t.Errorf("initialisation failed: %s", err)
		}
	}
	tests := []struct {
		name    string
		options domain.LocationQueryOptions
		want    []domain.LocationDTO
	}{
		{
			"four locations, sorted by distance",
			domain.LocationQueryOptions{
				Lat:         53.07,
				Lng:         8.81,
				MaxDistance: 1000000,
			},
			[]domain.LocationDTO{
				{
					Location: cities["Bremen"],
					Distance: 0,
				},
				{
					Location: cities["Hamburg"],
					Distance: 141381.42531616212,
				},
				{
					Location: cities["Berlin"],
					Distance: 513898.6834661818,
				},
				{
					Location: cities["Munich"],
					Distance: 621209.2911179818,
				},
			},
		},
		{
			"maxDistance just before Berlin",
			domain.LocationQueryOptions{
				Lat:         53.07,
				Lng:         8.81,
				MaxDistance: 500000,
			},
			[]domain.LocationDTO{
				{
					Location: cities["Bremen"],
					Distance: 0,
				},
				{
					Location: cities["Hamburg"],
					Distance: 141381.42531616212,
				},
			},
		},
		{
			"Pagination",
			domain.LocationQueryOptions{
				Lat:         53.07,
				Lng:         8.81,
				MaxDistance: 1000000,
				Skip:        1,
				Limit:       2,
			},
			[]domain.LocationDTO{
				{
					Location: cities["Hamburg"],
					Distance: 141381.42531616212,
				},
				{
					Location: cities["Berlin"],
					Distance: 513898.6834661818,
				},
			},
		},
		{
			"Only Skip",
			domain.LocationQueryOptions{
				Lat:         53.07,
				Lng:         8.81,
				MaxDistance: 1000000,
				Skip:        1,
				Limit:       0,
			},
			[]domain.LocationDTO{
				{
					Location: cities["Hamburg"],
					Distance: 141381.42531616212,
				},
				{
					Location: cities["Berlin"],
					Distance: 513898.6834661818,
				},
				{
					Location: cities["Munich"],
					Distance: 621209.2911179818,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := db.GetLocations(tt.options, context.Background())
			if err != nil {
				t.Errorf("GetLocations() error = %v", err)
				return
			}
			for i := range got {
				got[i].Created = time.Time{}
				got[i].Modified = time.Time{}
				got[i].Distance = math.Round(got[i].Distance)
			}
			for i := range tt.want {
				tt.want[i].Created = time.Time{}
				tt.want[i].Modified = time.Time{}
				tt.want[i].Distance = math.Round(tt.want[i].Distance)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetLocations()\n  got = %v,\n  want  %v", got, tt.want)
			}
		})
	}
}
