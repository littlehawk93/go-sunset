package gosunset

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const baseURL string = "https://api.sunrise-sunset.org/json"

// Result an API result returned from the sunrise-sunset.org API
type Result struct {
	Sunrise                   time.Time `json:"sunrise"`
	Sunset                    time.Time `json:"sunset"`
	DayLength                 int       `json:"day_length"`
	SolarNoon                 time.Time `json:"solar_noon"`
	CivilTwilightBegin        time.Time `json:"civil_twilight_begin"`
	CivilTwilightEnd          time.Time `json:"civil_twilight_end"`
	NauticalTwilightBegin     time.Time `json:"nautical_twilight_begin"`
	NauticalTwilightEnd       time.Time `json:"nautical_twilight_end"`
	AstronomicalTwilightBegin time.Time `json:"astronomical_twilight_begin"`
	AstronomicalTwilightEnd   time.Time `json:"astronomical_twilight_end"`
}

type results struct {
	Result *Result `json:"results"`
	Status string  `json:"status"`
}

// GetResult retrieve a new result from sunrise-sunset.org API. Returns the API result or any errors encountered.
func GetResult(lat, lng float32) (*Result, error) {

	u, _ := url.Parse(baseURL)

	q := u.Query()
	q.Add("lat", fmt.Sprintf("%10.7f", lat))
	q.Add("lng", fmt.Sprintf("%10.7f", lng))
	q.Add("formatted", fmt.Sprintf("%d", 0))
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())

	if err != nil {
		return nil, fmt.Errorf("Error retrieving API: %s", err.Error())
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("Response body is nil")
	}

	defer resp.Body.Close()

	var apiResults results

	if err := json.NewDecoder(resp.Body).Decode(&apiResults); err != nil {
		return nil, fmt.Errorf("Unable to parse API response: %s", err.Error())
	}

	if apiResults.Status != "OK" {
		return nil, fmt.Errorf("Error in API response. Status: %s", apiResults.Status)
	}

	if apiResults.Result == nil {
		return nil, fmt.Errorf("No results returned by API")
	}

	return apiResults.Result, nil
}

// ToLocal returns this API results with all times in local timezone time
func (me Result) ToLocal() *Result {

	return &Result{
		Sunrise:                   me.Sunrise.Local(),
		Sunset:                    me.Sunset.Local(),
		DayLength:                 me.DayLength,
		SolarNoon:                 me.SolarNoon.Local(),
		CivilTwilightBegin:        me.CivilTwilightBegin.Local(),
		CivilTwilightEnd:          me.CivilTwilightEnd.Local(),
		NauticalTwilightBegin:     me.NauticalTwilightBegin.Local(),
		NauticalTwilightEnd:       me.NauticalTwilightEnd.Local(),
		AstronomicalTwilightBegin: me.AstronomicalTwilightBegin.Local(),
		AstronomicalTwilightEnd:   me.AstronomicalTwilightEnd.Local(),
	}
}
