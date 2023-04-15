# Azure Container Apps deployment
# Makefile

# Variables
RESOURCE_GROUP_NAME = my-azure-rg
APP_NAME = my-app-name
LOCATION = eastus
IMAGE_NAME = my-image-name
DOCKERFILE_PATH = .
CONTAINER_PORT = 8082

# Build and push the Docker image
docker-build:
    docker build -t $(IMAGE_NAME) $(DOCKERFILE_PATH)

docker-push:
    az acr login --name <your-registry-name>.azurecr.io
    docker tag $(IMAGE_NAME) <your-registry-name>.azurecr.io/$(IMAGE_NAME)
    docker push <your-registry-name>.azurecr.io/$(IMAGE_NAME)

# Terraform deployment
init:
    terraform init

plan:
    terraform plan -var="resource_group_name=$(RESOURCE_GROUP_NAME)" -var="app_name=$(APP_NAME)" -var="location=$(LOCATION)" -var="image_name=$(IMAGE_NAME)" -var="container_port=$(CONTAINER_PORT)"

apply:
    terraform apply -var="resource_group_name=$(RESOURCE_GROUP_NAME)" -var="app_name=$(APP_NAME)" -var="location=$(LOCATION)" -var="image_name=$(IMAGE_NAME)" -var="container_port=$(CONTAINER_PORT)"

destroy:
    terraform destroy -var="resource_group_name=$(RESOURCE_GROUP_NAME)" -var="app_name=$(APP_NAME)" -var="location=$(LOCATION)" -var="image_name=$(IMAGE_NAME)" -var="container_port=$(CONTAINER_PORT)"

# Deploy the app
deploy: docker-build docker-push init apply
