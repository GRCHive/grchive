import axios from 'axios'
import * as qs from 'query-string'
import { 
    newNodeGLLinkUrl,
    deleteNodeGLLinkUrl,
    allNodeGLLinkUrl,
} from '../url'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { RawGeneralLedgerAccount, RawGeneralLedgerCategory } from '../generalLedger'
import { cleanProcessFlowFromJson } from '../processFlow'

export interface TNewNodeGLLinkInput {
    nodeId: number
    accountId: number
    orgId: number
}

export function newNodeGLLink(inp : TNewNodeGLLinkInput) : Promise<void> {
    return postFormJson<void>(newNodeGLLinkUrl, inp, getAPIRequestConfig())
}

export interface TDeleteNodeGLLinkInput {
    nodeId: number
    accountId: number
    orgId: number
}

export function deleteNodeGLLink(inp : TDeleteNodeGLLinkInput) : Promise<void> {
    return postFormJson<void>(deleteNodeGLLinkUrl, inp, getAPIRequestConfig())
}

export interface TAllNodeGLLinkInput {
    nodeId?: number
    accountId?: number
    orgId: number
}

export interface TAllNodeGLLinkOutput {
    data: {
        Accounts?: RawGeneralLedgerAccount[]
        Categories?: RawGeneralLedgerCategory[]
        Flows?:  ProcessFlowBasicData[]
    }
}

export function allNodeGLLink(inp : TAllNodeGLLinkInput) : Promise<TAllNodeGLLinkOutput> {
    return axios.get(allNodeGLLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllNodeGLLinkOutput) => {
        if (!!resp.data.Flows) {
            resp.data.Flows!.forEach(cleanProcessFlowFromJson)
        }
        return resp
    })
}
