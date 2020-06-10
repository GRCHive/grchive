package sap_api

func (c sapClient) RunRfc(fn string) (map[string]interface{}, error) {
	_, err := c.conn.GetFunctionDescription(fn)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{}
	return c.conn.Call(fn, params)
}
