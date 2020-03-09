import axios from 'axios'
import * as qs from 'query-string'
import { 
    allAuditTrailLinkUrl,
    getAuditTrailLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { AuditEventEntry, cleanAuditEventEntryFromJson, AuditTrailFilterData } from '../auditTrail'
import { ResourceHandle } from '../resourceUtils'

export interface TAllAuditTrailInput {
    orgId: number
    page: number
    numItems: number
    sortHeaders: string
    sortDesc: boolean
    filter: AuditTrailFilterData
}

export interface TAllAuditTrailOutput {
    data: {
        Entries: AuditEventEntry[]
        Total: number
    }
}

export function allAuditTrail(inp : TAllAuditTrailInput) : Promise<TAllAuditTrailOutput> {
    let passData : any = {
        ...inp
    }
    passData.filter = JSON.stringify(passData.filter)
        
    return axios.get(allAuditTrailLinkUrl + '?' + qs.stringify(passData), getAPIRequestConfig()).then((resp : TAllAuditTrailOutput) => {
        resp.data.Entries.forEach(cleanAuditEventEntryFromJson)
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
