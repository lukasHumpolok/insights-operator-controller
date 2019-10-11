package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const API_PREFIX = "/api/v1/"

var apiRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "api_endpoints_requests",
	Help: "The total number requests per API endpoint",
}, []string{"url"})

var apiResponses = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "response_time",
	Help:    "Response time",
	Buckets: prometheus.LinearBuckets(0, 20, 20),
}, []string{"url"})

func countEndpoint(request *http.Request, start time.Time) {
	url := request.URL.String()
	log.Printf("Request URL: %s\n", url)
	duration := time.Since(start)
	log.Printf("Time to serve the page: %s\n", duration)

	apiRequests.With(prometheus.Labels{"url": url}).Inc()

	apiResponses.With(prometheus.Labels{"url": url}).Observe(float64(duration.Microseconds()))
}

func mainEndpoint(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "Hello world!\n")
	countEndpoint(request, start)
}

func getClusters(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	start := time.Now()
	clusters := storage.ListOfClusters()
	json.NewEncoder(writer).Encode(clusters)
	countEndpoint(request, start)
}

func listConfigurationProfiles(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	start := time.Now()
	profiles := storage.ListConfigurationProfiles()
	json.NewEncoder(writer).Encode(profiles)
	countEndpoint(request, start)
}

func getConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "getConfigurationProfile\n")
	countEndpoint(request, start)
}

func setConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "setConfigurationProfile\n")
	countEndpoint(request, start)
}

func changeConfigurationProfile(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "changeConfigurationProfile\n")
	countEndpoint(request, start)
}

func getClusterConfiguration(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	cluster := mux.Vars(request)["cluster"]
	start := time.Now()
	configuration := storage.ListClusterConfiguration(cluster)
	json.NewEncoder(writer).Encode(configuration)
	countEndpoint(request, start)
}

func setClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "setClusterConfiguration")
	countEndpoint(request, start)
}

func enableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "enableClusterConfiguration")
	countEndpoint(request, start)
}

func disableClusterConfiguration(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "disableClusterConfiguration")
	countEndpoint(request, start)
}

func readConfigurationForOperator(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	io.WriteString(writer, "readConfigurationForOperator")
	countEndpoint(request, start)
}

func Initialize(address string, storage storage.Storage) {
	log.Println("Initializing HTTP server at", address)
	router := mux.NewRouter().StrictSlash(true)

	// common REST API endpoints
	router.HandleFunc(API_PREFIX, mainEndpoint)

	// REST API endpoints used by client
	// configuration profiles
	router.HandleFunc(API_PREFIX+"client/configuration_profile", func(w http.ResponseWriter, r *http.Request) { listConfigurationProfiles(w, r, storage) }).Methods("GET")
	router.HandleFunc(API_PREFIX+"client/configuration_profile/{id}", getConfigurationProfile).Methods("GET")
	router.HandleFunc(API_PREFIX+"client/configuration_profile/{id}", changeConfigurationProfile).Methods("PUT")
	router.HandleFunc(API_PREFIX+"client/configuration_profile", setConfigurationProfile).Methods("POST")

	// clusters and its configurations
	router.HandleFunc(API_PREFIX+"client/clusters", func(w http.ResponseWriter, r *http.Request) { getClusters(w, r, storage) }).Methods("GET")
	router.HandleFunc(API_PREFIX+"client/cluster/{cluster}/configuration", func(w http.ResponseWriter, r *http.Request) { getClusterConfiguration(w, r, storage) }).Methods("GET")
	router.HandleFunc(API_PREFIX+"client/cluster/{cluster}/configuration/{id}", setClusterConfiguration).Methods("POST", "PUT")
	router.HandleFunc(API_PREFIX+"client/cluster/{cluster}/configuration/{id}/enable", enableClusterConfiguration).Methods("POST", "PUT")
	router.HandleFunc(API_PREFIX+"client/cluster/{cluster}/configuration/{id}/disable", enableClusterConfiguration).Methods("POST", "PUT")

	// REST API endpoints used by operator
	router.HandleFunc(API_PREFIX+"operator/configuration", readConfigurationForOperator).Methods("GET")

	// Prometheus metrics
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Unable to initialize HTTP server", err)
		os.Exit(2)
	}
}
