import { deleteProcessFlowAPIUrl } from '../url'
import { postFormUrlEncoded } from '../http'

export interface TDeleteProcessFlowInput {
    csrf: string
    flowId: number
}

export interface TDeleteProcessFlowOutput {
}

export function deleteProcessFlow(inp : TDeleteProcessFlowInput) : 
        Promise<TDeleteProcessFlowOutput> {
    return postFormUrlEncoded<TDeleteProcessFlowOutput>(deleteProcessFlowAPIUrl, inp)
}
