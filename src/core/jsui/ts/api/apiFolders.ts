import { 
    newFolderUrl,
    updateFolderUrl,
    deleteFolderUrl,
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { FileFolder } from '../folders'

export interface TNewFolderInput {
    name : string
    orgId: number
    controlId: number
}

export interface TNewFolderOutput {
    data: FileFolder
}

export function newFolder(inp : TNewFolderInput) : Promise<TNewFolderOutput> {
    return postFormJson<TNewFolderOutput>(newFolderUrl, inp, getAPIRequestConfig())
}

export interface TUpdateFolderInput extends TNewFolderInput {
    folderId: number
}

export interface TUpdateFolderOutput {
    data: FileFolder
}

export function updateFolder(inp : TUpdateFolderInput) : Promise<TUpdateFolderOutput> {
    return postFormJson<TUpdateFolderOutput>(updateFolderUrl, inp, getAPIRequestConfig())
}

export interface TDeleteFolderInput {
    orgId: number
    folderId: number
}

export function deleteFolder(inp : TDeleteFolderInput) : Promise<void> {
    return postFormJson<void>(deleteFolderUrl, inp, getAPIRequestConfig())
}
