git pull

docker build -t ts-balancer .
docker run --publish 3031:3333 ts-balancer