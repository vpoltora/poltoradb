# Мини-SQL СУБД: архитектура и roadmap реализации

## Что уже реализовано
- SQL-парсер (CREATE / INSERT / SELECT + простые WHERE)
- REPL

---

# Модули, которые нужно реализовать дальше

## 1. Page Manager (Pager)
Низкоуровневый менеджер страниц.

**Задачи:**
- фиксированный размер страницы (например, 4096 байт)
- чтение страницы по `pageID`
- создание новой страницы
- запись страницы на диск
- (позже) кэш страниц

Pager работает только с бинарными страницами.

---

## 2. Page (структура страницы)
Страница — контейнер для строк.

**Структура:**
- header
- slot directory (список `{offset, length}`)
- data region

**Операции:**
- `AddRow(rowBytes) → (slotID, offset)`
- `GetRow(slotID) → []byte`
- проверка свободного места

---

## 3. Row Format (сериализация/десериализация)
Преобразование row ↔ `[]byte`.

**Типы:**
- INT64
- TEXT
- BOOL
- NULL

**Операции:**
- `Serialize(row) → []byte`
- `Deserialize([]byte) → row`

---

## 4. Table Storage
Высокоуровневый слой, работающий со строками.

**Задачи:**
- управлять файлом `.tbl`
- взаимодействовать с Pager
- использовать Page для вставок/чтения
- использовать Row Format

**API:**
- `Insert(row)`
- `Scan() → []Row`
- `GetRow(pageID, slotID)`

---

## 5. Catalog (metadata)
Каталог схем всех таблиц.

**Содержит:**
- список таблиц
- список колонок
- типы колонок
- primary key
- путь к `.tbl` и `.idx`

Минимальная реализация — `catalog.json`.

---

## 6. Query Executor
Связывает AST → операции над таблицами.

**Исполняет:**
- `CREATE TABLE`
- `INSERT`
- `SELECT`

**Использует:**
- Catalog
- Table Storage
- B+Tree index (для SELECT по PK)

**Алгоритмы:**
- INSERT: сериализация → вставка в таблицу → обновление индекса
- SELECT: индекс если PK, иначе полный скан

---

## 7. B+Tree Index (индекс по первичному ключу)
Отдельный файл `.idx`.

**Задачи:**
- хранить `key → (pageID, slotID)`
- быстрый поиск по PK
- вставки + split
- каждая нода = страница

---

## 8. WAL (Write-Ahead Log)
Журнал предзаписи.

**Протокол:**
- `BEGIN`
- `SET_PAGE(pageID, bytes)`
- `COMMIT` + fsync

**Recovery:**
- читать WAL
- применять committed операции
- игнорировать незавершённые
- checkpoint

---

## 9. Transactions (модель MVP)
Простая модель:

- один writer, много readers
- все записи через WAL
- блокировка таблицы или всего движка на время записи

---

## 10. Buffer Pool (LRU-кэш)
Оптимизация после WAL.

**Функции:**
- кэш страниц (64–256 страниц)
- LRU-алгоритм
- dirty pages → flush при commit

---

# Итоговая архитектура

Parser + REPL
↓
Query Executor
↓
Catalog —— Table Storage —— Index (B+Tree)
↓                   ↓
Page Manager ———— WAL —— Buffer Pool
↓
table.tbl / table.idx

---

# Порядок реализации (строгий roadmap)

1. Page Manager  
2. Page (структура, AddRow/GetRow)  
3. Row Format  
4. Table Storage  
5. Catalog  
6. Query Executor  
7. B+Tree Index (PK)  
8. WAL  
9. Transactions (один writer)  
10. Buffer Pool (LRU)
