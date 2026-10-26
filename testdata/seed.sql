INSERT INTO aws_cloud.account(
    account_id,
    account_name,
    account_arn,
    email,
    account_state,
    joined_method,
    joined_at
)
VALUES 
(
    '111111111111',
    'test-account-1',
    'arn:aws:organizations::111111111111:account/o-abc123/111111111111',
    'foo@bar.com',
    'ACTIVE',
    'CREATED',
    '2026-06-23 14:30:00+00'
),
(
    '222222222222',
    'test-account-2',
    'arn:aws:organizations::111111111111:account/o-abc123/222222222222',
    'bar@baz.com',
    'ACTIVE',
    'INVITED',
    '2026-06-23 14:30:00+00'
);
