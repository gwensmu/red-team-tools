package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/http"
	"net_helpers"
)

var schemaQuery string = `{"query":"query {\n  __schema {\n      types { \n        name \n        fields { \n          name \n          args {\n          name \n            type {\n            \tname\n          }\n        } \n      }\n    } \n  }\n}"}`
var potentiallySensitiveFields []string = []string{"ssn", "tin"}

func lookForFieldsThatSeemSensitive(schema string) (sensitiveFields []string) {
	for _, field := range potentiallySensitiveFields {
		if bytes.Contains([]byte(schema), []byte(field)) {
			sensitiveFields = append(sensitiveFields, field)
		}
	}
	return sensitiveFields
}

func sendQuery(endpoint string, query string) (result string, err error) {
	response, e := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(query)))
	if e != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	resolved, err := io.ReadAll(response.Body)
	schema := string(resolved)

	return schema, err
}

func Probe(endpoint string) (e error) {

	schema, err := sendQuery(endpoint, schemaQuery)
	if err != nil {
		log.Fatal(err)
	}

	if len(lookForFieldsThatSeemSensitive(schema)) > 0 {
		log.Println("Found potentially sensitive fields in the schema:", lookForFieldsThatSeemSensitive(schema))
	} else {
		log.Println("No sensitive fields found in the schema")
	}

	return err
}

func main() {
	net_helpers.InitLogFile("scans", "graphql")

	target := flag.String("target", "http://localhost:8080/query", "a Graphql API endpoint to probe")
	flag.Parse()

	if *target == "" {
		log.Fatal("Please specify a graphql endpoint to probe")
	}

	log.Println("Probing", *target)

	Probe(*target)
}
