## Инструкция по запуску приложения

### 1. Клонировать репозиторий
```bash
    git clone https://github.com/Jersonmade/todos-app-golang.git
    cd todos-app-golang
```

### 2. Собрать образ
```bash
    docker build -t todos-app .
```

### 3. Запустить сервис
```bash
    docker run -p 8080:8080 todos-app
```

```md
> Приложение будет доступно по адресу: http://localhost:8080
```

## Запуск тестов

### 1. Перейти в директорию с хэндлерами
```bash
    cd .\internal\handlers
```

### 2. Заупстить тесты
```bash
    go test -v
```


