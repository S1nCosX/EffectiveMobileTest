CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
	service_name TEXT NOT NULL,
    price INTEGER NOT NULL,
    start_date DATE NOT NULL ,
    end_date DATE
);
