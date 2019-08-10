package MonitoringService

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
	"github.com/underfisk/monitoring-service/core"
)


var (
	repo *Repository
)

/**
	This service collects metrics related with multiple servers/instances/apps
	Also will be used to provide logs monitoring, service monitoring and frontend application
	interface.
	In order to manage the dashboard (visually) we'll be using Vue.js as a spa
	and authorization based using JWT

	Support:
		- gRPC (For fast communication and create nestjs client for this)
		- REST (http support)
 */
func main () {
	print("Running monitoring Service on port: 4000 for now")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))


	if err != nil {
		log.Fatal("\n Can't connect to mongodb ")
	}

	client.Connect(ctx)
	repo = NewRepository(client)

	r := mux.NewRouter()
	r.HandleFunc("/log", onLogPost).Methods("POST")
	r.HandleFunc("/log", onLogsRetrieve).Methods("GET")


	srv := &http.Server{
		Handler:      r,
		Addr:         ":4000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func onLogPost (w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "We got a new log post to insert")
}

func onLogsRetrieve (w http.ResponseWriter, r *http.Request) {
	//cc, _ := context.WithTimeout(context.Background(), 5*time.Second)
/*	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection := mongoClient.Database("top").Collection("emaillogs")
	//count, _ := collection.CountDocuments(context.WithTimeout(context.Background(), 5*time.Second))
	list, err := collection.Find(ctx, bson.D{})

	if err != nil {
		log.Fatal(err)
	}

	defer list.Close(ctx)

	for list.Next(ctx) {
		var result bson.M
		err := list.Decode(&result)
		if err != nil { log.Fatal(err) }

		fmt.Fprintf(w, "%+v\n", result)
	}*/

	res, err := repo.CreateLog(Log{
		applicationId: 1,
		name: "TEST",
		_type: ERROR,
		payload: "SOMETHING",
	})

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	fmt.Fprint(w, res)


	fmt.Fprintf(w, "Displaying the logs")
}
