APP_NAME=go_clean_architecture

# echo "Downloading packages..."
# go mod download

echo "Compiling..."
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

echo "Docker building..."
docker build -t ${APP_NAME} -f ./Dockerfile .
echo "Docker saving..."
docker save -o ${APP_NAME}.tar ${APP_NAME}

echo "Done"