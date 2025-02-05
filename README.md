Task Manager Backend

1. Postgres Database and Login:
● Create a PostgresSQL database to store user information such as name
and email.
● Develop a login APIs that allows users to sign up and assigns them an
internal login ID.
● Develop APIs to create tasks. Any user can create a task. (You can
assume a schema for the task.). Taks should have status.
● Implement a bonus feature That lets other users provide approval for
the task

2. Multi approval Process:
● Develop APIs for multi-approval processes where a task needs to be
approved by 3 other users to be in approved status.
● Provide design pointers or hooks for
● User creating task to choose the 3 other users from a dropdown
list and send email notifications to each user when a new task is
created.
● Add the functionality for users to add comments.
● Ensure the process creator receives a notification on their page
when anyone signs off, and notify all parties involved via email
when everyone signs off.

3. API:
● Break down the multi-approval process into REST APIs.
● Ensure that the APIs can be integrated into any webpage.

--------------------------------------------------------------------------------------------------------------------------------------
Implemented Design:
1. We will be making use of 4 tables, schemas as follows:

CREATE TABLE users (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password        VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*USERS* table is responsible for storing user details with primary key as ID (which is also loginID)


CREATE TABLE tasks (
    id              SERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT NOT NULL,
    task_status     VARCHAR(50) NOT NULL,
    creator_id      INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    approvers       INT[],
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*TASKS* table is responsible for storing task details with primary key as ID (FK on creator_id)

CREATE TABLE task_approvals (
    id              SERIAL PRIMARY KEY,
    task_id         INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    approver_id     INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    approved_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (task_id, approver_id)       --------> Ensures that only a single approval can come from each user_id for a particular task
);
*TASK_APPROVALS* table is responsible for storing the approvals with primary key as ID (KS on task_id, approver_id)

CREATE TABLE task_comments (
    id              SERIAL PRIMARY KEY,
    task_id         INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    commentor       INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    comment         TEXT NOT NULL,
    commented_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
*TASK_COMMENTS* table is responsible for storing the comments made on a single task (FK on task_id, commentor)

We will be making use of 7 APIs:
We will be making use of JWT token with a expiry time of 30 minutes (Configurable)
*LOGIN AND TOKEN GENERATION API*
SIGNUP API 	    = "/api/user/signup"
LOGIN API  	    = "/api/user/login"

*JWT PROTECTED APIs*
GETUSERS API	= "/api/user/get"
CREATETASK API	= "/api/task/create"
APPROVETASK API = "/api/task/approve"
COMMENTTASK API = "/api/task/comment"
GETTASK API 	= "/api/task/get/{id}"

Testing of the flow is achieved via APIs in a sequential manner. We are making use of BRUNO API manager to test the flow.
The service is run using a docker.
*RUN*: docker-compose up

Repo shared with: satish-xalts, reuvab, vinay-xalts, vineet-xalts (Shared the repo link on the mail thread)