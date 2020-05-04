import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { CodeRunTestSummary } from '../tests'
import {
    getCodeRunTestsUrl
} from '../url'

export interface TGetCodeRunTestInput {
    orgId: number
    runId: number
    summary: boolean
}

export interface TGetCodeRunTestOutput {
    data : CodeRunTestSummary
}

export function getCodeRunTest(inp : TGetCodeRunTestInput) : Promise<TGetCodeRunTestOutput> {
    return axios.get(getCodeRunTestsUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
