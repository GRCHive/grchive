import axios from 'axios'
import * as qs from 'query-string'
import { getGLUrl,
         createNewGLAccUrl,
         createNewGLCatUrl,
         editGLCatUrl,
         deleteGLCatUrl,
         getGLAccUrl } from '../url'
import { getAPIRequestConfig } from './apiUtility'
import { RawGeneralLedgerCategory, RawGeneralLedgerAccount } from '../generalLedger'
import { postFormJson } from '../http'

export interface TGetGLInputs {
    orgId: number
}

export interface TGetGLOutputs {
    data: {
        Categories: RawGeneralLedgerCategory[]
        Accounts: RawGeneralLedgerAccount[]
    }
}

export function getGL(inp : TGetGLInputs) : Promise<TGetGLOutputs> {
    return axios.get(getGLUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}

export interface TNewGLCategoryInputs {
    orgId:            number
	parentCategoryId: number | null
	name:             string
	description:      string
}

export interface TNewGLCategoryOutputs {
    data: RawGeneralLedgerCategory
}

export function newGLCategory(inp : TNewGLCategoryInputs) : Promise<TNewGLCategoryOutputs> {
    return postFormJson<TNewGLCategoryOutputs>(createNewGLCatUrl, inp, getAPIRequestConfig())
}

export interface TEditGLCategoryInputs extends TNewGLCategoryInputs {
    catId: number
}

export interface TEditGLCategoryOutputs {
    data: RawGeneralLedgerCategory
}

export function editGLCategory(inp : TEditGLCategoryInputs) : Promise<TEditGLCategoryOutputs> {
    return postFormJson<TEditGLCategoryOutputs>(editGLCatUrl, inp, getAPIRequestConfig())
}

export interface TNewGLAccountInputs {
	orgId:               number
	parentCategoryId:    number
	accountId:           string
	accountName:         string
	accountDescription:  string
	financiallyRelevant: boolean
}

export interface TNewGLAccountOutputs {
    data: RawGeneralLedgerAccount
}

export function newGLAccount(inp : TNewGLAccountInputs) : Promise<TNewGLAccountOutputs> {
    return postFormJson<TNewGLAccountOutputs>(createNewGLAccUrl, inp, getAPIRequestConfig())
}

export interface TDeleteGLCategoryInputs {
    orgId: number
    catId: number
}

export interface TDeleteGLCategoryOutputs {
    data: RawGeneralLedgerCategory
}

export function deleteGLCategory(inp : TDeleteGLCategoryInputs) : Promise<TDeleteGLCategoryOutputs> {
    return postFormJson<TDeleteGLCategoryOutputs>(deleteGLCatUrl, inp, getAPIRequestConfig())
}

export interface TGetGLAccountInputs {
	orgId: number
	accId: number
}

export interface TGetGLAccountOutputs {
    data: {
        Account: RawGeneralLedgerAccount
        Parents: RawGeneralLedgerCategory[]
    }
}

export function getGLAccount(inp : TGetGLAccountInputs) : Promise<TGetGLAccountOutputs> {
    return axios.get(getGLAccUrl + '?' + qs.stringify(inp), getAPIRequestConfig())
}
