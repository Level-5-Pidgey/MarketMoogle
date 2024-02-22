create table public.listings
(
    listing_id             integer generated always as identity,
    universalis_listing_id varchar not null,
    item_id                integer not null,
    region_id              integer not null,
    data_center_id         integer not null,
    world_id               integer not null,
    price_per_unit         integer default 0
        constraint listings_price_check
            check (price_per_unit > 0),
    quantity               integer default 1,
    total_price            integer not null,
    is_high_quality        boolean default false,
    retainer_name          varchar(50),
    retainer_city          smallint,
    last_review_time       timestamp,
    constraint listings_total_price_check
        check ((price_per_unit * quantity) <= total_price)
) partition by list (data_center_id);

alter table public.listings add primary key (listing_id, data_center_id);

alter table public.listings add constraint universalis_listing_id_key
    unique (universalis_listing_id, data_center_id);

comment on table public.listings is 'Listings mirrored from Universalis server.';

comment on column public.listings.universalis_listing_id is 'A SHA256 Hash for the Universalis Listing Id.';

alter table public.listings
    owner to admin;

create table public.sales
(
    sales_id             integer generated always as identity,
    item_id              integer not null,
    world_id             integer not null,
    price_per_unit       integer default 0
        constraint sales_price_check
            check (price_per_unit > 0),
    quantity             integer default 1,
    total_price          integer not null,
    is_high_quality      boolean default false,
    buyer_name           varchar(50),
    sale_time            timestamp,
    constraint sales_total_price_check
        check ((price_per_unit * quantity) <= total_price)
) partition by range (sale_time);

alter table public.sales add primary key (sales_id, sale_time);

comment on column public.sales.sale_time is 'Unix timestamp for the sale time.';

alter table public.sales
    owner to admin;

