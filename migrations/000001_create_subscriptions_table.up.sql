CREATE TABLE IF NOT EXISTS subscriptions(
    id SERIAL PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date TEXT NOT NULL,
    end_date TEXT,
    UNIQUE (service_name, user_id)
);
