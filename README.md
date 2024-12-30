# Сервис передачи файлов по UDP

Этот проект реализует легковесный сервис передачи файлов по протоколу UDP, разработанный для эффективной обработки одновременных передач файлов от нескольких клиентов.

## Возможности
- **Передача файлов**: Отправка файлов любого типа по UDP.
- **Передача по частям**: Файлы разбиваются на части для эффективной передачи.
- **Проверка целостности**: Обеспечивает целостность файлов с помощью контрольных сумм (MD5-хеш).
- **Одновременная работа с клиентами**: Поддерживает одновременную отправку файлов от нескольких клиентов на сервер.
- **Избежание конфликтов имен файлов**: Автоматически переименовывает файлы с одинаковыми именами.
- **Обработка таймаутов**: Очищает память и ресурсы, если передача файла прерывается.

## Структура проекта
```
.
├── LICENSE.txt
├── README.md               # Описание проекта
├── go.mod                  
├── go.sum                  
├── receiver.go             # Приёмник
└── transmitter
    ├── doo                 # Тестовый файл 1
    ├── file.txt            # Тестовый файл 2
    ├── funny_monkey.jpg    # Тестовый файл 3
    ├── important.jpg       # Тестовый файл 4
    └── transmitter.go      # Передатчик
```

## Настройка и использование

### Сборка
Скомпилируйте приложения передатчика и приемника с помощью Go:
```bash
# Компиляция приемника
go build -o receiver receiver.go

# Компиляция передатчика
cd transmitter
go build -o transmitter transmitter.go
```

### Запуск 

#### Запуск приемника:
```bash
./receiver <порт>
```
Замените `<порт>` на желаемый номер порта (например, `8080`).

#### Запуск передатчика:
```bash
./transmitter <адрес> <файл>
```
- `<адрес>`: IP-адрес и порт сервера (например, `127.0.0.1:8080`).
- `<файл>`: Файл, который вы хотите отправить (файл должен быть в одной папке с передатчиком).

### Пример работы
1. Запустите приемник:
   ```bash
   ./receiver 8080
   ```
2. Отправьте файл с клиента:
   ```bash
   ./transmitter 127.0.0.1:8080 important.jpg
   ```

Сервер подтвердит успешную передачу, и файл будет сохранен на сервере с уникальным именем, если это необходимо.

## Используемые технологии
- **Go 1.19**: Основной язык программирования.
- **UUID v1.6.0**: Для генерации уникальных идентификаторов для управления сессиями клиентов.
- **Библиотека Net**: Встроенная библиотека Go для работы с сетью и обработки UDP-соединений.

## Дополнительные примечания
- Убедитесь, что настройки сети позволяют передачу UDP-трафика на указанный порт.
- Тестовые файлы, предоставленные в директории `transmitter`, можно использовать для проверки функциональности.

Лицензия
---
[MIT License](LICENSE.txt) — свободное использование, изменение и распространение проекта.