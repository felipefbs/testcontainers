CREATE TABLE IF NOT EXISTS 
    data_points 
    (id SERIAL PRIMARY KEY,
     inserted_at timestamp DEFAULT CURRENT_TIMESTAMP,
     sended_at timestamp,
     received_at timestamp,
     message varchar(255));
