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
    flowId: number
}

export interface TDeleteProcessFlowOutput {
}

export function deleteProcessFlow(inp : TDeleteProcessFlowInput) : 
        Promise<TDeleteProcessFlowOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowOutput>(deleteProcessFlowAPIUrl, inp, getAPIRequestConfig())
}

export interface TUpdateProcessFlowInput {
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
    name: string
    description: string
    organization: string
}

export interface TNewProcessFlowOutput {
    data: ProcessFlowBasicData
}

export function newProcessFlow(inp : TNewProcessFlowInput) : 
        Promise<TNewProcessFlowOutput> {
    return postFormUrlEncoded<TNewProcessFlowOutput>(newProcessFlowAPIUrl, inp, getAPIRequestConfig()).then(
        (resp : TNewProcessFlowOutput) => {
            resp.data.CreationTime = new Date(resp.data.CreationTime)
            resp.data.LastUpdatedTime = new Date(resp.data.LastUpdatedTime)
            return resp
        })
}

export interface TGetAllProcessFlowInput {
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
    return axios.get(getAllProcessFlowAPIUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then(
        (resp : TGetAllProcessFlowOutput) => {
            for (let d of resp.data.Flows) {
                d.CreationTime = new Date(d.CreationTime)
                d.LastUpdatedTime = new Date(d.LastUpdatedTime)
            }
            return resp
        })
}

export interface TGetFullProcessFlowInput {
}

export interface TGetFullProcessFlowOutput {
    data: {
        Basic: ProcessFlowBasicData
        Graph: FullProcessFlowResponseData
    }
}

export function getFullProcessFlow(id : number, inp : TGetFullProcessFlowInput) : 
        Promise<TGetFullProcessFlowOutput> {
    return axios.get(createGetProcessFlowFullDataUrl(id) + '?' + qs.stringify(inp), getAPIRequestConfig()).then(
        (resp : TGetFullProcessFlowOutput) => {
            resp.data.Basic.CreationTime = new Date(resp.data.Basic.CreationTime)
            resp.data.Basic.LastUpdatedTime = new Date(resp.data.Basic.LastUpdatedTime)
            return resp
        })
}
