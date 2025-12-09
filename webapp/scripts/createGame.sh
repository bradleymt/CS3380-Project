#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/create-game \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ5NjcxMTEsInVzZXJfaWQiOjF9.iAlAMXTOLCCAEIKbUhSexCclMprWrP65fFdHmnl59vI" \
  -d '{
    "name": "ABCDGame",
    "price": 0.001,
    "currency": "USD",
    "genre": "Adult Video-Game"
  }'