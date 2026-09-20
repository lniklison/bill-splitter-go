insert into bill_shares (
    bill_id,
    position,
    person_name,
    normalized_name,
    percentage_basis_points
)
values
    (1, 0, 'Alice', 'alice', 6000),
    (1, 1, 'Bob', 'bob', 4000),
    (2, 0, 'Alex', 'alex', 4500),
    (2, 1, 'Jordan', 'jordan', 3500),
    (2, 2, 'Sam', 'sam', 2000)
on conflict (bill_id, position) do nothing;
