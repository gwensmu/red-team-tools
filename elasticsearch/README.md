# red-team-tools/elasticsearch

## Use

`go build -o bin/$(basename $(pwd))`

`./bin/elasticsearch --block 127.0.0.1/30`
`./bin/elasticsearch --cloud gce --region us-east1`

## Tests

`go test -v`

# features to add

* `http://+host+:9200/_cat/indices/my-index-*?v=true&s=index&pretty`
* add flag for number of workers
* keep track of stopping point of ranges scanned, so process can be resumed
* filter out the node indexers

```
and it's compounded by the insecure defaults
11:08
redis is also another good one to scan for, 6379
11:08
only downside is it's a binary client, not http so you need redis cli to interact
:white_check_mark:
1

11:08
but most of them never have passwords configured
:eyes:
1

11:08
so doing KEYS * is fun
:eyes:
1


```
