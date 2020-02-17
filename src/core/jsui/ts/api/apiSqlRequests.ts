import axios from 'axios'
import * as qs from 'query-string'
import { 
    newSqlRequestUrl,
    allSqlRequestUrl,
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { 
    DbSqlQueryRequest,
    cleanDbSqlRequestFromJson
} from '../sql'

export interface TNewSqlRequestInput {
    queryId: number
    orgId : number
    name: string
    description: string
}

export interface TNewSqlRequestOutput {
    data: DbSqlQueryRequest
}

export function newSqlRequest(inp : TNewSqlRequestInput) : Promise<TNewSqlRequestOutput> {
    return postFormJson<TNewSqlRequestOutput>(newSqlRequestUrl, inp, getAPIRequestConfig()).then((resp : TNewSqlRequestOutput) => {
        cleanDbSqlRequestFromJson(resp.data)
        return resp
    })
}

export interface TAllSqlRequestInput {
    dbId? : number
    orgId : number
}

export interface TAllSqlRequestOutput {
    data: DbSqlQueryRequest[]
}

export function allSqlRequest(inp : TAllSqlRequestInput) : Promise<TAllSqlRequestOutput> {
    return axios.get(allSqlRequestUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllSqlRequestOutput) => {
        resp.data.forEach(cleanDbSqlRequestFromJson)
        return resp
    })
}
