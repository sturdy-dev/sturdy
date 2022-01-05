CREATE TABLE completed_onboarding_steps (
    user_id text NOT NULL,
    step_id text NOT NULL,
    created_at timestamp WITH TIME ZONE NOT NULL
);

CREATE UNIQUE INDEX completed_onboarding_steps_unique_idx ON completed_onboarding_steps(user_id, step_id);
