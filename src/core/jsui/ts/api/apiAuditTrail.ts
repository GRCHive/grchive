import axios from 'axios'
import * as qs from 'query-string'
import { 
    allAuditTrailLinkUrl,
    getAuditTrailLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { AuditEventEntry, cleanAuditEventEntryFromJson } from '../auditTrail'
import { ResourceHandle } from '../resourceUtils'

export interface TAllAuditTrailInput {
    orgId: number
}

export interface TAllAuditTrailOutput {
    data: AuditEventEntry[]
}

export function allAuditTrail(inp : TAllAuditTrailInput) : Promise<TAllAuditTrailOutput> {
    return axios.get(allAuditTrailLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllAuditTrailOutput) => {
        resp.data.forEach(cleanAuditEventEntryFromJson)
        return resp
    })
}

export interface TGetAuditTrailInput {
    orgId: number
    resourceHandleOnly: boolean
    entryId? : number
}

export interface TGetAuditTrailOutput {
    data: {
        Handle?: ResourceHandle
    }
}

export function getAuditTrail(inp : TGetAuditTrailInput) : Promise<TGetAuditTrailOutput> {
    return axios.get(getAuditTrailLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
