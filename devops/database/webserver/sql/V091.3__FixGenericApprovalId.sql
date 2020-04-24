DROP TABLE generic_approval;
CREATE TABLE generic_approval (
    id BIGSERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES generic_requests(id) ON DELETE CASCADE,
    response_time TIMESTAMPTZ NOT NULL,
    responder_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE NO ACTION,
    response BOOLEAN NOT NULL,
    reason TEXT
);
