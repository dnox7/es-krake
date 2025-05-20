CREATE OR REPLACE TRIGGER update_timestamp_berfore_update
BEFORE UPDATE ON attributes
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
