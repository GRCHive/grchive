import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import {
    newDeploymentUrl
} from '../url'
import { FullDeployment } from '../deployments'

export interface TNewDeploymentInput {
    orgId: number
    systemId: number | null
    dbId: number | null
}

export interface TNewDeploymentOutput {
    data : FullDeployment
}

export function newDeployment(inp : TNewDeploymentInput) : Promise<TNewDeploymentOutput> {
    return postFormJson<TNewDeploymentOutput>(newDeploymentUrl, inp, getAPIRequestConfig())
}
