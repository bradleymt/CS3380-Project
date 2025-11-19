#!/bin/bash
curl -X GET -i http://localhost:8081/api/secure/ping \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjM1OTA4MjEsInN0YXJ0IjoxNzYzNTg3MjIxLCJ1c2VyX2lkIjoxfQ.FodqrybE3gYbZPLRAlZ0i__NIYE7oKzgDx-kygcSWTw"