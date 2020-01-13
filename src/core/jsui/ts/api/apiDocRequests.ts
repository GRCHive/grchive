import axios from 'axios'
import * as qs from 'query-string'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { newDocRequestUrl,
         allDocRequestUrl,
         getDocRequestUrl,
         deleteDocRequestUrl,
         completeDocRequestUrl,
         updateDocRequestUrl } from '../url'
import { DocumentRequest } from '../docRequests'
import { ControlDocumentationCategory, ControlDocumentationFile } from '../controls'

export interface TNewDocRequestInput {
    name: string
    description: string
    catId: number
    orgId: number
    requestedUserId: number
    vendorProductId: number
}

export interface TNewDocRequestOutput {
    data: DocumentRequest
}

export function newDocRequest(inp : TNewDocRequestInput) : Promise<TNewDocRequestOutput> {
    return postFormJson<TNewDocRequestOutput>(newDocRequestUrl, inp, getAPIRequestConfig()).then((resp : TNewDocRequestOutput) => {
        if (!!resp.data.CompletionTime) {
            resp.data.CompletionTime = new Date(resp.data.CompletionTime)
        }
        resp.data.RequestTime = new Date(resp.data.RequestTime)
        return resp
    })
}

export interface TUpdateDocRequestInput extends TNewDocRequestInput {
    requestId: number
}

export interface TUpdateDocRequestOutput {
    data: DocumentRequest
}

export function updateDocRequest(inp : TUpdateDocRequestInput) : Promise<TUpdateDocRequestOutput> {
    return postFormJson<TUpdateDocRequestOutput>(updateDocRequestUrl, inp, getAPIRequestConfig()).then((resp : TUpdateDocRequestOutput) => {
        if (!!resp.data.CompletionTime) {
            resp.data.CompletionTime = new Date(resp.data.CompletionTime)
        }
        resp.data.RequestTime = new Date(resp.data.RequestTime)
        return resp
    })
}

export interface TGetAllDocumentRequestInput {
    orgId: number
    catId?: number
    vendorProductId?: number
}

export interface TGetAllDocumentRequestOutput {
    data: DocumentRequest[]
}

export function getAllDocRequests(inp : TGetAllDocumentRequestInput) : Promise<TGetAllDocumentRequestOutput> {
    return axios.get(allDocRequestUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetAllDocumentRequestOutput) => {
        resp.data = resp.data.map((ele : DocumentRequest) => {
            if (!!ele.CompletionTime) {
                ele.CompletionTime = new Date(ele.CompletionTime)
            }
            ele.RequestTime = new Date(ele.RequestTime)
            return ele
        })
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
        Category: ControlDocumentationCategory
    }
}

export function getSingleDocRequest(inp : TGetSingleDocumentRequestInput) : Promise<TGetSingleDocumentRequestOutput> {
    return axios.get(getDocRequestUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetSingleDocumentRequestOutput) => {
        if (!!resp.data.Request.CompletionTime) {
            resp.data.Request.CompletionTime = new Date(resp.data.Request.CompletionTime)
        }
        resp.data.Request.RequestTime = new Date(resp.data.Request.RequestTime)

        resp.data.Files = resp.data.Files.map((ele : ControlDocumentationFile) => {
            ele.UploadTime = new Date(ele.UploadTime)
            ele.RelevantTime = new Date(ele.RelevantTime)
            return ele
        })

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

export interface TCompleteDocumentRequestInput {
    requestId: number
    orgId: number
    complete: boolean
}

export function completeDocRequest(inp : TCompleteDocumentRequestInput) : Promise<void> {
    return postFormJson(completeDocRequestUrl, inp, getAPIRequestConfig())
}

