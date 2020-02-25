import { PageParamsStore } from '../pageParams'
import { getApiKey } from '../api/apiUtility'

function waitForWebsocket(ws : WebSocket) : Promise<void> {
    return new Promise<void>((resolve, reject) => {
        const timeoutSec : number = 15
        let elapsedSec: number = 0
        let intervalSec : number = 0.1
        let id : any = setInterval(() => {
            if (ws.readyState == WebSocket.OPEN) {
                clearInterval(id)
                resolve()
            } else if (elapsedSec >= timeoutSec) {
                clearInterval(id)
                reject()
            }
            elapsedSec += intervalSec
        }, intervalSec * 1000.0)
    })
}

export async function authWebsocket(ws : WebSocket) : Promise<WebSocket> {
    let data = {
        ApiKey: getApiKey(),
        OrgId: PageParamsStore.state.organization!.Id,
    }

    return new Promise<WebSocket>((resolve, reject) => {
        waitForWebsocket(ws).then(() => {
            try {
                ws.send(JSON.stringify(data))
                resolve(ws)
            } catch (err) {
                reject(err)
            }
        }).catch((err : any) => {
            reject(err)
        })

    })
}
