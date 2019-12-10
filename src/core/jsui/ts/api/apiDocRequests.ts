import axios from 'axios'
import * as qs from 'query-string'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { newDocRequestUrl,
         allDocRequestUrl,
         getDocRequestUrl } from '../url'
import { DocumentRequest } from '../docRequests'

export interface TNewDocRequestInput {
    name: string
    description: string
    catId: number
    orgId: number
    requestedUserId: number
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

export interface TGetAllDocumentRequestInput {
    orgId: number
    catId?: number
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
