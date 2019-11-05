import axios from 'axios'
import * as qs from 'query-string'
import { getAllProcessFlowNodeTypesAPIUrl,
         editProcessFlowNodeAPIUrl,
         deleteProcessFlowNodeAPIUrl,
         newProcessFlowNodeAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'
import { getAPIRequestConfig } from './apiUtility'

export interface TGetProcessFlowNodeTypesInput { 
}

export interface TGetProcessFlowNodeTypesOutput { 
    data : ProcessFlowNodeType[]
}

export function getProcessFlowNodeTypes(inp : TGetProcessFlowNodeTypesInput) : 
        Promise<TGetProcessFlowNodeTypesOutput> {
    return axios.get(getAllProcessFlowNodeTypesAPIUrl+ '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TEditProcessFlowNodeInput { 
    nodeId: number,
    name: string,
    description: string,
    type: number
}

export interface TEditProcessFlowNodeOutput { 
    data: ProcessFlowNode
}

export function editProcessFlowNode(inp : TEditProcessFlowNodeInput) : 
        Promise<TEditProcessFlowNodeOutput> {
    return postFormUrlEncoded<TEditProcessFlowNodeOutput>(editProcessFlowNodeAPIUrl, inp, getAPIRequestConfig())
}

export interface TDeleteProcessFlowNodeInput { 
    nodeId: number
}

export interface TDeleteProcessFlowNodeOutput { 
}

export function deleteProcessFlowNode(inp : TDeleteProcessFlowNodeInput) :
        Promise<TDeleteProcessFlowNodeOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowNodeOutput>(deleteProcessFlowNodeAPIUrl, inp, getAPIRequestConfig())
}

export interface TNewProcessFlowNodeInput {
    typeId: number
    flowId: number
}

export interface TNewProcessFlowNodeOutput {
    data: ProcessFlowNode
}

export function newProcessFlowNode(inp : TNewProcessFlowNodeInput) : 
        Promise<TNewProcessFlowNodeOutput> {
    return postFormUrlEncoded<TNewProcessFlowNodeOutput>(newProcessFlowNodeAPIUrl, inp, getAPIRequestConfig())
}
