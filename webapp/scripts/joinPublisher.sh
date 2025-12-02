#!/bin/bash
curl -X POST -i http://localhost:8081/api/secure/join-publisher \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQxMzE3OTUsInVzZXJfaWQiOjF9.3je8P6jo1Ycl-YFRuwzZZO-1hZ_dFacfGziRrooN3ew" \
  -d '{
    "studio_name": "ABCDStudios"
  }'