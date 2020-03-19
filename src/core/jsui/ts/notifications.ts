import Vuex, { StoreOptions } from 'vuex'
import { allNotifications, TAllNotificationOutput } from './api/apiNotifications'
import { PageParamsStore } from './pageParams'
import { connectUserNotificationWebsocket } from './websocket/wsNotifications'

export interface Notification {
    Id                  : number
    OrgId               : number
    Time                : Date
    SubjectType         : string
    SubjectId           : number
    Verb                : string
    ObjectType          : string
    ObjectId            : number
    IndirectObjectType  : string
    IndirectObjectId    : number
}

export interface NotificationWrapper {
    Notification: Notification
    OrgName: string
    Read: boolean
}

export function cleanJsonNotificationWrapper(n : NotificationWrapper) {
    n.Notification.Time = new Date(n.Notification.Time)
}

interface NotificationStoreState {
    allNotifications : NotificationWrapper[]
    recentNotifications : NotificationWrapper[]
    canPullMore: boolean
    requestInProgress: boolean
    wsConnected: boolean
}

let websocketConnection : WebSocket

function connectWebsocket(context : any, host : string) {
    if (context.state.wsConnected) {
        return
    }
    context.commit('connectWs')

    if (!!websocketConnection) {
        websocketConnection.close()
    }

    connectUserNotificationWebsocket(host, PageParamsStore.state.user!.Id).then(
        (ws : WebSocket) => {
            websocketConnection = ws

            websocketConnection.onclose = (e : CloseEvent) => {
                if (e.code != 1001) {
                    // Automatically try to reconnect?
                    connectWebsocket(context, host)
                }
            }

            websocketConnection.onmessage = (e : MessageEvent) => {
                let data : NotificationWrapper = JSON.parse(e.data)
                cleanJsonNotificationWrapper(data)

                context.dispatch('pushNotification', data)
            }
        })
}

const browserNotificationTimeMs = 5000

const notificationStoreOptions: StoreOptions<NotificationStoreState> = {
    state: {
        allNotifications: [],
        recentNotifications: [],
        canPullMore: true,
        requestInProgress : false,
        wsConnected: false
    },
    mutations: {
        startPull(state) {
            state.requestInProgress = true
        },
        stopPull(state) {
            state.requestInProgress = false
        },
        addNotifications(state, data) {
            state.allNotifications.push(...data)
        },
        pushNotification(state, notif) {
            state.allNotifications.unshift(notif)
            state.recentNotifications.unshift(notif)
        },
        removeRecentNotification(state, id) {
            let idx : number = state.recentNotifications.findIndex(
                    (ele : NotificationWrapper) => ele.Notification.Id == id
            )

            if (idx == -1) {
                return
            }
            state.recentNotifications.splice(idx, 1)
        },
        markAllAsRead(state) {
            state.allNotifications.forEach((ele : NotificationWrapper) => {
                ele.Read = true
            })
        },
        stopAllowingPull(state) {
            state.canPullMore = false
        },
        connectWs(state) {
            state.wsConnected = true
        }
    },
    actions: {
        pullNotifications(context) {
            if (!context.state.canPullMore || context.state.requestInProgress) {
                return
            }

            context.commit('startPull')
            allNotifications({
                userId: PageParamsStore.state.user!.Id,
                offset: context.state.allNotifications.length,
            }).then((resp : TAllNotificationOutput) => {
                context.commit(
                    'addNotifications',
                    resp.data.filter((ele : NotificationWrapper) => !context.getters.notificationIdSet.has(ele.Notification.Id))
                )

                if (resp.data.length == 0) {
                    context.commit('stopAllowingPull')
                }

                context.commit('stopPull')
            }).catch((err : any) => {
                console.log(err)
                context.commit('stopPull')
            })
        },
        pushNotification(context, notif) {
            context.commit('pushNotification', notif)
            setTimeout(() => {
                context.commit('removeRecentNotification', notif.Notification.Id)
            }, browserNotificationTimeMs)
        },
        initialize(context, {host}) {
            connectWebsocket(context, host)
        },
    },
    getters: {
        hasUnreadNotifications(state) : boolean {
            return state.allNotifications.some((ele : NotificationWrapper) => !ele.Read)
        },
        notificationIdSet(state): Set<number> {
            return new Set<number>(state.allNotifications.map((ele : NotificationWrapper) => ele.Notification.Id))
        }
    }
}

let store = new Vuex.Store<NotificationStoreState>(notificationStoreOptions)
export default store
