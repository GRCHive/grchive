import axios from 'axios'
import * as qs from 'query-string'
import { RoleMetadata, Permissions, FullRole } from '../roles'
import { getAPIRequestConfig } from './apiUtility'
import { getOrgRolesUrl, newRoleUrl } from '../url'
import { postFormJson } from '../http'

export interface TGetAllOrgRolesInput {
    orgId: number
}

export interface TGetAllOrgRolesOutput {
    data: RoleMetadata[]
}

export function getAllOrgRoles(inp : TGetAllOrgRolesInput) : Promise<TGetAllOrgRolesOutput> {
    return axios.get(getOrgRolesUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TNewRoleInput {
    orgId: number
    name: string
    description: string
    permissions: Permissions
}

export interface TNewRoleOutput {
    data: FullRole
}

export function newRole(inp : TNewRoleInput) : Promise<TNewRoleOutput> {
    return postFormJson<TNewRoleOutput>(newRoleUrl, inp, getAPIRequestConfig())
}
