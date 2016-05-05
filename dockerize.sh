mkdir ./docker/src/
cp -r ./src/g0Ch@_server ./docker/src/
cp -r ./src/GeneralParser ./docker/src/
cd ./docker
sudo docker build -t g0chat_server .
