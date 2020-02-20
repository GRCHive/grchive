import axios from 'axios'
import * as qs from 'query-string'
import { 
    allRiskSystemLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { System } from '../systems'

export interface TAllRiskSystemLinkInput {
    riskId?: number
    systemId?: number
    orgId: number
}

export interface TAllRiskSystemLinkOutput {
    data: {
        Systems?: System[]
        Risks?: ProcessFlowRisk[]
    }
}

export function allRiskSystemLink(inp : TAllRiskSystemLinkInput) : Promise<TAllRiskSystemLinkOutput> {
    return axios.get(allRiskSystemLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
