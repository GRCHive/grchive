import axios from 'axios'
import * as qs from 'query-string'
import { deleteProcessFlowAPIUrl,
         createUpdateProcessFlowApiUrl,
         newProcessFlowAPIUrl,
         getAllProcessFlowAPIUrl,
         createGetProcessFlowFullDataUrl } from '../url'
import { postFormUrlEncoded } from '../http'
import { getAPIRequestConfig } from './apiUtility'

export interface TDeleteProcessFlowInput {
    csrf: string
    flowId: number
}

export interface TDeleteProcessFlowOutput {
}

export function deleteProcessFlow(inp : TDeleteProcessFlowInput) : 
        Promise<TDeleteProcessFlowOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowOutput>(deleteProcessFlowAPIUrl, inp, getAPIRequestConfig())
}

export interface TUpdateProcessFlowInput {
    csrf: string
    name: string
    description: string
}

export interface TUpdateProcessFlowOutput {
    data: ProcessFlowBasicData
}

export function updateProcessFlow(id : number, inp : TUpdateProcessFlowInput) : 
        Promise<TUpdateProcessFlowOutput> {
    return postFormUrlEncoded<TUpdateProcessFlowOutput>(createUpdateProcessFlowApiUrl(id), inp, getAPIRequestConfig())
}

export interface TNewProcessFlowInput {
    csrf: string
    name: string
    description: string
    organization: string
}

export interface TNewProcessFlowOutput {
    data: {
        Name : string
        Id: number
    }
}

export function newProcessFlow(inp : TNewProcessFlowInput) : 
        Promise<TNewProcessFlowOutput> {
    return postFormUrlEncoded<TNewProcessFlowOutput>(newProcessFlowAPIUrl, inp, getAPIRequestConfig())
}

export interface TGetAllProcessFlowInput {
    csrf: string
    requested: number | null
    organization: string
}

export interface TGetAllProcessFlowOutput {
    data: {
        Flows: ProcessFlowBasicData[]
        RequestedIndex: number
    }
}

export function getAllProcessFlow(inp : TGetAllProcessFlowInput) : 
        Promise<TGetAllProcessFlowOutput> {
    return axios.get(getAllProcessFlowAPIUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetFullProcessFlowInput {
    csrf: string
}

export interface TGetFullProcessFlowOutput {
    data: FullProcessFlowResponseData
}

export function getFullProcessFlow(id : number, inp : TGetFullProcessFlowInput) : 
        Promise<TGetFullProcessFlowOutput> {
    return axios.get(createGetProcessFlowFullDataUrl(id) + '?' + qs.stringify(inp), getAPIRequestConfig())
}
