import axios from 'axios'
import * as qs from 'query-string'
import { allAuditTrailLinkUrl } from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { AuditEventEntry, cleanAuditEventEntryFromJson } from '../auditTrail'

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
