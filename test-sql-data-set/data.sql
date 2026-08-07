INSERT INTO users (
	email,
	password_hash
) VALUES (
	'piter@gmail.com',
	'fsdf98789s0-sdf789798s0df'
);

INSERT INTO scenario_versions(
	logical_id,
	version,
	role,
	title,
	description,
	is_active,
	content
) VALUES (
	gen_random_uuid(),
	1,
	'buyer',
	'buying videocard',                   
	'scam when buying',                     
	true,
	'{
  "start_node_id": "node-start",
  "nodes": [
    {
      "id": "node-start",
      "author": "seller",
      "text": "Сообщение продавца",
      "choices": [
        {
          "id": "choice-1",
          "text": "Ответ пользователя",
          "consequence": "Последствие выбора",
          "explanation": "Почему это безопасно или опасно",
          "weight": 3,
          "score": 100,
          "risk_categories": ["СМС оповещение", "Стороняя ссылка"],
          "next_node_id": "node-2",
          "ending_id": ""
        }
      ]
    }
  ],
  "endings": [
    {
      "id": "ending-safe",
      "header": "Безопасная концовка",
      "result": "Пользователь избежал мошенничества"
    }
  ]
}'::jsonb  
);

INSERT INTO attempts (
	user_id,
	scenario_id,
	status,
	current_node_id
) VALUES (
	(SELECT id FROM users ORDER BY created_at DESC LIMIT 1),
	(SELECT id FROM scenario_versions ORDER BY created_at DESC LIMIT 1),
	'in_progress',
	'start_node'
);

INSERT INTO answers (
	attempt_id,
	node_id,
	choice_id,
	idempotency_key,
	weight,
	choice_score,
	risk_categories,
	consequence,
	explanation,
	response
) VALUES (
	(SELECT id FROM attempts ORDER BY started_at DESC LIMIT 1), 
	'start_dialog',
	'get_suspicious link',
	gen_random_uuid(), 
	2,
	50,
	'["Not check url"]'::jsonb,
	'lost a prime',
	'You shouldnt click on suspicious links.',
	'{"emotion": "happy"}'::jsonb
);

INSERT INTO answers (
	attempt_id,
	node_id,
	choice_id,
	idempotency_key,
	weight,
	choice_score,
	risk_categories,
	consequence,
	explanation,
	response
) VALUES (
	(SELECT id FROM attempts ORDER BY started_at DESC LIMIT 1), 
	'check_url_node',           
	'ignore_link',              
	gen_random_uuid(),          
	1,                          
	100,                        
	'[]'::jsonb,                
	'safely avoided the scam',
	'User successfully recognized a suspicious URL.',
	'{"emotion": "relieved"}'::jsonb
);

SELECT * FROM answers;