CREATE TABLE t_account (
       account_id INTEGER NOT NULL PRIMARY KEY,
       balance VARCHAR(255)
);

CREATE TABLE t_transaction (
     uuid UUID DEFAULT gen_random_uuid() NOT NULL,
     source_account_id INTEGER NOT NULL,
     source_init VARCHAR(255) NOT NULL,
     source_after VARCHAR(255) NOT NULL,
     destination_account_id INTEGER NOT NULL,
     destination_init VARCHAR(255) NOT NULL,
     destination_after VARCHAR(255) NOT NULL,
     amount VARCHAR(255),
     success BOOLEAN DEFAULT FALSE NOT NULL,
     create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     modify_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
     CONSTRAINT fk_source_account FOREIGN KEY (source_account_id) REFERENCES t_account(account_id),
     CONSTRAINT fk_destination_account FOREIGN KEY (destination_account_id) REFERENCES t_account(account_id)
);