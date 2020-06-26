CREATE OR REPLACE FUNCTION get_pbc_request_status(r document_requests)
    RETURNS INT AS
$$
    DECLARE
        OPEN constant INT := 0;
        PROGRESS constant INT := 1;
        FEEDBACK constant INT := 2;
        COMPLETE constant INT := 3;
        OVERDUE constant INT := 4;
    BEGIN
        IF r.completion_time IS NULL THEN
            IF r.due_date IS NOT NULL THEN
                IF NOW() <= r.due_date THEN
                    RETURN OPEN;
                ELSE
                    RETURN OVERDUE;
                END IF;
            ELSIF r.progress_time IS NOT NULL THEN
                RETURN PROGRESS;
            ELSE
                RETURN OPEN;
            END IF;
        ELSIF r.feedback_time IS NOT NULL THEN
            IF r.due_date IS NOT NULL AND NOW() > r.due_date THEN
                RETURN OVERDUE;
            ELSIF r.feedback_time <= r.completion_time THEN
                RETURN COMPLETE;
            ELSE
                RETURN FEEDBACK;
            END IF;
        ELSE
            RETURN COMPLETE;
        END IF;
    END;
$$ LANGUAGE plpgsql;
