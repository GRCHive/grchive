import axios from 'axios'
import * as qs from 'query-string'
import { 
    allControlDocCatLinkUrl
} from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { ControlDocumentationCategory } from '../controls'

export interface TAllControlDocCatLinkInput {
    controlId?: number
    catId?: number
    orgId: number
}

export interface TAllControlDocCatLinkOutput {
    data: {
        Cats?: ControlDocumentationCategory[]
        Controls?: ProcessFlowControl[]
    }
}

export function allControlDocCatLink(inp : TAllControlDocCatLinkInput) : Promise<TAllControlDocCatLinkOutput> {
    return axios.get(allControlDocCatLinkUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
