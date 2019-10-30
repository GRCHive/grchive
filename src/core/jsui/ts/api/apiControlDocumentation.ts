import { postFormUrlEncoded } from '../http'
import { ControlDocumentationCategory } from '../controls'
import { newControlDocCatUrl } from '../url'

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
