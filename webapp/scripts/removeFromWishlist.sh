#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/remove-from-wishlist \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ5NjY5NTAsInVzZXJfaWQiOjF9.q0qGcUTon6HcOCesXSLRKd96GrgRnJzlTYfqna2dHNs" \
  -d '{
    "game_id": 1
  }'
