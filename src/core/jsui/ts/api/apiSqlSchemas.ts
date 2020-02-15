import axios from 'axios'
import * as qs from 'query-string'
import { 
    allSqlSchemasUrl,
    getSqlSchemaUrl,
} from '../url'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { 
    DbSchema,
    DbTable,
    DbColumn,
    DbFunction
} from '../sql'

export interface TAllSqlSchemasInput {
    refreshId: number
    orgId: number
}

export interface TAllSqlSchemasOutput {
    data: DbSchema[]
}

export function allSqlSchemas(inp : TAllSqlSchemasInput) : Promise<TAllSqlSchemasOutput> {
    return axios.get(allSqlSchemasUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetSqlSchemaInput {
    schemaId: number
    orgId: number
    fnMode: boolean
}

export interface TGetSqlSchemaOutput {
    data: {
        Schema: {
            Tables: DbTable[]
            Columns: Record<number, DbColumn[]>
        } | null
        
        Functions: DbFunction[] | null
    }
}

export function getSqlSchema(inp : TGetSqlSchemaInput) : Promise<TGetSqlSchemaOutput> {
    return axios.get(getSqlSchemaUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
