import axios from 'axios'
import * as qs from 'query-string'
import { getAPIRequestConfig } from './apiUtility'
import { FileFolder } from '../folders'
import { 
    apiv2DocRequestControlFolderLinks
} from '../url'

export interface TGetDocRequestControlFolderLinksInput {
    requestId: number
    orgId: number
    controlId: number
}

export interface TGetDocRequestControlFolderLinksOutput {
    data: FileFolder
}

export function getDocRequestControlFolderLink(inp : TGetDocRequestControlFolderLinksInput) : Promise<TGetDocRequestControlFolderLinksOutput> {
    return axios.get(
        apiv2DocRequestControlFolderLinks(inp.orgId, inp.requestId, inp.controlId),
        getAPIRequestConfig()
    )
}
