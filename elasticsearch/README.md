# red-team-tools/elasticsearch

## Use

`go build -o bin/$(basename $(pwd))`

`./bin/elasticsearch --block 127.0.0.1/30`
`./bin/elasticsearch --cloud gce --region us-east1`

## Tests

`go test -v`
