TRUNCATE TABLE users, scenario_versions, attempts, answers, user_inventory, sessions, password_reset_tokens CASCADE;

-- ==========================================
-- 1. Создаем 3-х пользователей
-- ==========================================
INSERT INTO users (id, email, password_hash) VALUES
('11111111-1111-1111-1111-111111111111', 'alice_pro@gmail.com', 'hash_1'),
('22222222-2222-2222-2222-222222222222', 'bob_noob@yandex.ru', 'hash_2'),
('33333333-3333-3333-3333-333333333333', 'charlie_new@mail.ru', 'hash_3');

-- ==========================================
-- 2. Создаем 2 сценария (за Покупателя и Продавца)
-- ==========================================
INSERT INTO scenario_versions (id, logical_id, version, role, title, description, is_active, content, reward_fragment_id) VALUES
-- Сценарий 1 (Покупатель): Покупка видеокарты. Категории: Фишинг, Предоплата. Дает фрагмент "frag-gpu"
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'aaaaaaaa-1111-aaaa-aaaa-aaaaaaaaaaaa', 1, 'buyer', 'Покупка RTX 4090', 'Покупаем видеокарту ниже рынка', true, '{"risk_categories": ["Фишинг", "Предоплата"]}'::jsonb, 'frag-gpu'),

-- Сценарий 2 (Продавец): Продажа ноутбука. Категории: Поддельная доставка, Фишинг. Дает фрагмент "frag-laptop"
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'bbbbbbbb-2222-bbbb-bbbb-bbbbbbbbbbbb', 1, 'seller', 'Продажа ноутбука', 'Покупатель кидает левую ссылку на доставку', true, '{"risk_categories": ["Поддельная доставка", "Фишинг"]}'::jsonb, 'frag-laptop');

-- ==========================================
-- 3. Создаем попытки прохождения (Attempts)
-- ==========================================
INSERT INTO attempts (id, user_id, scenario_id, status, current_node_id, ending_id, score, started_at, completed_at) VALUES
-- Алиса: Завершила Сценарий 1 ДВА ДНЯ НАЗАД (чтобы у нее был prev_rank) на 90 баллов
('c1111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'completed', NULL, 'ending-safe', 90, now() - interval '2 days', now() - interval '2 days'),
-- Алиса: Сейчас проходит Сценарий 2
('c1111112-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'in_progress', 'node-2', NULL, NULL, now(), NULL),

-- Боб: Завершил Сценарий 1 СЕГОДНЯ на 60 баллов
('c2222221-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'completed', NULL, 'ending-scam', 60, now(), now()),
-- Боб: Завершил Сценарий 2 СЕГОДНЯ на 100 баллов
('c2222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'completed', NULL, 'ending-safe', 100, now(), now()),

-- Чарли: Только начал Сценарий 1
('c3333331-3333-3333-3333-333333333333', '33333333-3333-3333-3333-333333333333', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'in_progress', 'node-1', NULL, NULL, now(), NULL);

-- ==========================================
-- 4. Заполняем ответы (Answers) для расчета весов
-- ==========================================
-- (Тут строгие проверки choice_score и risk_categories на уровне БД)
INSERT INTO answers (attempt_id, node_id, choice_id, idempotency_key, weight, choice_score, risk_categories, consequence, explanation, response) VALUES
-- Алиса (попытка 1 - успешная)
('c1111111-1111-1111-1111-111111111111', 'node-1', 'choice-a', gen_random_uuid(), 3, 100, '[]'::jsonb, 'Всё ок', 'Правильно', '{"ok": true}'::jsonb),
('c1111111-1111-1111-1111-111111111111', 'node-2', 'choice-b', gen_random_uuid(), 2, 50, '["Переход по ссылке"]'::jsonb, 'Опасно', 'Риск', '{"ok": true}'::jsonb),

-- Боб (попытка 1 - плохая)
('c2222221-2222-2222-2222-222222222222', 'node-1', 'choice-c', gen_random_uuid(), 3, 0, '["Скам", "Предоплата"]'::jsonb, 'Потерял деньги', 'Плохо', '{"ok": false}'::jsonb),

-- Чарли (в процессе - ответил на 1 вопрос)
('c3333331-3333-3333-3333-333333333333', 'node-1', 'choice-a', gen_random_uuid(), 2, 100, '[]'::jsonb, 'Всё ок', 'Молодец', '{"ok": true}'::jsonb);

-- ==========================================
-- 5. Инвентарь (Собранные пазлы)
-- ==========================================
INSERT INTO user_inventory (user_id, scenario_id, fragment_id, earned_at) VALUES
-- Алиса собрала 1 фрагмент
('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'frag-gpu', now() - interval '2 days'),
-- Боб собрал 2 фрагмента
('22222222-2222-2222-2222-222222222222', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'frag-gpu', now()),
('22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'frag-laptop', now());