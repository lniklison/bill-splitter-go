insert into bills (id, description, total_cents)
values
    (1, 'Team dinner', 12000),
    (2, 'Weekend groceries', 7345)
on conflict (id) do nothing;
