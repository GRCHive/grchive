import axios from 'axios'
import * as qs from 'query-string'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { Comment } from '../comments'
import { 
    newCommentUrl,
    allCommentUrl,
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

export interface TNewCommentInput {
    comment: TGenericNewCommentInput
    requestId?: number
    sqlRequestId?: number
    catId?: number
    fileId?: number
    orgId: number
}

function cleanComment(comment : Comment) : Comment {
    comment.PostTime = new Date(comment.PostTime)
    return comment
}

export function newComment(inp : TNewCommentInput) : Promise<TNewCommentOutput> {
    return postFormJson<TNewCommentOutput>(newCommentUrl, inp, getAPIRequestConfig()).then((resp : TNewCommentOutput) => {
        resp.data = cleanComment(resp.data)
        return resp
    })
}

export interface TAllCommentsInput {
    requestId?: number
    fileId?: number
    orgId: number
}

export function allComments(inp : TAllCommentsInput) : Promise<TGetAllCommentsOutput> {
    return axios.get(allCommentUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TGetAllCommentsOutput) => {
        resp.data = resp.data.map(cleanComment)
        return resp
    })
}
