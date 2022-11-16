We then need to pull an Elasticsearch Docker image. This will take some time to download.

docker pull docker.elastic.co/elasticsearch/elasticsearch:7.5.2
Now we need to create a Docker volume so that Elasticsearch data doesn’t get lost when a container exits:

docker volume create elasticsearch
The Docker command line to run an Elasticsearch container is quite long, so we will create a script called run-elastic.sh to run the Docker command for us:

#! /bin/bash

docker rm -f elasticsearch
docker run -d --name elasticsearch -p 9200:9200 -e discovery.type=single-node \
    -v elasticsearch:/usr/share/elasticsearch/data \
    docker.elastic.co/elasticsearch/elasticsearch:7.5.2
docker ps
The script needs to be made executable and then run.

chmod +x run-elastic.sh
./run-elastic.sh
Finally, verify that Elasticsearch is running:

curl http://localhost:9200