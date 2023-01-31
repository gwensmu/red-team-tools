package main

import (
	"context"
	"flag"
	"log"
	"net_helpers"

	"github.com/hasura/go-graphql-client"
)

var query struct {
	__Schema struct {
		Types struct {
			Name   string
			Fields struct {
				Name string
				Args struct {
					Name string
					Type struct {
						Name string
					}
				}
			}
		}
	}
}

func Probe(endpoint string) (e error) {
	client := graphql.NewClient(endpoint, nil)

	log.Println("querying", endpoint)

	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("query executed")

	schema := query.__Schema
	log.Println(schema)

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
