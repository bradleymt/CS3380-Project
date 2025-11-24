#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/join-publisher \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQwMTczNDksInVzZXJfaWQiOjF9.lmEjtTdn8HhgkJ7VrO6pWI1NoX_I-qQ4jnuWCWxgT6A" \
  -d '{
    "studio_name": "ABCDStudios"
  }'