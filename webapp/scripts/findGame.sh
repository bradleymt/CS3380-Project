#!/bin/bash
curl -X GET -i "http://localhost:8081/api/secure/find-games?name=&publisher=&price=&currency=&genre=" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQxMzE3OTUsInVzZXJfaWQiOjF9.3je8P6jo1Ycl-YFRuwzZZO-1hZ_dFacfGziRrooN3ew"