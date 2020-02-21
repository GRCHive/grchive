import axios from 'axios'
import * as qs from 'query-string'
import { 
    allDocRequestDocCatLinksUrl 
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { ControlDocumentationCategory } from '../controls'
import { DocumentRequest } from '../docRequests'

export interface TAllDocRequestDocCatLinksInput {
    requestId?: number
    catId?: number
    orgId: number
}

export interface TAllDocRequestDocCatLinksOutput {
    data: {
        Cat?: ControlDocumentationCategory
        Requests?: DocumentRequest[]
    }
}

export function allDocRequestDocCatLink(inp : TAllDocRequestDocCatLinksInput) : Promise<TAllDocRequestDocCatLinksOutput> {
    return axios.get(allDocRequestDocCatLinksUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
