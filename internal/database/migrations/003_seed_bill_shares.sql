insert into bill_shares (
    bill_id,
    position,
    person_name,
    normalized_name,
    percentage_basis_points
)
values
    (1, 0, 'Alice', 'alice', 6000),
    (1, 1, 'Bob', 'bob', 4000)
on conflict (bill_id, position) do nothing;
