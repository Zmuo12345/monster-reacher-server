cd ../docker
docker image prune --filter="dangling=true" -q
docker compose up --force-recreate --build