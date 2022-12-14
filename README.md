# Elastic-Search
we need to install Elasticsearch. A convenient way of doing this is to use a Docker image containing an already configured Elasticsearch. If you haven’t already got Docker on your machine Install Docker Engine.

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
You should see a JSON object containing details of the server.

### Follow this guide for more information
https://developer.okta.com/blog/2021/04/23/elasticsearch-go-developers-guide