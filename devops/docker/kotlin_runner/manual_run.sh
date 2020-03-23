#!/bin/bash
LOCAL_INPUT_MOUNT=$1

CONTAINER_NAME=$(uuidgen)
DATA_VOLUME_NAME="data-${CONTAINER_NAME}"

docker build --tag registry.gitlab.com/grchive/grchive/kotlin_runner:latest .

docker volume create --driver local \
    --opt type=tmpfs \
    --opt device=tmpfs \
    --opt o=size=1g \
    ${DATA_VOLUME_NAME}

docker run \
    --name ${CONTAINER_NAME} \
    --mount type=bind,source=$LOCAL_INPUT_MOUNT,target=/input,readonly \
    --mount source=datavolume,target=/data \
    registry.gitlab.com/grchive/grchive/kotlin_runner:latest \
    --script /input/main.kt

docker cp ${CONTAINER_NAME}:/data/script.jar ${LOCAL_INPUT_MOUNT}/script.jar
docker rm ${CONTAINER_NAME}
docker volume rm ${DATA_VOLUME_NAME}
