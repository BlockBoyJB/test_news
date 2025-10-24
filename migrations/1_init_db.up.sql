create table if not exists news
(
    id      bigserial primary key,
    title   varchar not null,
    content text    not null
);

create table if not exists news_categories
(
    news_id     bigint references news (id) not null,
    category_id bigint                      not null,
    unique (news_id, category_id)
);

create index if not exists idx_news_categories_news_id
    on news_categories (news_id);

