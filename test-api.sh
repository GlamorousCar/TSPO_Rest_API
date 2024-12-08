#!/bin/bash

# Store the API base URL
API_URL="http://localhost:8080"

echo "Тестирование зарпосов к 1 практике"

## First, let's get a JWT token (assuming we have an endpoint that generates tokens)
#TOKEN=$(curl -s -X POST "${API_URL}/auth/token" \
#  -H "Content-Type: application/json" \
#  -d '{"username": "test", "password": "test"}' | jq -r '.token')
#
#echo "JWT Token: $TOKEN"

# GET all books
echo -e "\nGetting all books:"
curl -s -X GET "${API_URL}/books"
#\
#  -H "Authorization: Bearer $TOKEN" | jq '.'

# GET a specific book
echo -e "\nGetting book with ID 1:"
curl -s -X GET "${API_URL}/books/1"
#\
#  -H "Authorization: Bearer $TOKEN" | jq '.'

# POST a new book
echo -e "\nCreating a new book:"
curl -s -X POST "${API_URL}/books" \
  -H "Content-Type: application/json" \
  -d '{
      "id": "5",
      "title": "The Pragmatic Programmer",
      "author": "Andy Hunt"
    }'
#    | jq '.'
#    \
#  -H "Authorization: Bearer $TOKEN" \
# DELETE a book
echo -e "\nDeleting book with ID 5:"
curl -s -X DELETE "${API_URL}/books/5" \
 -H "Authorization: Bearer $TOKEN"

# Verify deletion by getting all books again
echo -e "\nVerifying deletion - getting all books:"
curl -s -X GET "${API_URL}/books"
#\
#  -H "Authorization: Bearer $TOKEN" | jq '.'

# PUT (update) a book
echo -e "\nUpdating book with ID 4:"
curl -s -X PUT "${API_URL}/books/4" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "The Pragmatic Programmer: 20th Anniversary Edition",
    "author": "Andy Hunt and Dave Thomas"
  }'
#   | jq '.'

