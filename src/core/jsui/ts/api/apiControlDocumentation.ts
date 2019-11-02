import axios from 'axios'
import * as qs from 'query-string'
import { postFormUrlEncoded, postFormMultipart } from '../http'
import { ControlDocumentationCategory, ControlDocumentationFile } from '../controls'
import { newControlDocCatUrl,
         editControlDocCatUrl,
         deleteControlDocCatUrl,
         uploadControlDocUrl,
         getControlDocUrl } from '../url'

export interface TNewControlDocCatInput {
    csrf: string
    controlId: number
    name: string
    description: string
}

export interface TNewControlDocCatOutput {
    data: ControlDocumentationCategory
}

export function newControlDocCat(inp : TNewControlDocCatInput): Promise<TNewControlDocCatOutput> {
    return postFormUrlEncoded<TNewControlDocCatOutput>(newControlDocCatUrl, inp)
}

export interface TEditControlDocCatInput {
    csrf: string
    catId: number
    name: string
    description: string
}

export interface TEditControlDocCatOutput {
    data: ControlDocumentationCategory
}

export function editControlDocCat(inp : TEditControlDocCatInput): Promise<TEditControlDocCatOutput> {
    return postFormUrlEncoded<TEditControlDocCatOutput>(editControlDocCatUrl, inp)
}

export interface TDeleteControlDocCatInput {
    csrf: string
    catId: number
}

export interface TDeleteControlDocCatOutput {
}

export function deleteControlDocCat(inp : TDeleteControlDocCatInput): Promise<TDeleteControlDocCatOutput> {
    return postFormUrlEncoded<TDeleteControlDocCatOutput>(deleteControlDocCatUrl, inp)
}


export interface TUploadControlDocOutput {
    data: ControlDocumentationFile
}

export function uploadControlDoc(inp : FormData): Promise<TUploadControlDocOutput> {
    return postFormMultipart<TUploadControlDocOutput>(uploadControlDocUrl, inp)
}

export interface TGetControlDocumentsInput {
    csrf: string
    catId: number
    page: number
    needPages: boolean
}

export interface TGetControlDocumentsOutput {
    data: {
        Files: ControlDocumentationFile[]
        TotalPages: number
        CurrentPage: number
    }
}

export function getControlDocuments(inp: TGetControlDocumentsInput) : Promise<TGetControlDocumentsOutput> {
    return axios.get(getControlDocUrl + '?' + qs.stringify(inp))
}
