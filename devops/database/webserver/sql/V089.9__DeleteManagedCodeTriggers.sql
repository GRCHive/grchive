DROP FUNCTION IF EXISTS insert_audit_managed_code_change CASCADE;
DROP FUNCTION IF EXISTS update_audit_managed_code_change CASCADE;
DROP FUNCTION IF EXISTS delete_audit_managed_code_change CASCADE;

DROP TRIGGER IF EXISTS trigger_insert_managed_code ON managed_code CASCADE;
DROP TRIGGER IF EXISTS trigger_update_managed_code ON managed_code CASCADE;
DROP TRIGGER IF EXISTS trigger_delete_managed_code ON managed_code CASCADE;
