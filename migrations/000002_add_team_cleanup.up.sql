CREATE OR REPLACE FUNCTION delete_empty_team()
RETURNS TRIGGER AS $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM users WHERE team_name = OLD.team_name
    ) THEN
        DELETE FROM teams WHERE team_name = OLD.team_name;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_delete_empty_team
AFTER DELETE OR UPDATE OF team_name ON users
FOR EACH ROW
EXECUTE FUNCTION delete_empty_team();
