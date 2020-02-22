import axios from 'axios'
import * as qs from 'query-string'
import { 
    allDocRequestDocCatLinksUrl 
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { ControlDocumentationCategory } from '../controls'
import { DocumentRequest, cleanJsonDocumentRequest } from '../docRequests'

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
    return axios.get(allDocRequestDocCatLinksUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllDocRequestDocCatLinksOutput) => {
        if (!!resp.data.Requests) {
            resp.data.Requests.forEach(cleanJsonDocumentRequest)
        }
        return resp
    })

}
