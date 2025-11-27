#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/approve-refund \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "refund_id": 1
  }'
