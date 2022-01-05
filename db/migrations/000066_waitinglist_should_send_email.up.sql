alter table waitinglist
    add column if not exists should_send_email boolean ;