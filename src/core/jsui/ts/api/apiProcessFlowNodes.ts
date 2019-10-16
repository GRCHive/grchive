import axios from 'axios'
import * as qs from 'query-string'
import { getAllProcessFlowNodeTypesAPIUrl,
         editProcessFlowNodeAPIUrl,
         deleteProcessFlowNodeAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'

export function getProcessFlowNodeTypes(inp : TGetProcessFlowNodeTypesInput) : 
        Promise<TGetProcessFlowNodeTypesOutput> {
    return axios.get(getAllProcessFlowNodeTypesAPIUrl+ '?' + qs.stringify(inp))
}

export function editProcessFlowNode(inp : TEditProcessFlowNodeInput) : 
        Promise<TEditProcessFlowNodeOutput> {
    return postFormUrlEncoded<TEditProcessFlowNodeOutput>(editProcessFlowNodeAPIUrl, inp)
}

export function deleteProcessFlowNode(inp : TDeleteProcessFlowNodeInput) :
        Promise<TDeleteProcessFlowNodeOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowNodeOutput>(deleteProcessFlowNodeAPIUrl, inp)
}
