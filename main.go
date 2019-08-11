package main

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	. "github.com/underfisk/monitorservice/src"
	"github.com/labstack/echo"
	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/joho/godotenv"
	"os"
)


var (
	repo *Repository
	userRepo *UserRepository
	jwtSecret = "SOMETHING THAT WILL COME FROM PROCESS.ENV"
)

func main () {
	print("Running monitoring Service on port: 4000 for now")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret = os.Getenv("JWT_SECRET")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))


	if err != nil {
		log.Fatal("\n Can't connect to mongodb ")
	}

	client.Connect(ctx)
	repo = NewRepository(client)
	userRepo = NewUserRepository(client)


	e := echo.New()

	//Middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://localhost:4000"},
		AllowHeaders: []string{ echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept },
	}))

	e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-XSRF-TOKEN",
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} latency=${latency_human} \n",
	}))
	//e.Use(middleware.Static("./public"))

	//Routes
	e.GET("/logs", onLogsRetrieve, middleware.JWT([]byte(jwtSecret)))
	e.POST("/login", onLogin)
	e.POST("/register", onRegister)
	e.POST("/log", onLogPost, middleware.JWT([]byte(jwtSecret)))


	e.Logger.Fatal( e.Start(":4000"))
}


type LogDataDto struct {
	Name string `json:name`
	Env string `json:env`
	Content string `json:content`
	Severity SeverityType `json:severity`
	Type LogType `json:type`

}


type BadFieldErrorResponse struct {
	Field string	`json:field`
	Message string	`json:message`
}

type CreatedResponse struct {
	Status string `json:status`
}

func onRegister (c echo.Context) error {
	data := new(RegisterDto)
	if err := c.Bind(data); err != nil {
		return err
	}

	//Validate this data
	log.Println(data)
	if usernameErr := validation.Validate(
		data.Username,
		validation.Required,
		validation.Length(3, 50),
		is.Alphanumeric); usernameErr != nil {
		res := &BadFieldErrorResponse {
			Field: "username",
			Message: usernameErr.Error(),
		}
		return c.JSON(400, res)
	}


	if pwdErr := validation.Validate(
		data.Password,
		validation.Required,
		validation.Length(5, 255)); pwdErr != nil {
		res := &BadFieldErrorResponse {
			Field: "password",
			Message: pwdErr.Error(),
		}
		return c.JSON(400, res)
	}

	if emailErr := validation.Validate(data.Email, validation.Required, is.Email); emailErr != nil {
		res := &BadFieldErrorResponse {
			Field: "email",
			Message: emailErr.Error(),
		}
		return c.JSON(400, res)
	}

	if jobPosErr := validation.Validate(data.JobPosition, validation.Required, validation.In("Developer", "CTO", "CEO", "ETC")); jobPosErr != nil {
		res := &BadFieldErrorResponse {
			Field: "job_position",
			Message: jobPosErr.Error(),
		}
		return c.JSON(400, res)
	}

	if companyErr := validation.Validate(data.Company, validation.Required, is.Alpha); companyErr != nil {
		res := &BadFieldErrorResponse {
			Field: "company",
			Message: companyErr.Error(),
		}
		return c.JSON(400, res)
	}

	//If everything reaches here
	log.Println("Data is okay lets register")

	usrId, err := userRepo.Create(data)

	if err != nil {
		return err
	}

	res := &CreatedResponse{
		Status: "ok",
	}

	log.Println("The id is: ", usrId)
	return c.JSON(201, &res)
}

type LoginDto struct {
	Email string `json:email`
	Password string `json:password`
}

func onLogin(c echo.Context) error {
	data := new(LoginDto)
	if err := c.Bind(data); err != nil {
		return err
	}

	log.Println("Logging email: ", data.Email)
	user,err := userRepo.GetByEmail(data.Email)
	if err != nil {
		log.Println("Error getting user: ", err)
		return echo.ErrUnauthorized
	}

	compareError := bcrypt.CompareHashAndPassword(user.Password, []byte(data.Password))
	if compareError != nil {
		log.Println("Invalid password", compareError)
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
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
