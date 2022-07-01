cd ../docker
docker image prune --filter="dangling=true"
docker compose up --force-recreate --build