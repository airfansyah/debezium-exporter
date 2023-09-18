## build
```
docker build -t airfansyah/debezium-exporter:amd64 -f Dockerfile.amd64 .
docker build -t airfansyah/debezium-exporter:ppc64le -f Dockerfile.ppc64le .

docker push airfansyah/debezium-exporter:amd64
docker push airfansyah/debezium-exporter:ppc64le


docker manifest create airfansyah/debezium-exporter:latest \
    airfansyah/debezium-exporter:amd64 \
    airfansyah/debezium-exporter:ppc64le

docker manifest annotate airfansyah/debezium-exporter:latest airfansyah/debezium-exporter:amd64 --os linux --arch amd64
docker manifest annotate airfansyah/debezium-exporter:latest airfansyah/debezium-exporter:ppc64le --os linux --arch ppc64le

docker manifest push airfansyah/debezium-exporter:latest
```

## RUN
```
docker run -e DEBEZIUM_URL="yoururl" -p 9100:9100 airfansyah/debezium-exporter
```