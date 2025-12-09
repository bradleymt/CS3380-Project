#!/bin/bash
# Usage: ./getGameDiscounts.sh <game_id>
# Example: ./getGameDiscounts.sh 1

GAME_ID=${1:-1}

curl -X GET -i "http://localhost:8081/api/secure/game-discounts/${GAME_ID}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQ5NjY5NTAsInVzZXJfaWQiOjF9.q0qGcUTon6HcOCesXSLRKd96GrgRnJzlTYfqna2dHNs"
