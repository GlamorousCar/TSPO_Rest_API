#!/bin/bash

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

echo -e "\nПодтверждение удаления книги: - получение всех книг::\n"
curl -s -X GET "${API_URL}/books"
sleep 2

echo -e "\nОбновление книги по ID 3:"
curl -s -X PUT "${API_URL}/books/3" \
  -d '{
    "title": "The Pragmatic Programmer: 20th Anniversary Edition",
    "author": "Andy Hunt and Dave Thomas"
  }'
sleep 2
echo -e "\nПроверка обновления книги:\n"
curl -s -X GET "${API_URL}/books"
sleep 2


echo -e "\n \n ======Тестирование запросов по 2 практике (JWT токены)======\n"
sleep 4
echo "Регистрация нового пользователя:"


curl -X POST ${API_URL}/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "testuser4", "password": "testpass"}'

echo "Авторизация для получение токена:"
sleep 3


response=$(curl -X POST ${API_URL}/auth/login \
   -H "Content-Type: application/json" \
   -d '{"username": "testuser4", "password": "testpass"}')

# Парсинг access_token из ответа
access_token=$(echo "$response" | jq -r '.access_token')

echo "Access Token: $access_token"

echo "Использование токена для защищенного эндпоинта\n"

curl -X GET http://localhost:8080/books_with_auth \
  -H "Authorization: Bearer ${access_token}"


echo -e "\n \n ======Тестирование запросов по 4 практике (Пагинация)======\n"

echo -e "\nПолучение книг с пагинацией\n"

curl "http://localhost:8080/books?page=1&pageSize=5"

echo -e "\nПолучение книг с фильтрацией\n"

curl "http://localhost:8080/books?title=Go&author=Donovan"

echo -e "\nПолучение книг с сортировкой\n"

curl "http://localhost:8080/books?sort=title&order=desc"

