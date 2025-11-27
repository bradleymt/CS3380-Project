#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/remove-from-cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "game_id": 1
  }'
