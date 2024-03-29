# /bin/zsh

# This script is used to run the end-to-end tests for the elasticsearch scanner

cd ..
docker-compose down
docker-compose up -d &
cd elasticsearch
make build
./bin/elasticsearch --block 127.0.0.1/30
