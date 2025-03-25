echo "start building binaries..."

rm -rf ./docker-binaries/*

docker cp media-server-app:/usr/app ./docker-binaries
cp ./Dockerfile ./docker-binaries/app
cp ./docker-compose.yml ./docker-binaries

echo "done building binaries..."