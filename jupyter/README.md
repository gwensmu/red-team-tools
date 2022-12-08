# red-team-tools/jupyter

## Use

`go build -o bin/$(basename $(pwd))`

`./bin/jupyter --block 127.0.0.1/30`
`./bin/jupyter --cloud gce --region us-east1`

## Tests

`go test -v`
