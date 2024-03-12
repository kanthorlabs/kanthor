BEGIN;

DROP TRIGGER IF EXISTS kanthor_audit_trigger_log_workspace ON public.kanthor_workspace;

DROP TRIGGER IF EXISTS kanthor_audit_trigger_log_application ON public.kanthor_application;

DROP TRIGGER IF EXISTS kanthor_audit_trigger_log_endpoint ON public.kanthor_endpoint;

DROP TRIGGER IF EXISTS kanthor_audit_trigger_log_endpoint_rule ON public.kanthor_endpoint_rule;

DROP TABLE IF EXISTS kanthor_audit;

DROP FUNCTION IF EXISTS kanthor_audit_trigger;
 
COMMIT;