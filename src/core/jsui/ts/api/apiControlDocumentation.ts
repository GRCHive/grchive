import { postFormUrlEncoded, postFormMultipart } from '../http'
import { ControlDocumentationCategory } from '../controls'
import { newControlDocCatUrl, editControlDocCatUrl, deleteControlDocCatUrl, uploadControlDocUrl } from '../url'

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
    data: any
}

export function uploadControlDoc(inp : FormData): Promise<TUploadControlDocOutput> {
    return postFormMultipart<TUploadControlDocOutput>(uploadControlDocUrl, inp)
}
