import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import {
    allGenRequestScriptsUrl,
    allGenRequestsUrl
} from '../url'
import {
    GenericRequest,
    GenericApproval,
    cleanGenericRequestFromJson,
    cleanGenericApprovalFromJson
} from '../requests'
import { ClientScript } from '../../ts/clientScripts'
import { ManagedCode } from '../../ts/code'

export interface TAllGenericRequestsInput {
    orgId: number
    scriptsOnly: boolean
}

export interface TAllGenericRequestsOutput {
    data : GenericRequest[]
}

export function allGenericRequests(inp : TAllGenericRequestsInput) : Promise<TAllGenericRequestsOutput> {
    let url : string
    if (inp.scriptsOnly) {
        url = allGenRequestScriptsUrl
    } else {
        throw "Invalid parameters for retrieving generic requests."
    }

    return axios.get(url + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllGenericRequestsOutput) => {
        resp.data.forEach(cleanGenericRequestFromJson)
        return resp
    })
}

export interface TGetGenericRequestScriptInput {
    requestId: number
    orgId: number
}

export interface TGetGenericRequestScriptOutput {
    data: {
        Script: ClientScript
        Code: ManagedCode
        OneTime?: Date
        RRule? : string
        Params: Record<string, any>
    }
}

export function getGenericRequestScript(inp : TGetGenericRequestScriptInput) : Promise<TGetGenericRequestScriptOutput> {
    return axios.get(allGenRequestScriptsUrl + `/${inp.requestId}?` + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetGenericRequestScriptOutput) => {
        if (!!resp.data.OneTime) {
            resp.data.OneTime = new Date(resp.data.OneTime)
        }
        return resp
    })
}

export interface TGetGenericRequestInput {
    requestId: number
    orgId: number
}

export interface TGetGenericRequestOutput {
    data: {
        Request: GenericRequest
        Approval : GenericApproval | null
    }
}

export function getGenericRequest(inp : TGetGenericRequestInput) : Promise<TGetGenericRequestOutput> {
    return axios.get(allGenRequestsUrl + `/${inp.requestId}?` + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetGenericRequestOutput) => {
        cleanGenericRequestFromJson(resp.data.Request)
        if (!!resp.data.Approval) {
            cleanGenericApprovalFromJson(resp.data.Approval)
        }
        return resp
    })
}