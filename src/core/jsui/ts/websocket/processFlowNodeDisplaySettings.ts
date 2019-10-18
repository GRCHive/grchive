import { createProcessFlowNodeDisplaySettingsWebsocket } from '../url'

export function connectProcessFlowNodeDisplaySettingsWebsocket(
    host : string,
    csrf : string,
    flowId : number) : WebSocket {
    const url : string = createProcessFlowNodeDisplaySettingsWebsocket(host, csrf, flowId)
    return new WebSocket(url)
}
