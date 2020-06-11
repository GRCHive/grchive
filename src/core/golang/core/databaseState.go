package core

type FullDbSchema struct {
	DbSchema

	Tables    map[string]*DbTable
	Functions map[string]*DbFunction
}

func (s *FullDbSchema) AllTables() []*DbTable {
	ret := make([]*DbTable, 0)
	for _, tbl := range s.Tables {
		ret = append(ret, tbl)
	}
	return ret
}

func (s *FullDbSchema) AllFunctions() []*DbFunction {
	ret := make([]*DbFunction, 0)
	for _, fn := range s.Functions {
		ret = append(ret, fn)
	}
	return ret
}

func (s *FullDbSchema) GetTable(nm string) *DbTable {
	return s.Tables[nm]
}

func (s *FullDbSchema) AddTable(table *DbTable) {
	if s.Tables == nil {
		s.Tables = map[string]*DbTable{}
	}

	s.Tables[table.TableName] = table
}

func (s *FullDbSchema) AddFunction(fn *DbFunction) {
	if s.Functions == nil {
		s.Functions = map[string]*DbFunction{}
	}
	s.Functions[fn.Name] = fn
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

func (s *FullDbState) AllSchemas() []*FullDbSchema {
	ret := make([]*FullDbSchema, 0)
	for _, sch := range s.Schemas {
		ret = append(ret, sch)
	}
	return ret
}

func (s *FullDbState) AllTables() []*DbTable {
	schemas := s.AllSchemas()
	ret := make([]*DbTable, 0)
	for _, sch := range schemas {
		ret = append(ret, sch.AllTables()...)
	}
	return ret
}

func (s *FullDbState) GetSchema(nm string) *FullDbSchema {
	return s.Schemas[nm]
}

func (s *FullDbState) GetTable(schema string, nm string) *DbTable {
	return s.GetSchema(schema).GetTable(nm)
}

func (s *FullDbState) AddTable(schema *DbSchema, table *DbTable) {
	sch := s.Schemas[schema.SchemaName]
	sch.AddTable(table)
}

func (s *FullDbState) AddFunction(schema *DbSchema, fn *DbFunction) {
	sch := s.Schemas[schema.SchemaName]
	sch.AddFunction(fn)
}

//
// Diffs
//

func (a *RawDbColumn) HasDiff(b *RawDbColumn) bool {
	return a.Type != b.Type
}

func (a *DbFunction) HasDiff(b *DbFunction) bool {
	if a.Src != b.Src {
		return true
	}

	if a.RetType != b.RetType {
		return true
	}

	return false
}

func (a *DbTable) HasDiff(b *DbTable) bool {
	aColumns := map[string]*RawDbColumn{}
	for _, aCol := range a.Columns {
		aColumns[aCol.Name] = aCol
	}

	bColumns := map[string]*RawDbColumn{}
	for _, bCol := range b.Columns {
		bColumns[bCol.Name] = bCol
	}

	processedColumns := map[string]bool{}
	for aKey, aCol := range aColumns {
		_, processed := processedColumns[aKey]
		if processed {
			continue
		}
		processedColumns[aKey] = true

		bCol, ok := bColumns[aKey]
		if !ok {
			return true
		}

		diff := aCol.HasDiff(bCol)
		if diff {
			return true
		}
	}

	for bKey, _ := range bColumns {
		_, processed := processedColumns[bKey]
		if processed {
			continue
		}
		processedColumns[bKey] = true
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
