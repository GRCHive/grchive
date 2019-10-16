import { newProcessFlowEdgeAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'

export function newProcessFlowEdge(inp : TNewProcessFlowEdgeInput) : Promise<TNewProcessFlowEdgeOutput> {
    return postFormUrlEncoded<TNewProcessFlowEdgeOutput>(newProcessFlowEdgeAPIUrl, inp)
}
