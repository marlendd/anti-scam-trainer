UPDATE scenario_versions
SET content = jsonb_set(
    jsonb_set(
        content,
        '{nodes,3,choices,2,next_node_id}',
        '"n4_danger"'::jsonb
    ),
    '{nodes,4,choices,2,next_node_id}',
    '"n4_danger"'::jsonb
)
WHERE id = '45d4cc8c-f604-4a7c-b8c5-f2464717b71f';

UPDATE scenario_versions
SET content = jsonb_set(
    jsonb_set(
        jsonb_set(
            jsonb_set(
                jsonb_set(
                    jsonb_set(
                        content,
                        '{nodes,1,choices,0,next_node_id}',
                        '"n3_final_pressure"'::jsonb
                    ),
                    '{nodes,2,choices,0,next_node_id}',
                    '"n3_final_pressure"'::jsonb
                ),
                '{nodes,2,choices,1,next_node_id}',
                '"n3_final_pressure"'::jsonb
            ),
            '{nodes,3,choices,2,next_node_id}',
            '"n4_money_lost"'::jsonb
        ),
        '{nodes,4,choices,2,next_node_id}',
        '"n4_money_lost"'::jsonb
    ),
    '{nodes,5,choices,3,ending_id}',
    '"ending_loss"'::jsonb
)
WHERE id = '854fd38a-3a54-4a9b-a23c-6bf7a0b3b405';

UPDATE scenario_versions
SET content = jsonb_set(
    content,
    '{nodes,5,text}',
    to_jsonb('Продавец переходит на оскорбления и добавляет вас в чёрный список. Через десять минут площадка блокирует его профиль с пометкой «Подозрительная активность».'::text)
)
WHERE id = '854fd38a-3a54-4a9b-a23c-6bf7a0b3b405';

UPDATE scenario_versions
SET content = jsonb_set(
    content,
    '{nodes,3,choices,2,next_node_id}',
    '"n4_danger"'::jsonb
)
WHERE id = 'a8f01498-28bb-42fd-a0dc-b729468f3897';
