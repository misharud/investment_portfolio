-- +goose Up

CREATE SCHEMA IF NOT EXISTS public;

create table if not exists public.investment_history
(
    id         SERIAL PRIMARY KEY,
    ticker     varchar(50) not null, -- USD / RUB
    quantity   integer     not null,
    price      integer     not null,
    status     smallint    not null        default 0,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    updated_at timestamp without time zone
);

create index if not exists investment_history_ticker_index
    on public.investment_history (ticker);

-- +goose Down

DROP TABLE investment_history;

DROP SCHEMA public;
