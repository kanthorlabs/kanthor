BEGIN;

CREATE TABLE IF NOT EXISTS kanthor_audit (
  table_name VARCHAR(256),
  id VARCHAR(256),
  operation VARCHAR(256),
  operator VARCHAR(256),
  rand VARCHAR(36) DEFAULT gen_random_uuid(),
  PRIMARY KEY (table_name, id, operation, operator, rand),
  created_at TIMESTAMP DEFAULT now(),
  previous jsonb,
  current jsonb
);

CREATE OR REPLACE FUNCTION kanthor_audit_trigger() RETURNS TRIGGER AS $$
DECLARE
    new_data jsonb;
    old_data jsonb;
    key text;
    new_values jsonb;
    old_values jsonb;
    operator VARCHAR(256);
BEGIN
    new_values := '{}';
    old_values := '{}';

    IF TG_OP = 'INSERT' THEN
        new_data := to_jsonb(NEW);
        new_values := new_data;
        operator := NEW.modifier;

    ELSIF TG_OP = 'UPDATE' THEN
        new_data := to_jsonb(NEW);
        old_data := to_jsonb(OLD);
        operator := NEW.modifier;

        FOR key IN SELECT jsonb_object_keys(new_data) INTERSECT SELECT jsonb_object_keys(old_data)
        LOOP
            IF new_data ->> key != old_data ->> key THEN
                new_values := new_values || jsonb_build_object(key, new_data ->> key);
                old_values := old_values || jsonb_build_object(key, old_data ->> key);
            END IF;
        END LOOP;

    ELSIF TG_OP = 'DELETE' THEN
        old_data := to_jsonb(OLD);
        old_values := old_data;
        operator := OLD.modifier;

        FOR key IN SELECT jsonb_object_keys(old_data)
        LOOP
            old_values := old_values || jsonb_build_object(key, old_data ->> key);
        END LOOP;

    END IF;

    IF TG_OP = 'INSERT' OR TG_OP = 'UPDATE' THEN
        INSERT INTO kanthor_audit (table_name, id, operation, operator, created_at, previous, current)
        VALUES (TG_TABLE_NAME, NEW.id, TG_OP, operator, now(), old_values, new_values);

        RETURN NEW;
    ELSE
        INSERT INTO kanthor_audit (table_name, id, operation, operator, created_at, previous, current)
        VALUES (TG_TABLE_NAME, OLD.id, TG_OP, operator, now(), old_values, new_values);

        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER kanthor_audit_trigger_log_workspace BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_workspace
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

CREATE TRIGGER kanthor_audit_trigger_log_application BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_application
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

CREATE TRIGGER kanthor_audit_trigger_log_endpoint BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_endpoint
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

CREATE TRIGGER kanthor_audit_trigger_log_endpoint_rule BEFORE INSERT OR UPDATE OR DELETE 
ON public.kanthor_endpoint_rule
FOR EACH ROW
EXECUTE FUNCTION kanthor_audit_trigger(); 

COMMIT;