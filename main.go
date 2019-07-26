package main

import (
	"strconv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
)

var Places []Place

type Place struct {
	Id 		int `json:"id"`
	Country string `json:"country"`
	Desc    string `json:"desc"`
}

type errorMessage struct {
    Message string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to Heidi's Travel Page!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllPlaces)
	myRouter.HandleFunc("/places", createNewPlace).Methods("POST")
	myRouter.HandleFunc("/places/{id}", deletePlace).Methods("DELETE")
	myRouter.HandleFunc("/places/{id}", updatePlace).Methods("PUT")
	myRouter.HandleFunc("/places/{id}", returnSinglePlace)
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	Places = []Place{
		Place{Id: 1, Country: "Portugal", Desc: "First Surfing"},
		Place{Id: 2, Country: "Mexico", Desc: "Fav Country So Far"},
		Place{Id: 3, Country: "Japan", Desc: "Mountain Onsen"},
		Place{Id: 4, Country: "Malta", Desc: "Cliff Jumping"},
		Place{Id: 5, Country: "Croatia", Desc: "Swimming in Cold Water"},
		Place{Id: 6, Country: "Slovenia", Desc: "Good Honey and Surprisingly great pizza"},
		Place{Id: 7, Country: "Hungary", Desc: "Amazing Spa and Bars"},
		Place{Id: 8, Country: "HongKong", Desc: "Crowded and So Many Food Markets"},
		Place{Id: 9, Country: "Iceland", Desc: "Ice Hiking and Northen Lights"},
		Place{Id: 10, Country: "Germany", Desc: "Oktoberfest"},
		Place{Id: 11, Country: "France", Desc: "So Much Cheese and Wine"},
		Place{Id: 12, Country: "Korea", Desc: "Home Country"},
		Place{Id: 13, Country: "Taiwan", Desc: "Such Friendly People and Cute Small Towns"},
		Place{Id: 14, Country: "Turkey", Desc: "Beautiful Valleys, Temples and Spices"},
		Place{Id: 15, Country: "Morocco", Desc: "Camping in Sahara"},
	}
	handleRequests()
}

func returnAllPlaces(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Endpoint Hit: returnAllPlaces")
	json.NewEncoder(w).Encode(Places)
}

func returnSinglePlace(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
	key := vars["id"]
	s, _ := strconv.Atoi(key)

	for _, place := range Places {
        if place.Id == s {
            json.NewEncoder(w).Encode(place)
        }
    }
}

func createNewPlace(w http.ResponseWriter, r *http.Request) {
	place := Place{}
	json.NewDecoder(r.Body).Decode(&place)
	//fmt.Printf("hi: %+v\n", place)

	// validation?
	ok := isIDempty(place)
	if ok {
		// create id for new place
		//fmt.Printf("ok: %+v\n", place)
		place = createID(place)
		Places = append(Places, place)
	} else {
		ok = isIDexist(place)
		if ok {
			//fmt.Printf("not ok: %+v", place)
			errM := errorMessage {
				Message: "ID already exists, please try a different ID",
			}
			//fmt.Println(errM)
			json.NewEncoder(w).Encode(errM)
			//create a structcure for an error
			//build 
		} else {
			Places = append(Places, place)
		}
	}
}

func isIDempty(p Place) bool {
	if p.Id == 0 {
		return true
	}
	return false
}

func isIDexist(p Place) bool {
	for _, ep := range Places {
		if p.Id == ep.Id {
			return true
		}
	}
	return false
}

func createID(p Place) Place {
	var maxID int
	for _, ep := range Places {
		if ep.Id > maxID {
			maxID = ep.Id
			//fmt.Printf("maxID?: %+v\n", maxID)
		}
	}
	return Place{
		Id: maxID+1,
		Country: p.Country,
		Desc: p.Desc,
	}
	//return fmt.Printf("newID?: %+v", place)
} 

func deletePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	deleteID, _ := strconv.Atoi(id)
	for index, place := range Places {
		if place.Id == deleteID {
			//Places.splice(place.Id, 1)
			/*Places = Places.filter(function(obj),
				obj.Id !== deleteID
			);*/
			Places = append(Places[:index], Places[index+1:]...)
		}
		
	}

}

func updatePlace(w http.ResponseWriter, r *http.Request) {
	//getting info from user's put request
	reqBody, _ := ioutil.ReadAll(r.Body)
	var place Place
	json.Unmarshal(reqBody, &place)
	//fmt.Printf("update?: %+v", place)
	//converting id from string to integer
	vars := mux.Vars(r)
	id := vars["id"]
	updateID, _ := strconv.Atoi(id)
	//iterating Places to find the right id
	for i, ep := range Places {
		if ep.Id == updateID {
			newPlace := &Places[i]
			// updating the info
			newPlace.Country = place.Country
			newPlace.Desc = place.Desc
			
		}
		
	}


}