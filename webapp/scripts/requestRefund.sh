#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/request-refund \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "purchase_id": 1,
    "reason": "Accidental purchase"
  }'
