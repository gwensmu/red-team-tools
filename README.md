# Red Team Tools

For e2e tests

```
./elasticsearch/e2e.sh
./redis/e2e.sh
```

# Using

`go build -o bin/$(basename $(pwd))`

`./bin/red-team-tools --block 127.0.0.1/30`
`./bin/red-team-tools --cloud gce --region us-east1`


# Todo
refactor the runner, so that each host can be pingged for each open service
and setup separate logfiles for each

# features to add

* add flag for number of workers
* keep track of stopping point of ranges scanned, so process can be resumed
