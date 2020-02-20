import axios from 'axios'
import * as qs from 'query-string'
import { 
    newNodeSystemLinkUrl,
    deleteNodeSystemLinkUrl,
    allNodeSystemLinkUrl,
} from '../url'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { System } from '../systems'

export interface TNewNodeSystemLinkInput {
    nodeId: number
    systemId: number
    orgId: number
}

export function newNodeSystemLink(inp : TNewNodeSystemLinkInput) : Promise<void> {
    return postFormJson<void>(newNodeSystemLinkUrl, inp, getAPIRequestConfig())
}

export interface TDeleteNodeSystemLinkInput {
    nodeId: number
    systemId: number
    orgId: number
}

export function deleteNodeSystemLink(inp : TDeleteNodeSystemLinkInput) : Promise<void> {
    return postFormJson<void>(deleteNodeSystemLinkUrl, inp, getAPIRequestConfig())
}

export interface TAllNodeSystemLinkInput {
    nodeId?: number
    systemId?: number
    orgId: number
}

export interface TAllNodeSystemLinkOutput {
    data: System[] | ProcessFlowNode[]
}

export function allNodeSystemLink(inp : TAllNodeSystemLinkInput) : Promise<TAllNodeSystemLinkOutput> {
    return axios.get(allNodeSystemLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
