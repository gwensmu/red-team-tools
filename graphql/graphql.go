package main

import (
	"context"
	"flag"
	"log"
	"net_helpers"

	"github.com/hasura/go-graphql-client"
)

type ArgsQuery struct {
	__Schema struct {
		Types struct {
			Name   graphql.String
			Fields struct {
				Name graphql.String
				Args struct {
					Name graphql.String
					Type struct {
						Name graphql.String
					}
				}
			}
		}
	}
}

func probe(endpoint string) (e error) {
	client := graphql.NewClient(endpoint, nil)

	var argsQuery ArgsQuery
	err := client.Query(context.Background(), &argsQuery, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(argsQuery.__Schema)
	return err
}

func main() {
	net_helpers.InitLogFile("scans", "graphql")

	target := flag.String("target", "http://localhost:8080/graphql", "a Graphql API endpoint to probe")
	flag.Parse()

	if *target == "" {
		log.Fatal("Please specify a graphql endpoint to probe")
	}

	log.Println("Probing", *target)

	probe(*target)
}
