package configs

import (
	"fmt"
	"log"
	"net/http"
	"os"

	user_routes "github.com/http-crud/api/routes"

	"github.com/joho/godotenv"
)

func LoadEnvVarsAndStartApp() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error while loading ENV %v", err)
	}
	PORT := os.Getenv("PORT")

	mux := http.NewServeMux()

	user_routes.UserRoutes(mux)

	fmt.Printf("Server running on Port %v \n", PORT)

	if err := http.ListenAndServe(PORT, mux); err != nil {
		log.Fatal(err)
	}
}
