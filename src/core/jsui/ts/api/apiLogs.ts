import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import {
    getLogUrl
} from '../url'

export interface TGetLogInput {
    orgId: number
    commitHash? : string
    runId? : number
    runLog? : boolean
}

export interface TGetLogOutput {
    data: string
}

export function getLog(inp : TGetLogInput) : Promise<TGetLogOutput> {
    return axios.get(getLogUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
