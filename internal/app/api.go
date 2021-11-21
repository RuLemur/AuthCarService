package app

import (
	"AuthCarService/internal/controller"
	"context"
	"encoding/json"
	"fmt"
	"github.com/RuLemur/CarService/pkg/endpoint"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// ValidBearer is a hardcoded bearer token for demonstration purposes.
const ValidBearer = "123456"

type Api struct {
	ctx        context.Context
	grpcClient endpoint.CarServiceClient
}

func NewApi(grpcClient endpoint.CarServiceClient) *Api {
	return &Api{grpcClient: grpcClient, ctx: context.Background()}
}

func (a *Api) jsonResponse(w http.ResponseWriter, data interface{}, c int) {
	dj, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}

func (a *Api) parseRequest(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	return nil
}

// RequireAuthentication is an example middleware handler that checks for a
// hardcoded bearer token. This can be used to verify session cookies, JWTs
// and more.
func (a *Api) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Make sure an Authorization header was provided
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")
		// This is where token validation would be done. For this boilerplate,
		// we just check and make sure the token matches a hardcoded string
		if token != ValidBearer {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		// Assuming that passed, we can execute the authenticated handler
		next.ServeHTTP(w, r)
	})
}

// @title Test Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securitydefinitions.oauth2.application OAuth2Application
// @tokenUrl https://example.com/oauth/token
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.implicit OAuth2Implicit
// @authorizationUrl https://example.com/oauth/authorize
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl https://example.com/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants read and write access to administrative information

// NewRouter @securitydefinitions.oauth2.accessCode OAuth2AccessCode
// @tokenUrl https://example.com/oauth/token
// @authorizationUrl https://example.com/oauth/authorize
// @scope.admin Grants read and write access to administrative information
func (a *Api) NewRouter() http.Handler {
	r := gin.Default()

	//r.Use(a.RequireAuthentication)

	c := controller.NewController(a.ctx, a.grpcClient)

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("/user/{id}", c.GetUser)
			users.POST("/user/", c.AddUser)
		}
		cars := v1.Group("/cars")
		{
			cars.POST("/search", c.CarSearch)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Register the API routes

	return r
}
