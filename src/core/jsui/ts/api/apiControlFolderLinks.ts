import axios from 'axios'
import * as qs from 'query-string'
import { 
    allControlFolderLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { FileFolder } from '../folders'

export interface TAllControlFolderLinkInput {
    controlId?: number
    folderId?: number
    orgId: number
}

export interface TAllControlFolderLinkOutput {
    data: {
        Folders?: FileFolder[]
        Controls?: ProcessFlowControl[]
    }
}

export function allControlFolderLink(inp : TAllControlFolderLinkInput) : Promise<TAllControlFolderLinkOutput> {
    return axios.get(allControlFolderLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
