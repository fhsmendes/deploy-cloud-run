package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const (
	UrlViaCEP     = "https://viacep.com.br/ws/%s/json/"
	UrlWeatherAPI = "https://api.weatherapi.com/v1/current.json?key=%s&q=%s"
)

type ViaCEP struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro,omitempty"`
}

type WeatherAPI struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type Temperature struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func isValidCEP(cep string) bool {
	matched, _ := regexp.MatchString(`^\d{8}$`, cep)
	return matched
}

func getCityFromCEP(cep string) (string, error) {
	url := fmt.Sprintf(UrlViaCEP, cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var viaCEP ViaCEP
	if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil {
		return "", err
	}

	if viaCEP.Erro || viaCEP.Localidade == "" {
		return "", fmt.Errorf("zipcode not found")
	}

	return viaCEP.Localidade, nil
}

func getTemperature(city string) (float64, error) {
	apiKey := os.Getenv("APIKeyWeather")

	if apiKey == "" {
		return 0, fmt.Errorf("API key is not set")
	}

	encodedCity := url.QueryEscape(city)
	apiUrl := fmt.Sprintf(UrlWeatherAPI, apiKey, encodedCity)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	var weather WeatherAPI
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return 0, err
	}

	return weather.Current.TempC, nil
}

func convertTemperatures(celsius float64) Temperature {
	fahrenheit := celsius*1.8 + 32
	kelvin := celsius + 273

	return Temperature{
		TempC: celsius,
		TempF: fahrenheit,
		TempK: kelvin,
	}
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")

	fmt.Println("Received request for zipcode:", cep)

	if !isValidCEP(cep) {
		fmt.Println("Invalid zipcode:", cep)
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}

	fmt.Println("Valid zipcode:", cep)

	city, err := getCityFromCEP(cep)
	if err != nil {
		fmt.Println("Error getting city from zipcode:", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}

	fmt.Println("City found:", city)

	tempC, err := getTemperature(city)
	if err != nil {
		fmt.Println("Error getting temperature:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error getting temperature"))
		return
	}

	fmt.Println("Temperature in Celsius:", tempC)
	temps := convertTemperatures(tempC)

	fmt.Println("Converted temperatures:", temps)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(temps)
}

func main() {
	// Load .env file if it exists (for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using environment variables")
	}

	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/temperature", temperatureHandler)

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
