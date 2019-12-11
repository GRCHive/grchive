import axios from 'axios'
import * as qs from 'query-string'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { Comment } from '../comments'
import { 
    newDocRequestCommentUrl,
    allDocRequestCommentUrl,
} from '../url'

export interface TGenericNewCommentInput {
    userId: number
    content: string
}

export interface TNewCommentOutput {
    data: Comment
}

export interface TGetAllCommentsOutput {
    data: Comment[]
}

export interface TDocRequestNewCommentInput {
    comment: TGenericNewCommentInput
    requestId: number
    catId: number
    orgId: number
}

function cleanComment(comment : Comment) : Comment {
    comment.PostTime = new Date(comment.PostTime)
    return comment
}

export function newDocumentRequestComment(inp : TDocRequestNewCommentInput) : Promise<TNewCommentOutput> {
    return postFormJson<TNewCommentOutput>(newDocRequestCommentUrl, inp, getAPIRequestConfig()).then((resp : TNewCommentOutput) => {
        resp.data = cleanComment(resp.data)
        return resp
    })
}

export interface TDocRequestAllCommentsInput {
    requestId: number
    orgId: number
}

export function allDocumentRequestComments(inp : TDocRequestAllCommentsInput) : Promise<TGetAllCommentsOutput> {
    return axios.get(allDocRequestCommentUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetAllCommentsOutput) => {
        resp.data = resp.data.map(cleanComment)
        return resp
    })
}
