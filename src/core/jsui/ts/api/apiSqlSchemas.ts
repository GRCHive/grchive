import axios from 'axios'
import * as qs from 'query-string'
import { 
    allSqlRefreshUrl,
    allSqlSchemasUrl,
    getSqlSchemaUrl,
} from '../url'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { 
    DbRefresh,
    cleanDbRefreshFromJson,
    DbSchema,
    DbTable,
    DbColumn
} from '../sql'

export interface TAllSqlRefreshInput {
    dbId: number
    orgId: number
}

export interface TAllSqlRefreshOutput {
    data: DbRefresh[]
}

export function allSqlRefresh(inp : TAllSqlRefreshInput) : Promise<TAllSqlRefreshOutput> {
    return axios.get(allSqlRefreshUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then(
        (resp : TAllSqlRefreshOutput) => {
            resp.data.forEach((ele : DbRefresh) => { cleanDbRefreshFromJson(ele) })
            return resp
        })
}

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
}

export interface TGetSqlSchemaOutput {
    data: {
        Tables: DbTable[]
        Columns: Record<number, DbColumn[]>
    }
}

export function getSqlSchema(inp : TGetSqlSchemaInput) : Promise<TGetSqlSchemaOutput> {
    return axios.get(getSqlSchemaUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
