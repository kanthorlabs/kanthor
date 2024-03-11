BEGIN;

DROP TRIGGER IF EXISTS kanthor_audit_trigger_log_workspace ON public.kanthor_workspace;

DROP FUNCTION IF EXISTS kanthor_audit_trigger;

DROP TABLE IF EXISTS kanthor_audit;
 
COMMIT;