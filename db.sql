CREATE TABLE person (
                        id serial PRIMARY KEY,
                        passwords text[3],
                        current_password text,
                        login varchar(30),
                        content jsonb,
                        created_at timestamp DEFAULT now(),
                        updated_at timestamp DEFAULT now(),
                        deleted_at timestamp DEFAULT now(),
                        cache boolean
);

Create Table people_tracing
(
    id int not null,
    previous_content jsonb,
    date_of_update timestamp default now(),
    created_at timestamp default now(),
    FOREIGN KEY (id) REFERENCES person(id)
);
Create Table Card
(
    id serial primary key,
    person_id int not null,
    card_number int8,
    date_of_expire varchar(5),
    logotype varchar(200),
    money numeric,
    bank_name varchar(100),
    updated_at timestamp default now(),
    deleted_at timestamp default now(),
    created_at timestamp default now(),
    FOREIGN KEY (person_id) references person(id)
);
Create Table Transaction_log
(
    id serial primary key ,
    sender_id int not null ,
    receiver_id int not null,
    time timestamp,
    status varchar(20),
    sum numeric,
    percent float4,
    percentedSum float4,
    FOREIGN KEY(sender_id) references card(id),
    FOREIGN KEY (receiver_id) references card(id)
);
Create Table refresh_token
(
    id serial primary key,
    person_id int not null,
    token text,
    FOREIGN KEY(person_id) references person(id)
);

create or replace function AddPasswordInsert()
    returns trigger as $$
BEGIN
    New.passwords[1] = null;
    New.passwords[2] = null;
    New.passwords[3]  = NEW.current_password;
    return NeW;
end;
$$ LANGUAGE plpgsql;

Create trigger insert_current_password
    before insert on person
    for each row
    execute function AddPasswordInsert();


CREATE TRIGGER update_passwords_trigger
    BEFORE UPDATE OF current_password ON person
    FOR EACH ROW
    WHEN (OLD.current_password IS DISTINCT FROM NEW.current_password)
EXECUTE FUNCTION update_passwords();












Create or replace function auto_time_transaction()
    returns trigger as $$
BEGIN
    new.time = now();
    return new;
end;
    $$ language plpgsql;

Create TRIGGER auto_time_transaction_log
    BEFORE insert on transaction_log
    for each row
    EXECUTE function auto_time_transaction();

CREATE OR REPLACE FUNCTION log_person_update()
    RETURNS TRIGGER AS $$
BEGIN

    INSERT INTO people_tracing (id, previous_content, date_of_update, created_at)
    VALUES (OLD.id, OLD.content, NOW(), NOW());

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER person_update_trigger
    AFTER UPDATE ON person
    FOR EACH ROW
EXECUTE FUNCTION log_person_update();

Create or replace Function update_date_if_update_function()
    Returns TRIGGER AS $$
BEGIN
    NEW.date_of_update = now();
    return new;
END;
    $$LANGUAGE plpgsql;

Create Trigger person_update
    Before INSERT On person
    for each row
    EXECUTE function person_update_data();

Create or replace Function person_update_data()
    Returns TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    return new;
END;
$$LANGUAGE plpgsql;


Create or replace function update_person_updated_at()
    returns trigger as $$
begin
    NEW.updated_at = NOW();
    return NEW;
end;
$$ language plpgsql;
Create trigger person_updated_at_trigger
    before update ON person
    for each row
    Execute function update_person_updated_at();



