package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net_helpers"

	"github.com/juliangruber/go-intersect"
)

var schemaQuery string = `{"query":"query {\n  __schema {\n      types { \n        name \n        fields { \n          name \n          args {\n          name \n            type {\n            \tname\n          }\n        } \n      }\n    } \n  }\n}"}`
var potentiallySensitiveFields []string = []string{"ssn", "tin"}

func NestedKeys(m map[string]interface{}) (keys []string) {
	for k, v := range m {
		if _, ok := v.(map[string]interface{}); ok {
			keys = append(keys, NestedKeys(v.(map[string]interface{}))...)
		} else {
			keys = append(keys, k)
		}
	}

	return keys
}

func lookForFieldsThatSeemSensitive(schema map[string]interface{}) (sensitiveFields []interface{}) {
	types := schema["data"].(map[string]interface{})["__schema"].(map[string]interface{})["types"].([]interface{})
	if len(types) < 1 {
		log.Fatal("No types found in the schema")
	}
	fieldNames := []string{}

	for _, graphqlType := range types {
		fields := graphqlType.(map[string]interface{})["fields"]
		if fields == nil {
			continue
		}

		for _, field := range fields.([]interface{}) {
			fieldNames = append(fieldNames, field.(map[string]interface{})["name"].(string))
		}
	}

	matchedFields := intersect.Simple(fieldNames, potentiallySensitiveFields)

	return matchedFields
}

func sendQuery(endpoint string, query string) (result map[string]interface{}, err error) {
	response, e := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(query)))
	if e != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	schema, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonSchema map[string]interface{}
	err = json.Unmarshal([]byte(schema), &jsonSchema)
	if err != nil {
		log.Fatal(err)
	}

	return jsonSchema, err
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
