alter table waitinglist
    add column invited_at timestamp,
    add column ignored boolean;