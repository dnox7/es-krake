CREATE ROLE esk_dev_migrator WITH LOGIN PASSWORD 'initial_pw';
GRANT CONNECT ON DATABASE esk_dev_1 TO esk_dev_migrator;
GRANT USAGE, CREATE ON SCHEMA public TO esk_dev_migrator;
