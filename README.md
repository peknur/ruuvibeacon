# ruuvibeacon

Publish github.com/peknur/ruuvitag measurements. Currently supports only http and log publishers.

http publisher requires APP_PUBLISHER_HTTP_URI env variable:
```
APP_PUBLISHER_HTTP_URI="https://some/url/to/post/data" ./ruuvibeacon -tick=60 -output=http
```

ruuvibeacon also provides web view (JSON format):
```
curl http://127.0.0.1:8080/
```
