import axios from 'axios'
import * as qs from 'query-string'
import { 
    allSqlQueryUrl,
    getSqlQueryUrl,
    newSqlQueryUrl,
    updateSqlQueryUrl,
    deleteSqlQueryUrl,
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { 
    DbSqlQueryMetadata,
    DbSqlQuery,
    cleanDbSqlQueryFromJson
} from '../sql'

export interface TAllSqlQueryInput {
    dbId: number
    orgId: number
}

export interface TAllSqlQueryOutput {
    data: DbSqlQueryMetadata[]
}

export function allSqlQuery(inp : TAllSqlQueryInput) : Promise<TAllSqlQueryOutput> {
    return axios.get(allSqlQueryUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetSqlQueryInput {
    metadataId: number
    orgId: number
}

export interface TGetSqlQueryOutput {
    data: DbSqlQuery[]
}

export function getSqlQuery(inp : TGetSqlQueryInput) : Promise<TGetSqlQueryOutput> {
    return axios.get(getSqlQueryUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetSqlQueryOutput) => {
        resp.data.forEach(cleanDbSqlQueryFromJson)
        return resp
    })
}

export interface TNewSqlQueryInput {
    dbId: number
    orgId: number
    name: string
    description: string
    uploadUserId: number
    query : string
}

export interface TNewSqlQueryOutput {
    data: {
        Metadata: DbSqlQueryMetadata
        Query: DbSqlQuery
    }
}

export function newSqlQuery(inp : TNewSqlQueryInput) : Promise<TNewSqlQueryOutput> {
    return postFormJson<TNewSqlQueryOutput>(newSqlQueryUrl, inp, getAPIRequestConfig()).then((resp : TNewSqlQueryOutput) => {
        cleanDbSqlQueryFromJson(resp.data.Query)
        return resp
    })
}

export interface TUpdateSqlQueryInput {
    orgId: number
    metadataId: number
    metadata?: {
        name: string
        description: string
    }

    query? : {
        query: string
        uploadUserId: number
    }
}

export interface TUpdateSqlQueryOutput {
    data: {
        Metadata: DbSqlQueryMetadata | null
        Query: DbSqlQuery | null
    }
}

export function updateSqlQuery(inp : TUpdateSqlQueryInput) : Promise<TUpdateSqlQueryOutput> {
    return postFormJson<TUpdateSqlQueryOutput>(updateSqlQueryUrl, inp, getAPIRequestConfig()).then((resp : TUpdateSqlQueryOutput) => {
        if (!!resp.data.Query) {
            cleanDbSqlQueryFromJson(resp.data.Query)
        }
        return resp
    })
}

export interface TDeleteSqlQueryInput {
    orgId: number
    metadataId: number
}

export function deleteSqlQuery(inp : TDeleteSqlQueryInput) : Promise<void> {
    return postFormJson<void>(deleteSqlQueryUrl, inp, getAPIRequestConfig())
}
