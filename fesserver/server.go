package fesserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"

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

/**
Initializes the routing for REST services
*/
func initServices(router *mux.Router) {
	router.HandleFunc("/", appHome).Methods("GET")
	router.HandleFunc("/encrypt", encryptFile).Methods("GET")
	router.HandleFunc("/decrypt", decryptFile).Methods("GET")
}

func Start() {
	initServer("config.json", &prop)
	router := mux.NewRouter().StrictSlash(true)
	initServices(router)
	log.Fatalln(http.ListenAndServe(":8080", router))
}
