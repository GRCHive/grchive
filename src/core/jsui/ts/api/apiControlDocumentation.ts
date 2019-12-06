import axios from 'axios'
import * as qs from 'query-string'
import { postFormUrlEncoded, postFormMultipart } from '../http'
import { ControlDocumentationCategory, ControlDocumentationFile } from '../controls'
import { newControlDocCatUrl,
         editControlDocCatUrl,
         deleteControlDocCatUrl,
         uploadControlDocUrl,
         getControlDocUrl,
         deleteControlDocUrl,
         downloadControlDocUrl,
         allControlDocCatUrl,
         getControlDocCatUrl } from '../url'
import JSZip from 'jszip'
import { getAPIRequestConfig } from './apiUtility'

export interface TNewControlDocCatInput {
    name: string
    description: string
    orgId: number
}

export interface TNewControlDocCatOutput {
    data: ControlDocumentationCategory
}

export function newControlDocCat(inp : TNewControlDocCatInput): Promise<TNewControlDocCatOutput> {
    return postFormUrlEncoded<TNewControlDocCatOutput>(newControlDocCatUrl, inp, getAPIRequestConfig())
}

export interface TEditControlDocCatInput {
    catId: number
    name: string
    description: string
}

export interface TEditControlDocCatOutput {
    data: ControlDocumentationCategory
}

export function editControlDocCat(inp : TEditControlDocCatInput): Promise<TEditControlDocCatOutput> {
    return postFormUrlEncoded<TEditControlDocCatOutput>(editControlDocCatUrl, inp, getAPIRequestConfig())
}

export interface TDeleteControlDocCatInput {
    catId: number
    orgId: number
}

export interface TDeleteControlDocCatOutput {
}

export function deleteControlDocCat(inp : TDeleteControlDocCatInput): Promise<TDeleteControlDocCatOutput> {
    return postFormUrlEncoded<TDeleteControlDocCatOutput>(deleteControlDocCatUrl, inp, getAPIRequestConfig())
}


export interface TUploadControlDocOutput {
    data: ControlDocumentationFile
    orgId: number
}

export function uploadControlDoc(inp : FormData): Promise<TUploadControlDocOutput> {
    return postFormMultipart<TUploadControlDocOutput>(uploadControlDocUrl, inp, getAPIRequestConfig())
}

export interface TGetControlDocumentsInput {
    catId: number
    orgId: number
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
    return axios.get(getControlDocUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TDeleteControlDocumentsInput {
    fileIds: number[]
    orgId: number
}

export interface TDeleteControlDocumentsOutput {
}

export function deleteControlDocuments(inp: TDeleteControlDocumentsInput) : Promise<TDeleteControlDocumentsOutput> {
    return postFormUrlEncoded<TDeleteControlDocumentsOutput>(deleteControlDocUrl, inp, getAPIRequestConfig())
}

export interface TDownloadControlDocumentsInput {
    files: ControlDocumentationFile[]
    orgId: number
}

export interface TDownloadControlDocumentsOutput {
    data: Blob
}

export function downloadControlDocuments(inp: TDownloadControlDocumentsInput) : Promise<TDownloadControlDocumentsOutput> {
    return new Promise(async (resolve, reject) => {
        let zip = new JSZip()
        for (let file of inp.files) {
            try {
                let blobData = await axios.get<Blob>(downloadControlDocUrl + '?' + qs.stringify({
                    fileId: file.Id,
                    orgId: inp.orgId
                }), {
                    ...getAPIRequestConfig(),
                    responseType: "blob"
                })

                zip.folder(file.RelevantTime.toDateString()).file(`${file.Id}-${file.StorageName}`, blobData.data)
            } catch (e) {
                reject(e)
                return
            }
        }

        zip.generateAsync({
            type:"blob"
        }).then((blob : Blob) => {
            resolve({ data: blob })
        })
    })
}

export interface TGetAllDocumentationCategoriesInput {
    orgId: number
}

export interface TGetAllDocumentationCategoriesOutput {
    data: ControlDocumentationCategory[]
}

export function getAllDocumentationCategories(inp : TGetAllDocumentationCategoriesInput) : Promise<TGetAllDocumentationCategoriesOutput> {
    return axios.get(allControlDocCatUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetDocCatInput {
    orgId: number
    catId: number
}

export interface TGetDocCatOutput {
    data: ControlDocumentationCategory
}

export function getDocumentCategory(inp : TGetDocCatInput) : Promise<TGetDocCatOutput> {
    return axios.get(getControlDocCatUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
