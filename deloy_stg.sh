APP_NAME=go_clean_architecture

docker rm -f ${APP_NAME} #stop container go_clean_architecture
docker rmi -f ${APP_NAME}

docker load -i ${APP_NAME}.tar