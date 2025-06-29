CREATE TABLE role (
    id serial PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE user_account (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(100) UNIQUE, 
    email VARCHAR(255) UNIQUE NOT NULL,
    name TEXT,
    password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role_id INTEGER NOT NULL,
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES role(id) ON DELETE RESTRICT
);

CREATE TABLE packages (
    id SERIAL PRIMARY KEY,               
    name TEXT NOT NULL,                  
    price INTEGER NOT NULL,              
    description TEXT,                    
    benefits TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,  
    user_id INTEGER NOT NULL REFERENCES user_account(id) ON DELETE CASCADE,
    package_id INTEGER NOT NULL REFERENCES packages(id) ON DELETE RESTRICT,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (status IN ('PENDING', 'ACTIVE', 'EXPIRED', 'CANCELLED')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    subscription_id INTEGER NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    reference_id VARCHAR(100) UNIQUE NOT NULL,        
    xendit_payment_id VARCHAR(100) NOT NULL,         
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('QR_CODE', 'VIRTUAL_ACCOUNT')), 
    amount INTEGER NOT NULL,                         
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING', 'PAID', 'EXPIRED', 'FAILED')),
    paid_at TIMESTAMP,                               
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);