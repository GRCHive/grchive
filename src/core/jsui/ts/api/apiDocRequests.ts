import axios from 'axios'
import * as qs from 'query-string'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { newDocRequestUrl,
         allDocRequestUrl,
         getDocRequestUrl,
         deleteDocRequestUrl,
         updateDocRequestUrl,
         apiv2DocRequestFileLinks,
         apiv2DocRequestReopen,
         apiv2DocRequestComplete,
         apiv2DocRequestApprove,
} from '../url'
import { DocumentRequest, cleanJsonDocumentRequest, DocRequestFilterData } from '../docRequests'
import { ControlDocumentationCategory, ControlDocumentationFile, cleanJsonControlDocumentationFile } from '../controls'

export interface TNewDocRequestInput {
    name: string
    description: string
    catId: number
    controlId? : number
    folderId?: number
    orgId: number
    requestedUserId: number
    assigneeUserId: number | null
    dueDate: Date | null
    vendorProductId: number
}

export interface TNewDocRequestOutput {
    data: {
        Request: DocumentRequest
    }
}

export function newDocRequest(inp : TNewDocRequestInput) : Promise<TNewDocRequestOutput> {
    return postFormJson<TNewDocRequestOutput>(newDocRequestUrl, inp, getAPIRequestConfig()).then((resp : TNewDocRequestOutput) => {
        cleanJsonDocumentRequest(resp.data.Request)
        return resp
    })
}

export interface TUpdateDocRequestInput extends TNewDocRequestInput {
    requestId: number
}

export interface TUpdateDocRequestOutput {
    data: {
        Request: DocumentRequest
    }
}

export function updateDocRequest(inp : TUpdateDocRequestInput) : Promise<TUpdateDocRequestOutput> {
    return postFormJson<TUpdateDocRequestOutput>(updateDocRequestUrl, inp, getAPIRequestConfig()).then((resp : TUpdateDocRequestOutput) => {
        cleanJsonDocumentRequest(resp.data.Request)
        return resp
    })
}

export interface TGetAllDocumentRequestInput {
    orgId: number
    catId?: number
    vendorProductId?: number
    filter: DocRequestFilterData
}

export interface TGetAllDocumentRequestOutput {
    data: DocumentRequest[]
}

export function getAllDocRequests(inp : TGetAllDocumentRequestInput) : Promise<TGetAllDocumentRequestOutput> {
    let passData : any = {
        ...inp,
        filter: JSON.stringify(inp.filter),
    }

    return axios.get(allDocRequestUrl + '?' + qs.stringify(passData), getAPIRequestConfig()).then((resp : TGetAllDocumentRequestOutput) => {
        resp.data.forEach(cleanJsonDocumentRequest)
        return resp
    })
}

export interface TGetSingleDocumentRequestInput {
    requestId: number
    orgId: number
}

export interface TGetSingleDocumentRequestOutput {
    data: {
        Request: DocumentRequest
        Files: ControlDocumentationFile[]
        VendorProductId: number
        VendorId: number
    }
}

export function getSingleDocRequest(inp : TGetSingleDocumentRequestInput) : Promise<TGetSingleDocumentRequestOutput> {
    return axios.get(getDocRequestUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetSingleDocumentRequestOutput) => {
        cleanJsonDocumentRequest(resp.data.Request)
        resp.data.Files.forEach(cleanJsonControlDocumentationFile)
        return resp
    })
}

export interface TDeleteDocumentRequestInput {
    requestId: number
    orgId: number
}

export function deleteSingleDocRequest(inp : TDeleteDocumentRequestInput) : Promise<void> {
    return postFormJson(deleteDocRequestUrl, inp, getAPIRequestConfig())
}

export interface TDocumentRequestStatusInput {
    requestId: number
    orgId: number
}

export function completeDocRequest(inp : TDocumentRequestStatusInput) : Promise<void> {
    return postFormJson(apiv2DocRequestComplete(inp.orgId, inp.requestId), {}, getAPIRequestConfig())
}

export function reopenDocRequest(inp : TDocumentRequestStatusInput) : Promise<void> {
    return postFormJson(apiv2DocRequestReopen(inp.orgId, inp.requestId), {}, getAPIRequestConfig())
}

export function approveDocRequest(inp : TDocumentRequestStatusInput) : Promise<void> {
    return postFormJson(apiv2DocRequestApprove(inp.orgId, inp.requestId), {}, getAPIRequestConfig())
}

export interface TLinkFilesToDocumentRequestInput {
    requestId: number
    orgId : number
    files : number[]
}

export function linkFilesToDocRequest(inp : TLinkFilesToDocumentRequestInput) : Promise<void> {
    return postFormJson(
        apiv2DocRequestFileLinks(inp.orgId, inp.requestId),
        { files: inp.files },
        getAPIRequestConfig(),
    )
}
