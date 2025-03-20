package database

import (
	"context"
	"log"
)

func (repo *PostgresRepo) Init() error {
	dbConn, err := repo.pool.Acquire(context.Background())
	if err != nil {
		log.Fatalln("error occurred while aquiring db connection")
	}

	defer dbConn.Release()

	tx, err := dbConn.Begin(context.Background())

	if err != nil {
		return err
	}

	dbInitQueries := []string{
		`DO $$ 
		 BEGIN
    		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'employee_leave_status') THEN
        		CREATE TYPE employee_leave_status AS ENUM ('pending', 'granted', 'canceled');
    		END IF;
		END $$;
		`,
		`DO $$ 
		 BEGIN
    		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_type') THEN
        		CREATE TYPE user_type AS ENUM ('admin', 'user');
    		END IF;
		END $$;
		`,
		`CREATE OR REPLACE FUNCTION update_timestamp()
			RETURNS TRIGGER AS $$
		BEGIN
    		NEW.updated_at = NOW();
    		RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;`,
		`CREATE TABLE IF NOT EXISTS admins (
			admin_id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			dob VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			phone_number VARCHAR(255) NOT NULL,
			profile_url VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			position VARCHAR(255) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS employee_category(
			category_id VARCHAR(255) PRIMARY KEY,
			admin_id VARCHAR(255) NOT NULL,
			category_name VARCHAR(255) NOT NULL,
			category_description TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			FOREIGN KEY (admin_id) REFERENCES admins(admin_id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			user_id VARCHAR(255) PRIMARY KEY,
			admin_id VARCHAR(255) NOT NULL,
			category_id VARCHAR(255) NOT NULL,
			name VARCHAR(255) NOT NULL,
			dob VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL UNIQUE,
			phone_number VARCHAR(255) NOT NULL,
			profile_url VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			position VARCHAR(255) NOT NULL,
			work_date VARCHAR(255) DEFAULT '0000-00-00',
			login_time VARCHAR(255) DEFAULT '00:00',
			logout_time VARCHAR(255) DEFAULT '00:00',
			login_status BOOLEAN DEFAULT false,
			latitude VARCHAR(255) DEFAULT '0.0',
			longitude VARCHAR(255) DEFAULT '0.0',
			uploaded_work TEXT DEFAULT 'welcome to the ca app',
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			FOREIGN KEY (admin_id) REFERENCES admins(admin_id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES employee_category(category_id) ON DELETE CASCADE
		);`,
		`CREATE TABLE IF NOT EXISTS users_history (
			user_id VARCHAR(255) NOT NULL,
			work_date VARCHAR(255) NOT NULL,
			login_time VARCHAR(255) NOT NULL,
			logout_time VARCHAR(255) NOT NULL,
			latitude VARCHAR(255) DEFAULT '0.0',
			longitude VARCHAR(255) DEFAULT '0.0',
			uploaded_work TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE 
		);`,
		`CREATE TABLE IF NOT EXISTS users_leave_history (
			leave_id VARCHAR(255) PRIMARY KEY,
			user_id VARCHAR(255) NOT NULL,
			leave_from VARCHAR(255) NOT NULL,
			leave_to VARCHAR(255) NOT NULL,
			leave_reason TEXT NOT NULL,
			status employee_leave_status NOT NULl,
			status_updated_by user_type NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE 
		);`,

		`CREATE TABLE IF NOT EXISTS user_otps (
			email VARCHAR(255) NOT NULL,
			otp VARCHAR(255) NOT NULL,
			expire_time TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS admin_otps (
			email VARCHAR(255) NOT NULL,
			otp VARCHAR(255) NOT NULL,
			expire_time TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			FOREIGN KEY (email) REFERENCES admins (email) ON DELETE CASCADE
		)`,
		`DO $$ 
		 BEGIN
    		IF NOT EXISTS (
        		SELECT 1 FROM pg_trigger 
        		WHERE tgname = 'set_timestamp_admins'
    		) THEN
        		CREATE TRIGGER set_timestamp_admins
        		BEFORE UPDATE ON admins
        		FOR EACH ROW
        		EXECUTE FUNCTION update_timestamp();
    		END IF;
		END $$;`,

		`DO $$ 
		 BEGIN
    		IF NOT EXISTS (
        		SELECT 1 FROM pg_trigger 
        		WHERE tgname = 'set_timestamp_users'
    		) THEN
        		CREATE TRIGGER set_timestamp_users
        		BEFORE UPDATE ON users
        		FOR EACH ROW
        		EXECUTE FUNCTION update_timestamp();
    		END IF;
		END $$;`,
		`DO $$ 
		 BEGIN
    		IF NOT EXISTS (
        		SELECT 1 FROM pg_trigger 
        		WHERE tgname = 'set_timestamp_users_leave_history'
    		) THEN
        		CREATE TRIGGER set_timestamp_users_leave_history
        		BEFORE UPDATE ON users_leave_history
        		FOR EACH ROW
        		EXECUTE FUNCTION update_timestamp();
    		END IF;
		END $$;`,
	}

	for index, query := range dbInitQueries {
		if _, err := tx.Exec(context.Background(), query); err != nil {
			tx.Rollback(context.Background())
			log.Println(index + 1)
			return err
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}
