#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/purchase \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ5NjY5NTAsInVzZXJfaWQiOjF9.q0qGcUTon6HcOCesXSLRKd96GrgRnJzlTYfqna2dHNs" \
  -d '{
    "game_ids": [1],
    "street": "123 Main St",
    "city": "Anytown",
    "state": "CA",
    "zip_code": "12345",
    "country": "USA",
    "payment_method": "card"
  }'
