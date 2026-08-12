UPDATE scenario_versions
SET content = jsonb_set(
    jsonb_set(
        content,
        '{nodes,3,choices,2,next_node_id}',
        '"n4_protected"'::jsonb
    ),
    '{nodes,4,choices,2,next_node_id}',
    '"n4_protected"'::jsonb
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
                        '"n4_refusal_result"'::jsonb
                    ),
                    '{nodes,2,choices,0,next_node_id}',
                    '"n4_refusal_result"'::jsonb
                ),
                '{nodes,2,choices,1,next_node_id}',
                '"n4_refusal_result"'::jsonb
            ),
            '{nodes,3,choices,2,next_node_id}',
            '"n4_refusal_result"'::jsonb
        ),
        '{nodes,4,choices,2,next_node_id}',
        '"n4_refusal_result"'::jsonb
    ),
    '{nodes,5,choices,3,ending_id}',
    '"ending_uncertain"'::jsonb
)
WHERE id = '854fd38a-3a54-4a9b-a23c-6bf7a0b3b405';

UPDATE scenario_versions
SET content = jsonb_set(
    content,
    '{nodes,5,text}',
    to_jsonb('Переписка прекращается. Через десять минут площадка блокирует профиль продавца с пометкой «Подозрительная активность».'::text)
)
WHERE id = '854fd38a-3a54-4a9b-a23c-6bf7a0b3b405';

UPDATE scenario_versions
SET content = jsonb_set(
    content,
    '{nodes,3,choices,2,next_node_id}',
    '"n4_protected"'::jsonb
)
WHERE id = 'a8f01498-28bb-42fd-a0dc-b729468f3897';
