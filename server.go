package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mitri-dvp/go-graphql-challenge/post/schema"

	"github.com/graphql-go/graphql"
)

func executeQuery(query string, variables map[string]interface{}, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:         schema,
		RequestString:  query,
		VariableValues: variables,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var reqBody struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}

	if err := json.Unmarshal(body, &reqBody); err != nil {
		http.Error(w, "Error parsing JSON request body", http.StatusBadRequest)
		return
	}

	query := reqBody.Query
	variables := reqBody.Variables

	result := executeQuery(query, variables, schema.PostSchema)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/graphql", graphqlHandler)

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file")
	}

	port := os.Getenv("PORT")
	fmt.Printf("Now server is running on http://localhost:%v/graphql", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
