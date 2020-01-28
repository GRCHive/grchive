import axios from 'axios'
import * as qs from 'query-string'
import { postFormUrlEncoded, postFormMultipart, postFormJson } from '../http'
import {
    ControlDocumentationCategory, 
    ControlDocumentationFile,
    FileVersion,
    cleanJsonControlDocumentationFile
} from '../controls'
import { newControlDocCatUrl,
         editControlDocCatUrl,
         deleteControlDocCatUrl,
         uploadControlDocUrl,
         allControlDocUrl,
         getControlDocUrl,
         allControlDocVersionsUrl,
         editControlDocUrl,
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

export interface TUploadControlDocInput {
    catId: number
    orgId: number
    file: File
    relevantTime: Date
    altName: string
    description: string
    uploadUserId: number
    fulfilledRequestId?: number | null
}

export interface TUploadControlDocOutput {
    data: ControlDocumentationFile
}

export function uploadControlDoc(inp : TUploadControlDocInput): Promise<TUploadControlDocOutput> {
    let data = new FormData()
    data.set('catId', inp.catId.toString())
    data.set('orgId', inp.orgId.toString())
    data.set('file', inp.file)
    data.set('relevantTime', inp.relevantTime.toISOString())
    data.set('altName', inp.altName)
    data.set('description', inp.description)
    data.set('uploadUserId', inp.uploadUserId.toString())
    if (!!inp.fulfilledRequestId) {
        data.set('fulfilledRequestId', inp.fulfilledRequestId!.toString())
    }

    return postFormMultipart<TUploadControlDocOutput>(uploadControlDocUrl, data, getAPIRequestConfig()).then((resp : TUploadControlDocOutput) => {
        cleanJsonControlDocumentationFile(resp.data)
        return resp
    })
}

export interface TEditControlDocInput {
    fileId: number
    orgId: number
    relevantTime: Date
    altName: string
    description: string
    uploadUserId: number
}

export interface TEditControlDocOutput {
    data: {
        File: ControlDocumentationFile
        Category: ControlDocumentationCategory
        UploadUser: User
    }
}

export function editControlDoc(inp : TEditControlDocInput) : Promise<TEditControlDocOutput> {
    return postFormJson<TEditControlDocOutput>(editControlDocUrl, inp, getAPIRequestConfig()).then((resp : TEditControlDocOutput) => {
        cleanJsonControlDocumentationFile(resp.data.File)
        return resp
    })
}


export interface TAllControlDocumentsInput {
    catId: number
    orgId: number
}

export interface TAllControlDocumentsOutput {
    data: {
        Files: ControlDocumentationFile[]
    }
}

export function allControlDocuments(inp: TAllControlDocumentsInput) : Promise<TAllControlDocumentsOutput> {
    return axios.get(allControlDocUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllControlDocumentsOutput) => {
        resp.data.Files = resp.data.Files.map((ele : ControlDocumentationFile) => {
            cleanJsonControlDocumentationFile(ele)
            return ele
        })
        return resp
    })
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

export interface TDownloadSingleControlDocumentInput {
    fileId: number
    orgId: number
}

export interface TDownloadSingleControlDocumentOutput {
    data: Blob
}

export function downloadSingleControlDocument(inp : TDownloadSingleControlDocumentInput) : Promise<TDownloadSingleControlDocumentOutput> {
    return axios.get<Blob>(downloadControlDocUrl + '?' + qs.stringify(inp), {
        ...getAPIRequestConfig(),
        responseType: "blob"
    })
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
                let blobData = await downloadSingleControlDocument({
                    fileId: file.Id,
                    orgId: inp.orgId,
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
    lean: boolean
}

export interface TGetDocCatOutput {
    data: {
        Cat: ControlDocumentationCategory
        InputFor: ProcessFlowControl[]
        OutputFor: ProcessFlowControl[]
    }
}

export function getDocumentCategory(inp : TGetDocCatInput) : Promise<TGetDocCatOutput> {
    return axios.get(getControlDocCatUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TGetSingleControlDocumentInput {
    fileId: number
    orgId: number
}

export interface TGetSingleControlDocumentOutput {
    data: {
        File: ControlDocumentationFile
        Category: ControlDocumentationCategory
        PreviewFile: ControlDocumentationFile | null
        UploadUser: User
    }
}

export function getSingleControlDocument(inp: TGetSingleControlDocumentInput) : Promise<TGetSingleControlDocumentOutput> {
    return axios.get(getControlDocUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TAllFileVersionsInput {
    fileId: number
    orgId: number
}

export interface TAllFileVersionsOutput {
    data: FileVersion[]
}

export function allFileVersions(inp : TAllFileVersionsInput) : Promise<TAllFileVersionsOutput> {
    return axios.get(allControlDocVersionsUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
