import { createProcessFlowNodeDisplaySettingsWebsocket } from '../url'
import { authWebsocket } from './websocket'
import { getCurrentCSRF } from '../csrf'

export function connectProcessFlowNodeDisplaySettingsWebsocket(
    host : string,
    flowId : number) : Promise<WebSocket> {
    const url : string = createProcessFlowNodeDisplaySettingsWebsocket(host, getCurrentCSRF(), flowId)
    return authWebsocket(new WebSocket(url))
}
