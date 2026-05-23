-- English Tutor — initial schema.
-- Curriculum hierarchy: level > lesson > topic > exercise.
-- Quizzes group exercises independently of topics. Attempts record learner answers.

create table levels (
    id          bigint generated always as identity primary key,
    code        text not null unique,
    name        text not null,
    description text not null default '',
    position    int  not null default 0
);

create table lessons (
    id          bigint generated always as identity primary key,
    level_id    bigint not null references levels (id) on delete cascade,
    number      int    not null,
    title       text   not null,
    summary     text   not null default '',
    position    int    not null default 0,
    unique (level_id, number)
);

create table topics (
    id          bigint generated always as identity primary key,
    lesson_id   bigint not null references lessons (id) on delete cascade,
    title       text   not null,
    explanation text   not null default '',
    position    int    not null default 0
);

create table quizzes (
    id          bigint generated always as identity primary key,
    level_id    bigint references levels (id) on delete set null,
    title       text   not null,
    description text   not null default '',
    position    int    not null default 0
);

-- A single practice-item table. An exercise belongs to exactly one parent:
-- a topic (lesson practice) or a quiz (assessment).
create table exercises (
    id          bigint generated always as identity primary key,
    topic_id    bigint references topics (id) on delete cascade,
    quiz_id     bigint references quizzes (id) on delete cascade,
    kind        text   not null check (kind in ('mcq', 'fill_blank', 'true_false')),
    prompt      text   not null,
    choices     jsonb  not null default '[]'::jsonb,
    answer      text   not null,
    explanation text   not null default '',
    position    int    not null default 0,
    constraint exercises_single_parent check ((topic_id is null) <> (quiz_id is null))
);

create table attempts (
    id           bigint generated always as identity primary key,
    exercise_id  bigint not null references exercises (id) on delete cascade,
    given_answer text   not null,
    is_correct   boolean not null,
    created_at   timestamptz not null default now()
);

create table vocabulary (
    id          bigint generated always as identity primary key,
    level_id    bigint references levels (id) on delete set null,
    lesson_id   bigint references lessons (id) on delete set null,
    category    text   not null default '',
    term        text   not null,
    definition  text   not null default '',
    example     text   not null default '',
    position    int    not null default 0
);

create index lessons_level_id_idx     on lessons (level_id);
create index topics_lesson_id_idx     on topics (lesson_id);
create index exercises_topic_id_idx   on exercises (topic_id);
create index exercises_quiz_id_idx    on exercises (quiz_id);
create index attempts_exercise_id_idx on attempts (exercise_id);
create index attempts_created_at_idx  on attempts (created_at);
create index vocabulary_level_id_idx  on vocabulary (level_id);
create index vocabulary_lesson_id_idx on vocabulary (lesson_id);
