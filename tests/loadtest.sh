#!/bin/bash

URL="http://localhost:8080/put"

for i in $(seq 1 10000); do
  key=$(openssl rand -hex 4)
  value=$(openssl rand -hex 8)
  
  curl -s "${URL}?key=${key}&value=${value}" > /dev/null &
done

wait
echo "Done sending 1 lakh requests"
