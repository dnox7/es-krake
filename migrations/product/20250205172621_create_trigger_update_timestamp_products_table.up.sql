CREATE OR REPLACE TRIGGER updated_at_timestamp_before_update
BEFORE UPDATE ON products;
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();


