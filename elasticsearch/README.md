# red-team-tools

# features to add

* `http://+host+:9200/_cat/indices/my-index-*?v=true&s=index&pretty`
* add flag for number of workers
* keep track of stopping point of ranges scanned, so process can be resumed

## Use

`go build`

`elasticsearch --block 127.0.0.1/30`
`./elasticsearch --cloud gce --region us-east1`

## Tests

`go test -v`
