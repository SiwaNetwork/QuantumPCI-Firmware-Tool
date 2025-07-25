# Анализ программы Quantum PCI FT

## Обзор
Quantum PCI FT (ранее The Firmware Tool) - это утилита командной строки, написанная на языке Go, предназначенная для добавления специального заголовка к бинарным файлам прошивок. Программа обеспечивает защиту от случайной установки несовместимых прошивок на PCI-устройства.

## Архитектура программы

### Структура проекта
```
quantum-pci-ft/
├── main.go           # Главный файл программы
├── header/          # Пакет для работы с заголовками
│   ├── config.go    # Конфигурация и валидация параметров
│   ├── config_test.go
│   ├── firmware.go  # Основная логика работы с прошивками
│   └── firmware_test.go
├── go.mod           # Модуль Go и зависимости
├── go.sum
└── README.md        # Документация

```

### Основные компоненты

#### 1. Структура заголовка (Header)
Заголовок занимает 16 байт и содержит:
- **Магические байты** (4 байта): 'SHIW' - идентификатор формата
- **PCI Vendor ID** (2 байта): идентификатор производителя устройства
- **PCI Device ID** (2 байта): идентификатор устройства
- **Размер образа** (4 байта): размер прошивки без заголовка
- **HW Revision ID** (2 байта): ревизия оборудования
- **CRC16** (2 байта): контрольная сумма прошивки

Все значения хранятся в сетевом порядке байтов (big-endian) для обеспечения кроссплатформенной совместимости.

#### 2. Конфигурация (Config)
Структура `Config` содержит:
- Параметры командной строки (пути файлов, ID устройств)
- Файловые дескрипторы для входного и выходного файлов
- Флаг `Apply` для фактического создания выходного файла

#### 3. Основной процесс работы

1. **Парсинг аргументов командной строки**
   - Использует стандартный пакет `flag`
   - Проверяет обязательные параметры

2. **Валидация конфигурации**
   - Проверка корректности PCI ID (должны помещаться в uint16)
   - Проверка наличия обязательных параметров

3. **Работа с файлами**
   - Открытие входного файла прошивки
   - Создание выходного файла (если указан флаг -apply)
   - Проверка наличия существующего заголовка

4. **Создание заголовка**
   - Формирование структуры заголовка с указанными параметрами
   - Вычисление размера прошивки

5. **Вычисление CRC16**
   - Использует полином CRC16_ARC
   - Читает файл блоками по 16KB
   - Одновременно копирует данные в выходной файл

6. **Запись результата**
   - Записывает заголовок в начало выходного файла
   - Копирует оригинальную прошивку после заголовка

### Особенности реализации

1. **Обработка ошибок**
   - Использует пакет logrus для логирования
   - Завершает работу при критических ошибках

2. **Производительность**
   - Буферизованное чтение/запись файлов
   - Одновременное вычисление CRC и копирование данных

3. **Безопасность**
   - Проверка существующих заголовков
   - Валидация всех входных параметров
   - Предупреждение о перезаписи заголовков

### Зависимости
- `github.com/sigurn/crc16` - вычисление контрольных сумм
- `github.com/sirupsen/logrus` - логирование
- `github.com/stretchr/testify` - тестирование (dev dependency)

### Применение
Программа используется в процессе подготовки прошивок для PCI-устройств, особенно для Time Card устройств в рамках Open Compute Project. Заголовок позволяет драйверу устройства проверить совместимость прошивки перед её установкой через интерфейс devlink.

## Возможные улучшения
1. Добавление поддержки различных форматов заголовков
2. Реализация верификации прошивки после создания
3. Добавление поддержки пакетной обработки файлов
4. Расширение набора поддерживаемых алгоритмов контрольных сумм
5. Добавление GUI интерфейса для удобства использования