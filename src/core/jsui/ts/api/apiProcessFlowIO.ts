import axios from 'axios'
import * as qs from 'query-string'
import { getAllProcessFlowIOTypesAPIUrl } from '../url'

export function getProcessFlowIOTypes(inp : TGetProcessFlowIOTypesInput) : 
        Promise<TGetProcessFlowIOTypesOutput> {
    return axios.get(getAllProcessFlowIOTypesAPIUrl + '?' + qs.stringify(inp))
}
