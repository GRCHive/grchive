import axios from 'axios'
import { getAPIRequestConfig } from './apiUtility'
import {
    codeParamTypeMetadataUrl
} from '../url'
import { SupportedParamType } from '../code'

export interface TCodeParameterTypeMetadataOutput {
    data: SupportedParamType[]
}

export function getCodeParameterTypeMetadata() : Promise<TCodeParameterTypeMetadataOutput> {
    return axios.get(codeParamTypeMetadataUrl, getAPIRequestConfig())
}
