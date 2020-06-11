package pg_api

func (pg *PgDriver) Close() {
	pg.connection.Close()
	pg.db.Close()
}
