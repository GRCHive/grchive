import axios from 'axios'
import * as qs from 'query-string'
import {
    allNotificationUrl,
    readNotificationUrl,
} from '../url'
import {
    NotificationWrapper,
    cleanJsonNotificationWrapper,
} from '../notifications'
import { getAPIRequestConfig } from './apiUtility'
import { postFormJson } from '../http'

export interface TAllNotificationInput {
    userId: number
}

export interface TAllNotificationOutput {
    data: NotificationWrapper[]
}

export function allNotifications(inp : TAllNotificationInput) : Promise<TAllNotificationOutput> {
    return axios.get(allNotificationUrl + '?' + qs.stringify(inp), getAPIRequestConfig()).then((resp : TAllNotificationOutput) => {
        resp.data.forEach(cleanJsonNotificationWrapper)
        return resp
    })
}

export interface TMarkNotificationReadInput {
    userId: number
    notificationIds: number[]
    all: boolean
}

export function markNotificationRead(inp : TMarkNotificationReadInput) : Promise<void> {
    return postFormJson<void>(readNotificationUrl, inp, getAPIRequestConfig())
}
