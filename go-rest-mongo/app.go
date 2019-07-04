package mongo

import(
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var client *mongo.Client

type Person struct {
	ID			primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname	string 				`json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname	string 				`json:"lastname,omitempty" bson:"lastname,omitempty"`

}

func main() {
	fmt.Println("Starting application..")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, "mongodb://localhost:27017")
	router := mux.NewRouter()
	router.HandleFunc("/person", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/person/{id}", GetPersonEndpoint).Methods("GET")
	log.Fatal(http.ListenAndServe(":12345", router))

}

func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application-json")
	var person Person
	_ = json.NewDecoder(request.Body).Decode(&person)
	collection := client.Database("go-rest-mongo-1").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application-json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person Person
	collection:= client.Database("go-rest-mongo-1").Collection("people")
	ctx, _= context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "`+ err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(person)

}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request){

}

