#!/bin/bash
curl -X GET -i http://localhost:8081/api/secure/ping \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjM2NzI3MDgsInVzZXJfaWQiOjF9.M_1hmc8TmXugfizRsOU2luWdMF1byQor3GLaXVmZaxQ"