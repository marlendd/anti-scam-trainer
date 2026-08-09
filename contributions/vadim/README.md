# Личный вклад – Вадим Лезинов

## Проект

**«Антискам-тренажёр»** – учебное веб-приложение для отработки безопасного поведения на классифайде через ветвящиеся диалоги.

- [Исходный репозиторий](https://github.com/marlendd/anti-scam-trainer)
<<<<<<< HEAD
- [Личный репозиторий с материалами вклада](https://github.com/marlendd/anti-scam-trainer-vadim-contribution)
- [Коммиты автора в ветке `develop`](https://github.com/marlendd/anti-scam-trainer/commits/develop/?author=marlendd)
=======
- [Коммиты автора в ветке `main`](https://github.com/marlendd/anti-scam-trainer/commits/main/?author=marlendd)
>>>>>>> feature/add-llm-feedback-recomendation
- GitHub: [@marlendd](https://github.com/marlendd)

## Моя роль

Технический координатор и backend-разработчик на Go. Помимо реализации сценарного движка, слоя хранения сценариев и попыток в PostgreSQL, API прохождения, seed-loader, механизма выдачи наград, интеграционных тестов и части CI, я координировал техническую работу команды и принимал решения по устройству проекта.

## Организационный и управленческий вклад

### Распределение обязанностей

- декомпозировал MVP на отдельные направления: frontend, авторизация, сценарии, попытки, прогресс, инфраструктура и тестирование;
- распределял задачи между участниками с учётом зависимостей между модулями;
- определял порядок интеграции работ, чтобы параллельная разработка не блокировала команду;
- синхронизировал изменения участников при объединении веток и разрешении пересечений.

### Принятие технических решений

- участвовал в определении границ MVP и фиксировал требования к продукту;
- выбрал модульный монолит вместо преждевременного разделения backend на микросервисы;
- определил модель хранения графа сценария в PostgreSQL и правила его строгой валидации;
- зафиксировал правила расчёта результата, идемпотентности ответов и атомарного завершения попытки;
- принимал решения о структуре API и разделении ответственности между модулями `scenario`, `engine`, `attempt` и `progress`;
- определял критерии готовности backend-задач: unit-тесты, PostgreSQL integration-тесты, race detector и lint.

### Контроль интеграции и качества

- проверял совместимость изменений перед объединением веток;
- устранял интеграционные проблемы на стыке API, базы данных, seed-сценариев и тестового окружения;
- организовал проверку полного пользовательского пути на сценарии, загруженном из PostgreSQL;
- включил реальные integration-тесты в CI и обеспечил их последовательный запуск на общей тестовой базе;
- поддерживал актуальность технических требований и следующих шагов команды.

Подтверждающие материалы:

- [`docs/requirements.md`](../../docs/requirements.md) – требования и границы MVP;
- [`ff552da`](https://github.com/marlendd/anti-scam-trainer/commit/ff552da) – первичная фиксация требований;
- [`d383749`](https://github.com/marlendd/anti-scam-trainer/commit/d383749) – уточнение требований по итогам технических решений;
- [история Pull Requests](https://github.com/marlendd/anti-scam-trainer/pulls?q=is%3Apr+author%3Amarlendd) – декомпозиция работы на изолированные изменения и их последовательная интеграция.

## Реализованные части

### 1. Модель и валидация сценариев

- модель ветвящегося сценария с узлами, вариантами ответа и концовками;
- структурная проверка контента;
- проверка достижимости узлов и концовок, циклов и корректности переходов;
- проверка путей, весов ответов, категорий риска и успешных концовок;
- PostgreSQL-схема и repository для версий сценариев.

Код и тесты:

- [`backend/internal/scenario`](../../backend/internal/scenario)
- [`backend/migrations/0002_create_scenarios_attempts_and_answers_tables.up.sql`](../../backend/migrations/0002_create_scenarios_attempts_and_answers_tables.up.sql)

### 2. Движок прохождения и расчёт результата

- чистая функция перехода по графу сценария;
- обработка обычных и завершающих переходов;
- расчёт итоговой оценки с учётом веса каждого ответа;
- unit-тесты успешных, рискованных и ошибочных переходов.

Код и тесты:

- [`backend/internal/engine`](../../backend/internal/engine)

### 3. Попытки и ответы пользователя

- создание, продолжение и перезапуск попытки;
- получение текущего состояния прохождения;
- транзакционная запись ответа и переход к следующему узлу;
- защита от повторного выбора и конфликтующих запросов;
- идемпотентная обработка повторного запроса;
- фиксация результата и завершение попытки одной транзакцией;
- HTTP API и интеграционные тесты на PostgreSQL.

Код и тесты:

- [`backend/internal/attempt`](../../backend/internal/attempt)

### 4. Каталог и загрузка сценариев

- API каталога сценариев по роли buyer/seller;
- строгий JSON seed-loader с отклонением неизвестных полей;
- проверка UUID, уникальности, версии, графа, награды и успешных концовок;
- идемпотентная загрузка всех seed-файлов одной транзакцией;
- четыре готовых сценария для покупателя и продавца.

Код и тесты:

- [`backend/internal/scenario/seed_loader.go`](../../backend/internal/scenario/seed_loader.go)
- [`backend/internal/scenario/seed_loader_integration_test.go`](../../backend/internal/scenario/seed_loader_integration_test.go)
- [`backend/seeds`](../../backend/seeds)

### 5. Награды, тестирование и CI

- атомарная выдача фрагмента пазла за успешную концовку;
- защита инвентаря пользователя от повторной выдачи;
- E2E-проверка прохождения сценария, загруженного из seed-файла;
- настройка `golangci-lint`;
- включение реальных PostgreSQL integration-тестов с `-race` в GitHub Actions.

Код и конфигурация:

- [`backend/internal/attempt/answer_service.go`](../../backend/internal/attempt/answer_service.go)
- [`backend/internal/attempt/answer_service_integration_test.go`](../../backend/internal/attempt/answer_service_integration_test.go)
- [`backend/cmd/api/main_test.go`](../../backend/cmd/api/main_test.go)
- [`.github/workflows/check.yml`](../../.github/workflows/check.yml)

<<<<<<< HEAD
### 6. OpenAPI-контракт backend

Я дополнял общий [`backend/openapi.yaml`](../../backend/openapi.yaml) контрактами реализованных мной сценарных API. В частности, описал:

- запуск новой попытки прохождения сценария;
- продолжение активной попытки и её перезапуск;
- отправку ответа и формат результата перехода;
- каталог сценариев с фильтрацией по роли пользователя;
- получение текущего состояния попытки;
- структуры запросов, ответов и типовые ошибки этих endpoint-ов.

Также я синхронизировал контракт с изменением границ модулей, убрав `best_score` из каталога сценариев после переноса этой ответственности в модуль прогресса.

`openapi.yaml` является общим командным файлом: контракт авторизации в нём подготовлен другим участником. Мой вклад относится к разделам сценариев и попыток и подтверждается коммитами [`ece26d6`](https://github.com/marlendd/anti-scam-trainer/commit/ece26d6), [`fed6c5e`](https://github.com/marlendd/anti-scam-trainer/commit/fed6c5e), [`0acb77f`](https://github.com/marlendd/anti-scam-trainer/commit/0acb77f) и [`f4fbcb4`](https://github.com/marlendd/anti-scam-trainer/commit/f4fbcb4).

=======
>>>>>>> feature/add-llm-feedback-recomendation
## Основные Pull Requests

| PR | Результат |
|---|---|
| [#4 – Scenario engine](https://github.com/marlendd/anti-scam-trainer/pull/4) | Модель сценария и базовая валидация |
| [#6 – Scenario persistence](https://github.com/marlendd/anti-scam-trainer/pull/6) | PostgreSQL-схема сценариев, попыток и ответов |
| [#7 – Scenario transitions](https://github.com/marlendd/anti-scam-trainer/pull/7) | Движок переходов и тесты |
| [#13 – Attempt application service](https://github.com/marlendd/anti-scam-trainer/pull/13) | Создание, продолжение и перезапуск попыток |
| [#18 – Scenario repository](https://github.com/marlendd/anti-scam-trainer/pull/18) | Хранение и чтение сценариев в PostgreSQL |
| [#25 – Answer service](https://github.com/marlendd/anti-scam-trainer/pull/25) | Транзакционная обработка ответов и тесты |
| [#26 – Scenario catalog](https://github.com/marlendd/anti-scam-trainer/pull/26) | Каталог сценариев по роли пользователя |
| [#27 – Attempt state](https://github.com/marlendd/anti-scam-trainer/pull/27) | API текущего состояния прохождения |
| [#28 – Scenario seeds](https://github.com/marlendd/anti-scam-trainer/pull/28) | Seed-loader и четыре игровых сценария |
| [#30 – Seed-backed E2E](https://github.com/marlendd/anti-scam-trainer/pull/30) | Сквозной integration-тест реального сценария |
| [#35 – Fragment granting](https://github.com/marlendd/anti-scam-trainer/pull/35) | Атомарная выдача фрагментов пазла |
| [#36 – CI integration tests](https://github.com/marlendd/anti-scam-trainer/pull/36) | Запуск PostgreSQL-тестов в GitHub Actions |

## Ключевые коммиты

- [`7317c1e`](https://github.com/marlendd/anti-scam-trainer/commit/7317c1e) – модель сценария и структурная валидация;
- [`d05edc6`](https://github.com/marlendd/anti-scam-trainer/commit/d05edc6) – проверка достижимости графа и циклов;
- [`4730be2`](https://github.com/marlendd/anti-scam-trainer/commit/4730be2) – PostgreSQL-схема сценариев и попыток;
- [`dad64f5`](https://github.com/marlendd/anti-scam-trainer/commit/dad64f5) – переходы сценарного движка;
- [`c620e6f`](https://github.com/marlendd/anti-scam-trainer/commit/c620e6f) – расчёт итоговой оценки;
- [`01a1686`](https://github.com/marlendd/anti-scam-trainer/commit/01a1686) – сервис обработки ответа и набор тестов;
<<<<<<< HEAD
- [`ece26d6`](https://github.com/marlendd/anti-scam-trainer/commit/ece26d6) – API попыток и соответствующий OpenAPI-контракт;
- [`fed6c5e`](https://github.com/marlendd/anti-scam-trainer/commit/fed6c5e) – каталог сценариев и его OpenAPI-контракт;
- [`0acb77f`](https://github.com/marlendd/anti-scam-trainer/commit/0acb77f) – синхронизация контракта с границами модуля прогресса;
- [`f4fbcb4`](https://github.com/marlendd/anti-scam-trainer/commit/f4fbcb4) – API и OpenAPI-схема текущего состояния попытки;
=======
>>>>>>> feature/add-llm-feedback-recomendation
- [`0f20a67`](https://github.com/marlendd/anti-scam-trainer/commit/0f20a67) – строгий seed-loader;
- [`09a9119`](https://github.com/marlendd/anti-scam-trainer/commit/09a9119) – E2E-тест прохождения seed-сценария;
- [`4acfe5a`](https://github.com/marlendd/anti-scam-trainer/commit/4acfe5a) – выдача фрагментов пазла;
- [`006a9db`](https://github.com/marlendd/anti-scam-trainer/commit/006a9db) – integration-тесты в CI.

## Как проверить результат

Для проверки backend unit-тестов из корня проекта:

```bash
make test
```

Для полного прогона с PostgreSQL в Docker:

```bash
make test-e2e
```

Для статического анализа:

```bash
make lint
```

История коммитов и Pull Requests в исходном репозитории позволяет сопоставить перечисленные изменения с автором и увидеть полный diff каждой части работы.
