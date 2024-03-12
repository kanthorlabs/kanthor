BEGIN;

CREATE TABLE IF NOT EXISTS kanthor_audit (
  table_name VARCHAR(256),
  id VARCHAR(256),
  operation VARCHAR(256),
  operator VARCHAR(256),
  rand VARCHAR(36) DEFAULT gen_random_uuid(),
  PRIMARY KEY (table_name, id, operation, operator, rand),
  created_at BIGINT	DEFAULT ROUND(EXTRACT(EPOCH FROM NOW()) * 1000),
  previous jsonb,
  current jsonb
);

CREATE OR REPLACE FUNCTION kanthor_audit_trigger() RETURNS TRIGGER AS $$
DECLARE
    new_data jsonb;
    old_data jsonb;
    diffkey text;
    new_values jsonb;
    old_values jsonb;
    operator VARCHAR(256);
    ts BIGINT;
BEGIN
    new_values := '{}';
    old_values := '{}';
    ts := ROUND(EXTRACT(EPOCH FROM NOW()) * 1000);

    IF TG_OP = 'INSERT' THEN
        new_data := to_jsonb(NEW);
        new_values := new_data;
        operator := NEW.modifier;

    ELSIF TG_OP = 'UPDATE' THEN
        new_data := to_jsonb(NEW);
        old_data := to_jsonb(OLD);
        operator := NEW.modifier;

        FOR diffkey IN SELECT jsonb_object_keys(new_data) INTERSECT SELECT jsonb_object_keys(old_data)
        LOOP
            IF new_data ->> diffkey != old_data ->> diffkey THEN
                new_values := new_values || jsonb_build_object(diffkey, new_data ->> diffkey);
                old_values := old_values || jsonb_build_object(diffkey, old_data ->> diffkey);
            END IF;
        END LOOP;

    ELSIF TG_OP = 'DELETE' THEN
        old_data := to_jsonb(OLD);
        old_values := old_data;
        operator := OLD.modifier;
        ts := ts + 1;

        FOR diffkey IN SELECT jsonb_object_keys(old_data)
        LOOP
            old_values := old_values || jsonb_build_object(diffkey, old_data ->> diffkey);
        END LOOP;

    END IF;

    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        INSERT INTO kanthor_audit (table_name, id, operation, operator, created_at, previous, current)
        VALUES (TG_TABLE_NAME, NEW.id, TG_OP, operator, ts, old_values, new_values);

        RETURN NEW;
    ELSE
        INSERT INTO kanthor_audit (table_name, id, operation, operator, created_at, previous, current)
        VALUES (TG_TABLE_NAME, OLD.id, TG_OP, operator, ts, old_values, new_values);

        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER kanthor_audit_trigger_log_workspace BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_workspace
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

CREATE OR REPLACE TRIGGER kanthor_audit_trigger_log_application BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_application
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

CREATE OR REPLACE TRIGGER kanthor_audit_trigger_log_endpoint BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_endpoint
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

CREATE OR REPLACE TRIGGER kanthor_audit_trigger_log_endpoint_rule BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_endpoint_rule
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

COMMIT;