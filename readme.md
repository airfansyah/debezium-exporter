##build
```
docker build -t airfansyah/debezium-exporter .
docker push airfansyah/debezium-exporter
```

##RUN
```
docker run -e DEBEZIUM_URL="yoururl" -p 9100:9100 airfansyah/debezium-exporter
```