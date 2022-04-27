begin transaction;

alter table views
    add column last_used_at_tz timestamp with time zone,
    add column created_at_tz timestamp with time zone;

update views
    set last_used_at_tz = last_used_at,
        created_at_tz = created_at;

alter table views drop column last_used_at, drop column created_at;
alter table views rename last_used_at_tz to last_used_at;
alter table views rename created_at_tz to created_at;

commit;