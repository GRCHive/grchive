import { 
    newFolderUrl,
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
