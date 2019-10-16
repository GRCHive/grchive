import { newProcessFlowEdgeAPIUrl, deleteProcessFlowEdgeAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'

export function newProcessFlowEdge(inp : TNewProcessFlowEdgeInput) : Promise<TNewProcessFlowEdgeOutput> {
    return postFormUrlEncoded<TNewProcessFlowEdgeOutput>(newProcessFlowEdgeAPIUrl, inp)
}

export function deleteProcessFlowEdge(inp : TDeleteProcessFlowEdgeInput) : Promise<TDeleteProcessFlowEdgeOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowEdgeOutput>(deleteProcessFlowEdgeAPIUrl, inp)
}
