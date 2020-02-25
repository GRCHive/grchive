import { createProcessFlowNodeDisplaySettingsWebsocket } from '../url'
import { authWebsocket } from './websocket'

export function connectProcessFlowNodeDisplaySettingsWebsocket(
    host : string,
    csrf : string,
    flowId : number) : Promise<WebSocket> {
    const url : string = createProcessFlowNodeDisplaySettingsWebsocket(host, csrf, flowId)
    return authWebsocket(new WebSocket(url))
}
