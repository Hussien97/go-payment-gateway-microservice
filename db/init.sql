-- Check if the transactions table exists and being used by docker to auto create the table
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'transactions') THEN
        CREATE TABLE transactions (
            id SERIAL PRIMARY KEY,
            transaction_id VARCHAR(255) NOT NULL UNIQUE,
            amount DECIMAL(10, 2) NOT NULL,
            type VARCHAR(50) NOT NULL,
            status VARCHAR(50) NOT NULL,
            data_format VARCHAR(50) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    END IF;
END $$;