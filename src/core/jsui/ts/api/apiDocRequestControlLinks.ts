import axios from 'axios'
import * as qs from 'query-string'
import { 
    allDocRequestControlLinksUrl 
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { DocumentRequest, cleanJsonDocumentRequest } from '../docRequests'

export interface TAllDocRequestControlLinksInput {
    requestId?: number
    controlId?: number
    orgId: number
}

export interface TAllDocRequestControlLinksOutput {
    data: {
        Control?: ProcessFlowControl
        Requests?: DocumentRequest[]
    }
}

export function allDocRequestControlLink(inp : TAllDocRequestControlLinksInput) : Promise<TAllDocRequestControlLinksOutput> {
    return axios.get(allDocRequestControlLinksUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllDocRequestControlLinksOutput) => {
        if (!!resp.data.Requests) {
            resp.data.Requests.forEach(cleanJsonDocumentRequest)
        }
        return resp
    })
}
