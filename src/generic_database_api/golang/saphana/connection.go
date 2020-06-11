package saphana_api

func (hdb *SapHanaDriver) Close() {
	hdb.connection.Close()
	hdb.db.Close()
}
