create table bills (
    id bigint primary key,
    description text not null,
    total_cents bigint not null check (total_cents > 0)
);

create table bill_shares (
    bill_id bigint not null references bills(id) on delete cascade,
    position integer not null check (position >= 0),
    person_name text not null,
    normalized_name text not null,
    percentage_basis_points integer not null check (
        percentage_basis_points > 0 and percentage_basis_points <= 10000
    ),
    primary key (bill_id, position),
    unique (bill_id, normalized_name)
);
