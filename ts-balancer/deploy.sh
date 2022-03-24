git pull

docker build -t ts-balancer .
docker run --publish 3333:3031 ts-balancer