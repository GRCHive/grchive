import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { 
    saveCodeUrl,
    allCodeUrl,
    getCodeUrl,
} from '../url'
import { ManagedCode, cleanManagedCodeFromJson } from '../code'

export interface TSaveCodeInput {
    orgId: number
    code : string
    dataId? : number
}

export interface TSaveCodeOutput {
    data : ManagedCode
}

export function saveCode(inp : TSaveCodeInput) : Promise<TSaveCodeOutput> {
    return postFormJson<TSaveCodeOutput>(saveCodeUrl, inp, getAPIRequestConfig()).then((resp : TSaveCodeOutput) => {
        cleanManagedCodeFromJson(resp.data)
        return resp
    })
}

export interface TAllCodeInput {
    orgId: number
    dataId? : number
}

export interface TAllCodeOutput {
    data : ManagedCode[]
}

export function allCode(inp : TAllCodeInput) : Promise<TAllCodeOutput> {
    return axios.get(allCodeUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllCodeOutput) => {
        resp.data.forEach(cleanManagedCodeFromJson)
        return resp
    })
}

export interface TGetCodeInput {
    orgId: number
    dataId? : number
    codeId : number
}

export interface TGetCodeOutput {
    data: string
}

export function getCode(inp : TGetCodeInput) : Promise<TGetCodeOutput> {
    return axios.get(getCodeUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
