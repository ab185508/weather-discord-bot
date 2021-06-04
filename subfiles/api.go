package discordbot

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Data struct {
	City string `json:"city"`
	Temp string `json:"temp"`
	Log  string `json:"log"`
}

var Datas []Data

func GetDataEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for _, item := range Datas {
		if item.City == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Data{})
}

func GetDatasEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Datas)
}

func CreateDataEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	var data Data
	_ = json.NewDecoder(req.Body).Decode(&data)
	data.City = params["id"]
	Datas = append(Datas, data)
	json.NewEncoder(w).Encode(Datas)
}

func DeleteDataEndpoint(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	for index, item := range Datas {
		if item.City == params["id"] {
			Datas = append(Datas[:index], Datas[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(Datas)
}

func mux_start() {
	Router := mux.NewRouter()
	Router.HandleFunc("/Datas", GetDatasEndpoint).Methods("GET")
	Router.HandleFunc("/Datas/{id}", GetDataEndpoint).Methods("GET")
	Router.HandleFunc("/Datas", CreateDataEndpoint).Methods("POST")
	Router.HandleFunc("/Datas/{id}", DeleteDataEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8081", Router))
}

func mains(c string, t string, l string) {
	Datas = append(Datas, Data{City: c, Temp: t, Log: l})
}
