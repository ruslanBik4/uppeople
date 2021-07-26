CREATE OR REPLACE FUNCTION after_insert_candidates()
    RETURNS trigger
    LANGUAGE plpgsql
AS
$$

BEGIN
  insert into logs (candidate_id,  changed, action_code)
    VALUES (new.id, json_build_object('comment', new.comments), (select id from log_actions where name = 'CODE_LOG_INSERT'));

  return NEW;
END;

$$;

create trigger insert_additives_trg
    after insert
    on candidates
    for each row
execute function after_insert_candidates();