import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'
import { 
    enableFeatureUrl,
} from '../url'


export interface TEnableFeatureInput {
    orgId: number
    featureId: number
}

export function enableFeature(inp : TEnableFeatureInput) : Promise<void> {
    return postFormJson<void>(enableFeatureUrl, inp, getAPIRequestConfig())
}
