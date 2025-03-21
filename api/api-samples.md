# API Sample requests

## Register

```shell
curl  -X POST \
  'http://localhost:8080/auth/register' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "login": "login",
  "password": "password"
}'
```

## Tasks search

```shell
curl  -X GET \
  'http://localhost:8080/tasks?status=in_progress&due_after=2025-02-02&title=re' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1MjgwNjksInN1YiI6ImxvZ2luIn0.aBizoWzK3Jq4fO7dl6f9pOwX27HY7uOpwUiuksY9EeU'
```

## Task update

```shell
curl  -X PUT \
  'http://localhost:8080/tasks/22222222-2222-2222-2222-222222222222' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1MjgwNjksInN1YiI6ImxvZ2luIn0.aBizoWzK3Jq4fO7dl6f9pOwX27HY7uOpwUiuksY9EeU' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "title": "foo",
  "status": "done",
  "priority": "low",
  "due_date": "2025-02-03"
}'
```

## Tasks export

```shell
curl  -X GET \
  'http://localhost:8080/tasks/export' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDI1MjgwNjksInN1YiI6ImxvZ2luIn0.aBizoWzK3Jq4fO7dl6f9pOwX27HY7uOpwUiuksY9EeU'
```
