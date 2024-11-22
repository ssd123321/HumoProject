create table person
(
    id         serial
        primary key,
    content    jsonb,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    deleted_at timestamp default now(),
    cache      boolean
);
Create Table people_tracing
(
    id int not null,
    previous_content jsonb,
    date_of_update timestamp default now(),
    created_at timestamp default now(),
    FOREIGN KEY (id) REFERENCES person(id)
);
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


