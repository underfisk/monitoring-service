package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
	. "github.com/underfisk/monitorservice/src"
	"github.com/labstack/echo"
)


var (
	repo *Repository
	userRepo *UserRepository
	jwtSecret = "SOMETHING THAT WILL COME FROM PROCESS.ENV"
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
	userRepo = NewUserRepository(client)

/*	r := mux.NewRouter()
	r.HandleFunc("/log", onLogPost).Methods("POST")
	r.HandleFunc("/log", onLogsRetrieve).Methods("GET")
*/

	e := echo.New()

	//Middlewares
	e.Use(middleware.Recover())

	//Routes
	e.GET("/log", onLogsRetrieve)
	e.POST("/login", login)

	r := e.Group("/")
	r.Use(middleware.JWT([]byte(jwtSecret)))
	r.POST("/log", onLogPost)


	e.Logger.Fatal( e.Start(":4000"))
}

type LogDataDto struct {
	Name string `json:name`
	Env string `json:env`
	Content string `json:content`
	Severity SeverityType `json:severity`
	Type LogType `json:type`

}


func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user,err := userRepo.GetByCredentials(username, password)
	if err != nil {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["id"] = user.Id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"access_token": t,
	})
}

func onLogPost (c echo.Context) error {
	log.Println("We got a new log post to insert")

	data := new(LogDataDto)
	if err := c.Bind(data); err != nil {
		return err
	}

	//Replace with "token" userId
	tokenUserId := uuid.UUID {}

	row, err := repo.CreateLog(Log{
		Name: data.Name,
		Environment: data.Env,
		Type: data.Type,
		Content: data.Content,
		UserId: tokenUserId,
	})

	if err != nil {
		return err
	}

	var res struct {
		Id string `json:id`
	}
	res.Id = row.Id.String()

	return c.JSON(200, &row)
}

func onLogsRetrieve (c echo.Context) error {
	id := c.QueryParam("id")
	parsedId,_ := uuid.Parse(id)

	log.Println( parsedId)
	data, err := repo.FindLog( parsedId )

	if err != nil {
		log.Println("Error retriveing the logs")
		log.Println(err)
		return err
	}

	return c.JSON(http.StatusOK, data)
}
