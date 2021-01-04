package fesserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"

	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
)

var prop FesServer

/**
Initializes the server with appropriate properties.
Reads the configuration file of server and stores them in the instance of properites struct.
*/
func initServer(configFile string, prop *FesServer) {
	jsonFile, err := os.Open(configFile)
	defer jsonFile.Close()
	if err != nil {
		log.Fatalln("Error opening server configuration file.", err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln("Problem reading server config file.", err)
	}

	json.Unmarshal(byteValue, prop)

	(*prop).Fesengine.Init("des_input.txt")
}

func appHome(w http.ResponseWriter, r *http.Request) {
	log.Println("FES Home!!")
	fmt.Println("json type:", reflect.TypeOf(simplejson.New()))
}

func encryptFile(w http.ResponseWriter, r *http.Request) {
	var response OperationResponse
	filename, ok := r.URL.Query()["filename"]
	if !ok {
		log.Println("No parameter 'filename' in the URL")
	}
	// var fesengine S_DES.DesEncryptor
	response.Result = prop.Fesengine.EncryptFile(filename[0])
	if response.Result {
		response.Message = "File encrypted successfully"
	} else {
		response.Message = "File encryption failed!"
	}
	// response.SetToSimpleJSON()
	WriteToResponseStream(w, response.SetToSimpleJSON(), "application/json")
}

func decryptFile(w http.ResponseWriter, r *http.Request) {
	var response OperationResponse
	filename, ok := r.URL.Query()["filename"]
	if !ok {
		log.Println("No parameter 'filename' in the URL")
	}
	// var fesengine S_DES.DesEncryptor
	prop.Fesengine.DecryptFile(filename[0])
	response.Result = true
	response.Message = "file decrypted successfully"
	WriteToResponseStream(w, response.SetToSimpleJSON(), "application/json")
}

func encryptFolder(w http.ResponseWriter, r *http.Request) {
	// var response OperationResponse
	var isRecursive bool
	foldername, ok := r.URL.Query()["foldername"]
	if !ok {
		log.Println("No parameter 'foldername' in the URL")
	}
	recursive, ok := r.URL.Query()["recursive"]
	if recursive[0] == "true" {
		isRecursive = true
	} else {
		isRecursive = false
	}
	prop.Fesengine.EncryptFolder(foldername[0], isRecursive)
}

func decryptFolder(w http.ResponseWriter, r *http.Request) {
	var isRecursive bool
	foldername, ok := r.URL.Query()["foldername"]
	if !ok {
		log.Println("No parameter 'foldername' in the URL")
	}
	recursive, ok := r.URL.Query()["recursive"]
	if recursive[0] == "true" {
		isRecursive = true
	} else {
		isRecursive = false
	}
	prop.Fesengine.DecryptFolder(foldername[0], isRecursive)
}

/**
Initializes the routing for REST services
*/
func initServices(router *mux.Router) {
	router.HandleFunc("/", appHome).Methods("GET")
	router.HandleFunc("/encrypt", encryptFile).Methods("GET")
	router.HandleFunc("/decrypt", decryptFile).Methods("GET")
	router.HandleFunc("/encryptFolder", encryptFolder).Methods("GET")
	router.HandleFunc("/decryptFolder", decryptFolder).Methods("GET")
}

func Start() {
	initServer("config.json", &prop)
	router := mux.NewRouter().StrictSlash(true)
	initServices(router)
	port := ":" + strconv.Itoa(prop.Port)
	log.Println("Starting application on port: ", prop.Port)
	log.Fatalln(http.ListenAndServe(port, router))
}
