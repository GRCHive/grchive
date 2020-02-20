import axios from 'axios'
import * as qs from 'query-string'
import { 
    allRiskGLLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { RawGeneralLedgerAccount, RawGeneralLedgerCategory } from '../generalLedger'

export interface TAllRiskGLLinkInput {
    riskId?: number
    accountId?: number
    orgId: number
}

export interface TAllRiskGLLinkOutput {
    data: {
        Accounts?: RawGeneralLedgerAccount[]
        Categories?: RawGeneralLedgerCategory[]
        Risks?: ProcessFlowRisk[]
    }
}

export function allRiskGLLink(inp : TAllRiskGLLinkInput) : Promise<TAllRiskGLLinkOutput> {
    return axios.get(allRiskGLLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
