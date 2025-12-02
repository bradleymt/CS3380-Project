#!/bin/bash
curl -X GET -i http://localhost:8081/api/secure/get-cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>"
