import axios from 'axios'
import * as qs from 'query-string'
import {
    allNotificationUrl
} from '../url'
import {
    NotificationWrapper,
    cleanJsonNotificationWrapper,
} from '../notifications'
import { getAPIRequestConfig } from './apiUtility'

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
