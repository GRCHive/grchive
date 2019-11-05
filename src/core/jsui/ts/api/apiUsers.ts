import axios from 'axios'
import * as qs from 'query-string'
import { createGetAllOrgUsersAPIUrl, createUserProfileEditAPIUrl } from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormUrlEncoded } from '../http'

export interface TGetAllOrgUsersInput {
    org : string
}

export interface TGetAllOrgUsersOutput {
    data: User[]
}

export function getAllOrgUsers(inp : TGetAllOrgUsersInput) : Promise<TGetAllOrgUsersOutput> {
    return axios.get(createGetAllOrgUsersAPIUrl(inp.org) + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TEditUserProfileInput {
    firstName : string
    lastName : string
}

export interface TEditUserProfileOutput {
}

export function editUserProfile(email : string, inp : TEditUserProfileInput) : 
        Promise<TEditUserProfileOutput> {
    return postFormUrlEncoded<TEditUserProfileOutput>(createUserProfileEditAPIUrl(email), inp, getAPIRequestConfig())
}
