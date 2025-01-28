# File Server
Сервер с загрузкой-выгрузкой файлов и разными SQL-ными операциями по выбору этих файлов.
Этот проект представляет собой RESTful API, который позволяет пользователям загружать, скачивать, искать и управлять файлами, хранящимися в базе данных. Поддерживаются различные операции, такие как поиск файлов по имени, получение файлов по ID, создание, обновление и удаление файлов, а также подсчет количества файлов.
Сервер развернут в docker-контейнер, что обеспечивает изоляцию, портативность и упрощает управление.

## Основные возможности
 - Загрузка и скачивание файлов: Файлы можно скачивать по их уникальному ID или имени.

 - Поиск файлов: Поддержка поиска файлов по имени с использованием SQL-запросов.

 - Управление файлами: Возможность обновления и удаления файлов.

 - Подсчет файлов: Получение общего количества файлов или количества файлов, соответствующих определенному критерию.

 - ZIP-архивирование: Возможность скачивания нескольких файлов в виде ZIP-архива.

## Технологии
Язык: Go

База данных: PostgreSQL

Библиотеки:
- net/http для обработки HTTP-запросов.
- database/sql для работы с базой данных и написания запросов с динамической типизацией.
- github.com/google/uuid для генерации уникальных идентификаторов.
- archive/zip для создания ZIP-архивов.

Инструмены: sqlc (для генерации type-safe кода для SQL-запросов), Docker Compose для осуществления взаимодействия между базой данных, сервером, Docker для контейнеризации.
## Установка и запуск
Установить Docker Desktop по оф. гайду [ https://docs.docker.com/engine/install/ ] (или Docker Engine [ https://docs.docker.com/engine/install/ ], если ОС не поддерживает виртуализацию, как, например, виртуальная машина)
>git clone https://github.com/wisp167/file-server.git

В директории проекта

Скачивание пакетов (vendor решил на всякий случай оставить, так что не обязательно):
>go get .

Генерация кода для статически-типизированных SQL-запросов:
>sqlc generate

Установка образа базы данных и сборка контейнеров по Dockerfile и docker-compose.yml
>docker compose up -d --build

API доступно по адресу http://localhost:8000 (по дефолту)
## Примеры запросов
Загрузка на сервер
>curl -X POST -F "file=@path/to/your/file.extension" http://localhost:8000/v1/upload

Загрузка с сервера по имени
>curl -O -J  http://localhost:8000/v1/file_by_name?name=file.ext

Загрузка с сервера по id
>curl -O -J http://localhost:8000/v1/download?id=file-id

Количество файлов с таким именем
>curl http://localhost:8000/v1/count_by_name?name=file.ext

и т.д.
Все виды запросов можно посмотреть в ../docker/sql_init/queries/queries.sql
