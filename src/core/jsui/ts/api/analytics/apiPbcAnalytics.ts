import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from '../apiUtility'
import { DocRequestFilterData, DocRequestStatus } from '../../docRequests'
import {
    apiv2PbcOverallAnalyticsLink,
    apiv2PbcCategoryAnalyticsLink,
} from '../../url'

export interface TGetPbcOverallProgressInputs {
    orgId: number
    filter: DocRequestFilterData
}

export interface TGetPbcOverallProgressOutputs {
    data: Record<DocRequestStatus, number>
}

export function getPbcOverallProgress(inp : TGetPbcOverallProgressInputs) : Promise<TGetPbcOverallProgressOutputs> {
    return axios.get(
        apiv2PbcOverallAnalyticsLink(inp.orgId) + '?' + qs.stringify({
            filter: JSON.stringify(inp.filter),
        }),
        getAPIRequestConfig()
    )
}

export interface TGetPbcCategoryProgressInputs {
    orgId: number
    category: string
    filter: DocRequestFilterData
}

export interface TGetPbcCategoryProgressOutputs {
    data: Array<{
        Name: string
        Data: Record<DocRequestStatus, number>
    }>
}

export function getPbcCategoryProgress(inp : TGetPbcCategoryProgressInputs) : Promise<TGetPbcCategoryProgressOutputs> {
    return axios.get(
        apiv2PbcCategoryAnalyticsLink(inp.orgId, inp.category) + '?' + qs.stringify({
            filter: JSON.stringify(inp.filter),
        }),
        getAPIRequestConfig()
    )
}
