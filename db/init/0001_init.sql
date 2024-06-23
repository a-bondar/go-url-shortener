-- Создание пользователя
CREATE USER gopher WITH PASSWORD 'gopher';

-- Создание базы данных
CREATE DATABASE go_advanced_praktikum
    OWNER gopher
    ENCODING 'UTF8'
    LC_COLLATE = 'en_US.utf8'
    LC_CTYPE = 'en_US.utf8';