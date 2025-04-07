# E-commerce Platform

Микросервисная платформа электронной коммерции, построенная на принципах Clean Architecture.

## Архитектура проекта

Проект состоит из трех основных микросервисов:

1. **api-gateway** (порт 8080)

   - Единая точка входа для всех запросов
   - Маршрутизация запросов к соответствующим микросервисам
   - Проксирование запросов

2. **inventory-service** (порт 8081)

   - Управление товарами
   - CRUD операции для товаров
   - Отслеживание количества товаров

3. **order-service** (порт 8082)
   - Управление заказами
   - Создание и обработка заказов
   - Отслеживание статуса заказов

## Технологии

- Go 1.21+
- Gin Web Framework
- MongoDB Atlas
- Clean Architecture

## Структура каждого микросервиса

Каждый микросервис построен по принципам Clean Architecture и имеет следующую структуру:

```
service/
├── cmd/
│   └── main.go              # Точка входа приложения
├── internal/
│   ├── entity/              # Бизнес-сущности
│   │   └── [entity].go      # Определение сущностей и интерфейсов репозиториев
│   ├── usecase/             # Бизнес-логика
│   │   └── [usecase].go     # Реализация бизнес-правил
│   ├── adapter/             # Адаптеры для внешних сервисов
│   │   └── mongodb/         # Реализация репозиториев для MongoDB
│   └── delivery/            # Доставка (интерфейсы)
│       └── http/            # HTTP-обработчики
├── go.mod                   # Зависимости
└── .env                     # Конфигурация
```

### Принципы Clean Architecture

1. **Entity Layer**

   - Содержит бизнес-сущности и интерфейсы репозиториев
   - Не зависит от других слоев
   - Определяет основные бизнес-правила

2. **UseCase Layer**

   - Содержит бизнес-логику приложения
   - Зависит только от Entity Layer
   - Реализует основные сценарии использования

3. **Adapter Layer**

   - Адаптирует внешние сервисы к интерфейсам репозиториев
   - Реализует интерфейсы, определенные в Entity Layer
   - Изолирует внешние зависимости

4. **Framework Layer**
   - Содержит детали реализации (HTTP, база данных и т.д.)
   - Зависит от всех внутренних слоев
   - Обеспечивает взаимодействие с внешним миром

## Запуск проекта

### Предварительные требования

1. Установите Go 1.21 или выше
2. Создайте кластер MongoDB Atlas
3. Склонируйте репозиторий

### Настройка

1. Создайте файлы `.env` в каждой директории сервиса:

```bash
# api-gateway/.env
PORT=8080
INVENTORY_SERVICE_URL=http://localhost:8081
ORDER_SERVICE_URL=http://localhost:8082

# inventory-service/.env
PORT=8081
MONGODB_URI=your_mongodb_uri
MONGODB_DATABASE=ecommerce

# order-service/.env
PORT=8082
MONGODB_URI=your_mongodb_uri
MONGODB_DATABASE=ecommerce
```

2. Установите зависимости для каждого сервиса:

```bash
cd api-gateway && go mod download
cd inventory-service && go mod download
cd order-service && go mod download
```

### Запуск

Запустите каждый сервис в отдельном терминале:

```bash
# Терминал 1
cd api-gateway && go run cmd/main.go

# Терминал 2
cd inventory-service && go run cmd/main.go

# Терминал 3
cd order-service && go run cmd/main.go
```

## API Endpoints

### Inventory Service

- `POST /api/products` - Создание товара
- `GET /api/products` - Получение списка товаров
- `GET /api/products/:id` - Получение товара по ID
- `PUT /api/products/:id` - Обновление товара
- `DELETE /api/products/:id` - Удаление товара

### Order Service

- `POST /api/orders` - Создание заказа
- `GET /api/orders` - Получение списка заказов
- `GET /api/orders/:id` - Получение заказа по ID
- `PUT /api/orders/:id` - Обновление заказа
- `DELETE /api/orders/:id` - Удаление заказа
- `GET /api/orders/user/:user_id` - Получение заказов пользователя

## Дальнейшее развитие

1. Добавление аутентификации и авторизации
2. Реализация валидации входных данных
3. Добавление логирования
4. Настройка мониторинга
5. Реализация тестов
6. Улучшение обработки ошибок
7. Добавление документации API
