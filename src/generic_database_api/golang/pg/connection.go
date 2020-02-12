package pg_api

func (pg *PgDriver) ConnectionReadOnly() bool {
	return false
}
