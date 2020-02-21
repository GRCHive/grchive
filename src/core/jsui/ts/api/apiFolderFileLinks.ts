import axios from 'axios'
import * as qs from 'query-string'
import { 
    allFolderFileLinkUrl,
    newFolderFileLinkUrl,
    deleteFolderFileLinkUrl,
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { ControlDocumentationFile, cleanJsonControlDocumentationFile } from '../controls'

export interface TAllFolderFileLinkInput {
    folderId?: number
    orgId: number
}

export interface TAllFolderFileLinkOutput {
    data: {
        Files?: ControlDocumentationFile[]
    }
}

export function allFolderFileLink(inp : TAllFolderFileLinkInput) : Promise<TAllFolderFileLinkOutput> {
    return axios.get(allFolderFileLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllFolderFileLinkOutput) => {
        if (!!resp.data.Files) {
            resp.data.Files.forEach(cleanJsonControlDocumentationFile)
        }
        return resp
    })
}

export interface TNewFolderFileLinkInput {
    folderId: number
    fileIds: number[]
    orgId: number
}

export function newFolderFileLink(inp : TNewFolderFileLinkInput) : Promise<void> {
    return postFormJson<void>(newFolderFileLinkUrl, inp, getAPIRequestConfig())
}

export interface TDeleteFolderFileLinkInput {
    folderId: number
    fileId: number
    orgId: number
}

export function deleteFolderFileLink(inp : TDeleteFolderFileLinkInput) : Promise<void> {
    return postFormJson<void>(deleteFolderFileLinkUrl, inp, getAPIRequestConfig())
}
