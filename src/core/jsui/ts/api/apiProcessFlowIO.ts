import axios from 'axios'
import * as qs from 'query-string'
import { getAllProcessFlowIOTypesAPIUrl, 
         deleteProcessFlowIOAPIUrl,
         editProcessFlowIOAPIUrl,
         newProcessFlowIOAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'
import { getAPIRequestConfig } from './apiUtility'

export function getProcessFlowIOTypes(inp : TGetProcessFlowIOTypesInput) : 
        Promise<TGetProcessFlowIOTypesOutput> {
    return axios.get(getAllProcessFlowIOTypesAPIUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export function deleteProcessFlowIO(inp : TDeleteProcessFlowIOInput) : Promise<TDeleteProcessFlowIOOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowIOOutput>(deleteProcessFlowIOAPIUrl, inp, getAPIRequestConfig())
}

export function editProcessFlowIO(inp : TEditProcessFlowIOInput) : Promise<TEditProcessFlowIOOutput> {
    return postFormUrlEncoded<TEditProcessFlowIOOutput>(editProcessFlowIOAPIUrl, inp, getAPIRequestConfig())
}


export interface TNewProcessFlowIOInput {
    csrf: string
    name: string
    isInput: boolean
    nodeId: number
    typeId: number
}

export interface TNewProcessFlowIOOutput {
    data: ProcessFlowInputOutput
}

export function newProcessFlowIO(inp : TNewProcessFlowIOInput) : 
        Promise<TNewProcessFlowIOOutput> {
    return postFormUrlEncoded<TNewProcessFlowIOOutput>(newProcessFlowIOAPIUrl, inp, getAPIRequestConfig())
}

