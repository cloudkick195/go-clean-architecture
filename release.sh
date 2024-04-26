APP_NAME=go_clean_architecture

docker rm -f ${APP_NAME}
echo "Docker building..."
docker build -t ${APP_NAME} -f ./Dockerfile .

