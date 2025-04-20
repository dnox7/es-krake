CREATE OR REPLACE TRIGGER updated_at_timestamp_before_update
BEFORE UPDATE ON category_parent
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


