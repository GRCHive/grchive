CREATE OR REPLACE FUNCTION user_to_human_name(u users)
    RETURNS TEXT AS
$$
    BEGIN
        RETURN u.first_name || ' ' || u.last_name || ' (' || u.email || ')';
    END;
$$ LANGUAGE plpgsql;
