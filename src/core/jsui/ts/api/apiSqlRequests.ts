import axios from 'axios'
import * as qs from 'query-string'
import { 
    newSqlRequestUrl,
    updateSqlRequestUrl,
    deleteSqlRequestUrl,
    allSqlRequestUrl,
    getSqlRequestUrl,
    statusSqlRequestUrl,
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { 
    DbSqlQueryRequest,
    DbSqlQueryRequestApproval,
    cleanDbSqlRequestFromJson,
    cleanDbSqlRequestApprovalFromJson
} from '../sql'

export interface TNewSqlRequestInput {
    queryId: number
    orgId : number
    name: string
    description: string
    assigneeUserId: number | null
    dueDate: Date | null
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

export interface TGetSqlRequestInput {
    requestId : number
    orgId : number
}

export interface TGetSqlRequestOutput {
    data: {
        Request: DbSqlQueryRequest,
        Approval: DbSqlQueryRequestApproval | null
    }
}

export function getSqlRequest(inp : TGetSqlRequestInput) : Promise<TGetSqlRequestOutput> {
    return axios.get(getSqlRequestUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetSqlRequestOutput) => {
        cleanDbSqlRequestFromJson(resp.data.Request)
        if (!!resp.data.Approval) {
            cleanDbSqlRequestApprovalFromJson(resp.data.Approval)
        }
        return resp
    })
}

export interface TStatusSqlRequestInput {
    requestId : number
    orgId : number
}

export interface TStatusSqlRequestOutput {
    data: DbSqlQueryRequestApproval | null
}

export function statusSqlRequest(inp : TStatusSqlRequestInput) : Promise<TStatusSqlRequestOutput> {
    return axios.get(statusSqlRequestUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TStatusSqlRequestOutput) => {
        if (!!resp.data) {
            cleanDbSqlRequestApprovalFromJson(resp.data)
        }
        return resp
    })
}

export interface TUpdateSqlRequestInput {
    requestId: number
    orgId : number
    name: string
    description: string
    assigneeUserId: number | null
    dueDate: Date | null
}

export interface TUpdateSqlRequestOutput {
    data: DbSqlQueryRequest
}

export function updateSqlRequest(inp : TUpdateSqlRequestInput) : Promise<TUpdateSqlRequestOutput> {
    return postFormJson<TUpdateSqlRequestOutput>(updateSqlRequestUrl, inp, getAPIRequestConfig()).then((resp : TUpdateSqlRequestOutput) => {
        cleanDbSqlRequestFromJson(resp.data)
        return resp
    })
}

export interface TDeleteSqlRequestInput {
    requestId: number
    orgId : number
}

export function deleteSqlRequest(inp : TDeleteSqlRequestInput) : Promise<void> {
    return postFormJson<void>(deleteSqlRequestUrl, inp, getAPIRequestConfig())
}

export interface TModifyStatusSqlRequestInput {
    requestId : number
    orgId : number
    approve: boolean
    reason : string
}

export interface TModifyStatusSqlRequestOutput {
    data: DbSqlQueryRequestApproval
}

export function modifyStatusSqlRequest(inp : TModifyStatusSqlRequestInput) : Promise<TModifyStatusSqlRequestOutput> {
    return postFormJson<TModifyStatusSqlRequestOutput>(statusSqlRequestUrl, inp, getAPIRequestConfig()).then((resp : TModifyStatusSqlRequestOutput) => {
        cleanDbSqlRequestApprovalFromJson(resp.data)
        return resp
    })
}


