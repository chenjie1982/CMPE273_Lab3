
package main

import (
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"PutCache",
		"PUT",
		"/keys/{key_id}/{value}",
		PutCache,
	},
	Route{
		"GetKey",
		"GET",
		"/keys/{key_id}",
		GetKey,
	},
	Route{
		"GetKeys",
		"GET",
		"/keys",
		GetKeys,
	},
}

type Data struct {
	Key   int    `json:"key"`
	Value string `json:"value"`
}

func main() {

	router := NewRouter()
	
	cache = make(map[int]string)

	log.Fatal(http.ListenAndServe(":3002", router))

}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		//handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	return router
}

var cache map[int]string


func PutCache(w http.ResponseWriter, r *http.Request) {

	var value string
	var key int

	vars := mux.Vars(r)
	key,_ = strconv.Atoi(vars["key_id"])
	value,_ = vars["value"]//strconv.Atoi(vars["value"])

	cache[key] = value

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	//if err := json.NewEncoder(w).Encode(reply); err != nil {
	//	panic(err)
	//}
}


func GetKey(w http.ResponseWriter, r *http.Request) {

	var key int
	var reply Data

	vars := mux.Vars(r)
	key,_ = strconv.Atoi(vars["key_id"])
	//value,_ = vars["value"]//strconv.Atoi(vars["value"])
	reply.Key = key
	reply.Value = cache[key]

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
}

func GetKeys(w http.ResponseWriter, r *http.Request) {

	var i = 0

	reply := make([]Data, len(cache))
	//fmt.Println("cache:",cache[1])


	for key, value := range cache {
		//fmt.Println("Key:", key, "Value:", value)
    	reply[i].Key = key
    	reply[i].Value = value
    	i++
	}	 
	fmt.Println("GetKeys:",reply)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(reply); err != nil {
		panic(err)
	}
}

