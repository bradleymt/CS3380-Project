#!/bin/bash
curl -X GET -i http://localhost:8081/api/secure/library \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ5NjY5NTAsInVzZXJfaWQiOjF9.q0qGcUTon6HcOCesXSLRKd96GrgRnJzlTYfqna2dHNs"
