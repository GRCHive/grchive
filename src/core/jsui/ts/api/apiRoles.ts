import axios from 'axios'
import * as qs from 'query-string'
import { RoleMetadata, Permissions, FullRole } from '../roles'
import { getAPIRequestConfig } from './apiUtility'
import { getOrgRolesUrl,
         newRoleUrl,
         editRoleUrl,
         deleteRoleUrl,
         getSingleOrgRoleUrl } from '../url'
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

export interface TGetSingleRoleInput {
    orgId: number
    roleId: number
}

export interface TGetSingleRoleOutput {
    data: {
        role: FullRole
        userIds: number[]
    }
}

export function getSingleRole(inp : TGetSingleRoleInput) : Promise<TGetSingleRoleOutput> {
    return axios.get(getSingleOrgRoleUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
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

export interface TEditRoleInput extends TNewRoleInput {
    roleId: number
}

export interface TEditRoleOutput {
    data: FullRole
}

export function editRole(inp : TEditRoleInput) : Promise<TEditRoleOutput> {
    return postFormJson<TEditRoleOutput>(editRoleUrl, inp, getAPIRequestConfig())
}

export interface TDeleteRoleInput {
    roleId: number
    orgId: number
}

export interface TDeleteRoleOutput {
}

export function deleteRole(inp : TDeleteRoleInput) : Promise<TDeleteRoleOutput> {
    return postFormJson<TEditRoleOutput>(deleteRoleUrl, inp, getAPIRequestConfig())
}
