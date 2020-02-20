import axios from 'axios'
import * as qs from 'query-string'
import { 
    allControlGLLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { RawGeneralLedgerAccount, RawGeneralLedgerCategory } from '../generalLedger'

export interface TAllControlGLLinkInput {
    controlId?: number
    accountId?: number
    orgId: number
}

export interface TAllControlGLLinkOutput {
    data: {
        Accounts?: RawGeneralLedgerAccount[]
        Categories?: RawGeneralLedgerCategory[]
        Controls?: ProcessFlowControl[]
    }
}

export function allControlGLLink(inp : TAllControlGLLinkInput) : Promise<TAllControlGLLinkOutput> {
    return axios.get(allControlGLLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
