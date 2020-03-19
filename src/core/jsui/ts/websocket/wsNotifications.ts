import { createUserNotificationWebsocket } from '../url'
import { authWebsocket } from './websocket'
import { getCurrentCSRF } from '../csrf'

export function connectUserNotificationWebsocket (
    host : string,
    userId : number) : Promise<WebSocket> {
    const url : string = createUserNotificationWebsocket(host, getCurrentCSRF(), userId)
    return authWebsocket(new WebSocket(url))
}
