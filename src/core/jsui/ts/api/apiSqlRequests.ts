import axios from 'axios'
import * as qs from 'query-string'
import { 
    newSqlRequestUrl,
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { 
    DbSqlQueryRequest,
    cleanDbSqlRequestFromJson
} from '../sql'

export interface TNewSqlRequestInput {
    queryId: number
    orgId : number
    name: string
    description: string
}

export interface TNewSqlRequestOutput {
    data: DbSqlQueryRequest
}

export function newSqlRequest(inp : TNewSqlRequestInput) : Promise<TNewSqlRequestOutput> {
    return postFormJson<TNewSqlRequestOutput>(newSqlRequestUrl, inp, getAPIRequestConfig()).then((resp : TNewSqlRequestOutput) => {
        cleanDbSqlRequestFromJson(resp.data)
        return resp
    })
}
