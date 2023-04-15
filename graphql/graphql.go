package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"net_helpers"
)

var schemaQuery string = `{"query":"query {\n  __schema {\n      types { \n        name \n        fields { \n          name \n          args {\n          name \n            type {\n            \tname\n          }\n        } \n      }\n    } \n  }\n}"}`
var potentiallySensitiveFields = make(map[string]int)

func buildSensitiveFieldsMap() {
	potentiallySensitiveFields["ssn"] = 1
	potentiallySensitiveFields["tin"] = 1
	potentiallySensitiveFields["taxId"] = 1
}

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

func lookForFieldsThatSeemSensitive(schema map[string]interface{}) (typesWithSensitiveFields []interface{}) {
	buildSensitiveFieldsMap()

	types := schema["data"].(map[string]interface{})["__schema"].(map[string]interface{})["types"].([]interface{})
	if len(types) < 1 {
		log.Fatal("No types found in the schema")
	}

	fieldNames := []string{}
	var sensitiveTypes []interface{}

	for _, graphqlType := range types {
		fields := graphqlType.(map[string]interface{})["fields"]
		if fields == nil {
			continue
		}

		for _, field := range fields.([]interface{}) {
			fieldName := field.(map[string]interface{})["name"].(string)
			fieldNames = append(fieldNames, fieldName)

			_, ok := potentiallySensitiveFields[fieldName]
			if ok {
				sensitiveTypes = append(sensitiveTypes, graphqlType)
			}
		}
	}

	log.Printf("Found %d fields in the schema", len(fieldNames))
	log.Print(fieldNames)

	return sensitiveTypes
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

	sensativeTypes := lookForFieldsThatSeemSensitive(schema)
	if len(sensativeTypes) > 0 {
		log.Println("Types with potentially sensitive fields:", sensativeTypes)
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
