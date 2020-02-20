import axios from 'axios'
import * as qs from 'query-string'
import { 
    allControlSystemLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { System } from '../systems'

export interface TAllControlSystemLinkInput {
    controlId?: number
    systemId?: number
    orgId: number
}

export interface TAllControlSystemLinkOutput {
    data: {
        Systems?: System[]
        Controls?: ProcessFlowControl[]
    }
}

export function allControlSystemLink(inp : TAllControlSystemLinkInput) : Promise<TAllControlSystemLinkOutput> {
    return axios.get(allControlSystemLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
