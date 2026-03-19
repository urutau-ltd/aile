# REST API Example

Run:

```bash
$ make example-rest
```

Try:

```bash
$ curl http://localhost:9094/api/articles
$ curl -X POST http://localhost:9094/api/articles -H 'content-type: application/json' -d '{"title": "Hello", "body": "World"}'
```
