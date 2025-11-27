#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/create-game \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQxMzE3OTUsInVzZXJfaWQiOjF9.3je8P6jo1Ycl-YFRuwzZZO-1hZ_dFacfGziRrooN3ew" \
  -d '{
    "name": "ABCDGame",
    "price": 0.001,
    "currency": "USD",
    "genre": "Adult Video-Game"
  }'