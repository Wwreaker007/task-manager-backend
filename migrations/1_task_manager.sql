-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE tasks (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT NOT NULL,
    task_status     VARCHAR(50) NOT NULL,
    creator_id      INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    approvers       INT[],
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE task_approvals (
    id              SERIAL PRIMARY KEY,
    task_id         INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    approver_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    approved_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (task_id, approver_id)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE task_comments (
    id              SERIAL PRIMARY KEY,
    task_id         INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    commentor       INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment         TEXT NOT NULL,
    commented_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX idx_task_ids_on_task_approvals ON task_approvals(task_id);
CREATE INDEX idx_task_ids_on_task_comments ON task_comments(task_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks CASCADE;
DROP TABLE users CASCADE;
DROP TABLE task_approvals CASCADE;
DROP TABLE task_comments CASCADE;
-- +goose StatementEnd