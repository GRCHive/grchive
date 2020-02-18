import axios from 'axios'
import * as qs from 'query-string'
import { postFormJson } from '../http'
import { getAPIRequestConfig } from './apiUtility'
import { Comment } from '../comments'
import { 
    newCommentUrl,
    allCommentUrl,
    updateCommentUrl,
    deleteCommentUrl,
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
    sqlRequestId?: number
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

export interface TUpdateCommentInput {
    commentId: number
    content: string
}

export interface TUpdateCommentOutput {
    data: Comment
}

export function updateComment(inp : TUpdateCommentInput) : Promise<TUpdateCommentOutput> {
    return postFormJson<TUpdateCommentOutput>(updateCommentUrl, inp, getAPIRequestConfig()).then((resp : TUpdateCommentOutput) => {
        resp.data = cleanComment(resp.data)
        return resp
    })
}

export interface TDeleteCommentInput {
    commentId: number
}

export function deleteComment(inp : TDeleteCommentInput) : Promise<void> {
    return postFormJson<void>(deleteCommentUrl, inp, getAPIRequestConfig())
}
