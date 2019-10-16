import axios from 'axios'
import * as qs from 'query-string'
import { getAllProcessFlowIOTypesAPIUrl, 
         deleteProcessFlowIOAPIUrl,
         editProcessFlowIOAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'

export function getProcessFlowIOTypes(inp : TGetProcessFlowIOTypesInput) : 
        Promise<TGetProcessFlowIOTypesOutput> {
    return axios.get(getAllProcessFlowIOTypesAPIUrl + '?' + qs.stringify(inp))
}

export function deleteProcessFlowIO(inp : TDeleteProcessFlowIOInput) : Promise<TDeleteProcessFlowIOOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowIOOutput>(deleteProcessFlowIOAPIUrl, inp)
}

export function editProcessFlowIO(inp : TEditProcessFlowIOInput) : Promise<TEditProcessFlowIOOutput> {
    return postFormUrlEncoded<TEditProcessFlowIOOutput>(editProcessFlowIOAPIUrl, inp)
}