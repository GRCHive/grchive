import axios from 'axios'
import * as qs from 'query-string'
import { ResourceHandle } from '../resourceUtils'
import {
    getResourceHandleUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'

export interface TGetResourceHandleInput {
    orgId : number
    resourceType: string
    resourceId: number
}

export interface TGetResourceHandleOutput {
    data: ResourceHandle
}

export function getResourceHandle(inp : TGetResourceHandleInput) : Promise<TGetResourceHandleOutput> {
    return axios.get(getResourceHandleUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
