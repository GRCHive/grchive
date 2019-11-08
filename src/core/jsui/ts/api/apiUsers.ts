import axios from 'axios'
import * as qs from 'query-string'
import { createGetAllOrgUsersAPIUrl, createUserProfileEditAPIUrl, requestVerificationEmailUrl } from '../url'
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

export function editUserProfile(userId : number, inp : TEditUserProfileInput) : 
        Promise<TEditUserProfileOutput> {
    return postFormUrlEncoded<TEditUserProfileOutput>(createUserProfileEditAPIUrl(userId), inp, getAPIRequestConfig())
}

export interface TRequestVerificationEmailInput {
    userId: number
}

export interface TRequestVerificationEmailOutput {
}


export function requestResendVerificationEmail(inp : TRequestVerificationEmailInput) : Promise<TRequestVerificationEmailOutput> {
    return postFormUrlEncoded<TRequestVerificationEmailOutput>(requestVerificationEmailUrl, inp, getAPIRequestConfig())
}
