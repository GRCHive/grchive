package core

type FullDbTable struct {
	DbTable

	Columns map[string]DbColumn
}

func (s *FullDbTable) AddColumn(col *DbColumn) {
	if s.Columns == nil {
		s.Columns = map[string]DbColumn{}
	}
	s.Columns[col.ColumnName] = *col
}

type FullDbSchema struct {
	DbSchema

	Tables    map[string]*FullDbTable
	Functions map[string]DbFunction
}

func (s *FullDbSchema) AddTable(table *DbTable) {
	if s.Tables == nil {
		s.Tables = map[string]*FullDbTable{}
	}

	s.Tables[table.TableName] = &FullDbTable{
		DbTable: *table,
	}
}

func (s *FullDbSchema) AddColumn(table *DbTable, col *DbColumn) {
	tbl := s.Tables[table.TableName]
	tbl.AddColumn(col)
}

func (s *FullDbSchema) AddFunction(fn *DbFunction) {
	if s.Functions == nil {
		s.Functions = map[string]DbFunction{}
	}
	s.Functions[fn.Name] = *fn
}

type FullDbState struct {
	Schemas map[string]*FullDbSchema
}

func (s *FullDbState) AddSchema(schema *DbSchema) {
	if s.Schemas == nil {
		s.Schemas = map[string]*FullDbSchema{}
	}

	s.Schemas[schema.SchemaName] = &FullDbSchema{
		DbSchema: *schema,
	}
}

func (s *FullDbState) AddTable(schema *DbSchema, table *DbTable) {
	sch := s.Schemas[schema.SchemaName]
	sch.AddTable(table)
}

func (s *FullDbState) AddColumn(schema *DbSchema, table *DbTable, col *DbColumn) {
	sch := s.Schemas[schema.SchemaName]
	sch.AddColumn(table, col)
}

func (s *FullDbState) AddFunction(schema *DbSchema, fn *DbFunction) {
	sch := s.Schemas[schema.SchemaName]
	sch.AddFunction(fn)
}

//
// Diffs
//

func (a DbFunction) HasDiff(b DbFunction) bool {
	if a.Src != b.Src {
		return true
	}

	if a.RetType != b.RetType {
		return true
	}

	return false
}

func (a DbColumn) HasDiff(b DbColumn) bool {
	if a.ColumnName != b.ColumnName {
		return true
	}

	if a.ColumnType != b.ColumnType {
		return true
	}

	return false
}

func (a *FullDbTable) HasDiff(b *FullDbTable) bool {
	processedColumns := map[string]bool{}
	for aName, aCol := range a.Columns {
		_, processed := processedColumns[aName]
		if processed {
			continue
		}
		processedColumns[aName] = true

		bCol, ok := b.Columns[aName]
		if !ok {
			return true
		}

		diff := aCol.HasDiff(bCol)
		if diff {
			return true
		}
	}

	for bName, _ := range b.Columns {
		_, processed := processedColumns[bName]
		if processed {
			continue
		}
		processedColumns[bName] = true
		return true
	}
	return false
}

func (a *FullDbSchema) HasDiff(b *FullDbSchema) bool {
	processedFunctions := map[string]bool{}
	for aName, aFn := range a.Functions {
		_, processed := processedFunctions[aName]
		if processed {
			continue
		}
		processedFunctions[aName] = true

		bFn, ok := b.Functions[aName]
		if !ok {
			return true
		}

		diff := aFn.HasDiff(bFn)
		if diff {
			return true
		}
	}

	for bName, _ := range b.Functions {
		_, processed := processedFunctions[bName]
		if processed {
			continue
		}
		processedFunctions[bName] = true
		return true
	}

	processedTables := map[string]bool{}
	for aName, aTbl := range a.Tables {
		_, processed := processedTables[aName]
		if processed {
			continue
		}
		processedTables[aName] = true

		bTbl, ok := b.Tables[aName]
		if !ok {
			return true
		}

		diff := aTbl.HasDiff(bTbl)
		if diff {
			return true
		}
	}

	for bName, _ := range b.Tables {
		_, processed := processedTables[bName]
		if processed {
			continue
		}
		processedTables[bName] = true
		return true
	}
	return false
}

func (a *FullDbState) HasDiff(b *FullDbState) bool {
	processedSchemas := map[string]bool{}

	for aName, aSchema := range a.Schemas {
		_, processed := processedSchemas[aName]
		if processed {
			continue
		}
		processedSchemas[aName] = true

		bSchema, ok := b.Schemas[aName]
		if !ok {
			return true
		}

		diff := bSchema.HasDiff(aSchema)
		if diff {
			return true
		}
	}

	for bName, _ := range b.Schemas {
		_, processed := processedSchemas[bName]
		if processed {
			continue
		}
		processedSchemas[bName] = true

		// At this point, if the schema was in the aSchema
		// then we would've already checked the diff. Therefore
		// these are only unique schemas in b.
		return true
	}

	return false
}
