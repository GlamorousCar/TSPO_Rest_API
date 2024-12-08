#!/bin/bash

# Store the API base URL
API_URL="http://localhost:8080"

echo -e "\n \n ======Тестирование запросов по 1 практике (CRUD операции)======\n"


echo -e "\nПолучение всех книг :\n"
curl -s -X GET "${API_URL}/books"
sleep 2

echo -e "\n Получение книги по id 1:\n"
curl -s -X GET "${API_URL}/books/1"
sleep 2

echo -e "\nСоздание новой записи книги:\n Запись {
                             \"id\": \"5\",
                             \"title\": \"The Pragmatic Programmer\",
                             \"author\": \"Andy Hunt\"
                                            }\":"
curl -s -X POST "${API_URL}/books" \
  -H "Content-Type: application/json" \
  -d '{
      "id": "5",
      "title": "The Pragmatic Programmer",
      "author": "Andy Hunt"
    }'
sleep 2
echo -e "\nПроверка наличия книги по ID 5 перед удалением:\n"
curl -s -X GET "${API_URL}/books"
sleep 2

echo -e "\nУдаление книги по ID 5:"
curl -s -X DELETE "${API_URL}/books/5"
sleep 2
# Verify deletion by getting all books again
echo -e "\nПодтверждение удаления книги: - получение всех книг::\n"
curl -s -X GET "${API_URL}/books"
sleep 2
#\
#  -H "Authorization: Bearer $TOKEN" | jq '.'

# PUT (update) a book
echo -e "\nОбновление книги по ID 3:"
curl -s -X PUT "${API_URL}/books/3" \
  -d '{
    "title": "The Pragmatic Programmer: 20th Anniversary Edition",
    "author": "Andy Hunt and Dave Thomas"
  }'
sleep 2
echo -e "\nПроверка j,обновления:\n"
curl -s -X GET "${API_URL}/books"
sleep 2
