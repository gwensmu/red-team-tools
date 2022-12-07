# /bin/zsh

# This script is used to run the end-to-end tests for the redis scanner

docker-compose down
docker-compose up -d &
go build -o bin/$(basename $(pwd))
./bin/redis --block 127.0.0.1/30
