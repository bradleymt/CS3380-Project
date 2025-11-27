#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/join-family \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQwMDU2MTIsInVzZXJfaWQiOjF9.R0Z7z3ycOCCAef7zMK_6pfSTBHts3vgiVoDdChZwhAs" \
  -d '{
    "family_id": 2
  }'