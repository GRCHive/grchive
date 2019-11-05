import { newProcessFlowEdgeAPIUrl, deleteProcessFlowEdgeAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'
import { getAPIRequestConfig } from './apiUtility'

export function newProcessFlowEdge(inp : TNewProcessFlowEdgeInput) : Promise<TNewProcessFlowEdgeOutput> {
    return postFormUrlEncoded<TNewProcessFlowEdgeOutput>(newProcessFlowEdgeAPIUrl, inp, getAPIRequestConfig())
}

export function deleteProcessFlowEdge(inp : TDeleteProcessFlowEdgeInput) : Promise<TDeleteProcessFlowEdgeOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowEdgeOutput>(deleteProcessFlowEdgeAPIUrl, inp, getAPIRequestConfig())
}
