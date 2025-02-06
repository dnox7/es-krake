CREATE OR REPLACE TRIGGER before_update_set_updated_at
BEFORE UPDATE ON categories;
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
