import axios from 'axios'
import * as qs from 'query-string'
import { createGetAllOrgUsersAPIUrl,
         createUserProfileEditAPIUrl,
         createUserGetOrgsAPIUrl,
         requestVerificationEmailUrl,
         inviteUsersToOrgUrl } from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { postFormUrlEncoded } from '../http'
import { Organization } from '../organizations'

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

export interface TGetUserOrgsInput {
    userId: number
}

export interface TGetUserOrgsOutput {
    data: Organization[]
}

export function getAllOrgsForUser(inp: TGetUserOrgsInput) : Promise<TGetUserOrgsOutput> {
    return axios.get(createUserGetOrgsAPIUrl(inp.userId) + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TInviteUsersInput {
    fromUserId: number
    fromOrgId: number
    toEmails: string[]
}

export interface TInviteUsersOutput {
}

export function inviteUsers(inp: TInviteUsersInput) : Promise<TInviteUsersOutput> {
    return postFormUrlEncoded<TInviteUsersOutput>(inviteUsersToOrgUrl, inp, getAPIRequestConfig())
}
