import axios from 'axios'
import * as qs from 'query-string'
import { 
    allSqlRefreshUrl,
    newSqlRefreshUrl,
    getSqlRefreshUrl,
} from '../url'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { 
    DbRefresh,
    cleanDbRefreshFromJson,
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

export interface TNewSqlRefreshInput {
    dbId: number
    orgId: number
}

export interface TNewSqlRefreshOutput {
    data: DbRefresh
}

export function newSqlRefresh(inp : TNewSqlRefreshInput) : Promise<TNewSqlRefreshOutput> {
    return postFormJson<TNewSqlRefreshOutput>(newSqlRefreshUrl, inp, getAPIRequestConfig()).then(
        (resp : TNewSqlRefreshOutput) => {
            cleanDbRefreshFromJson(resp.data)
            return resp
        })
}

export interface TGetSqlRefreshInput {
    refreshId: number
    orgId: number
}

export interface TGetSqlRefreshOutput {
    data: DbRefresh
}

export function getSqlRefresh(inp : TGetSqlRefreshInput) : Promise<TGetSqlRefreshOutput> {
    return axios.get(getSqlRefreshUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then(
        (resp : TGetSqlRefreshOutput) => {
            cleanDbRefreshFromJson(resp.data)
            return resp
        })
}
