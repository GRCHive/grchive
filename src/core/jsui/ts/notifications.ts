import Vuex, { StoreOptions } from 'vuex'
import { allNotifications, TAllNotificationOutput } from './api/apiNotifications'
import { PageParamsStore } from './pageParams'

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
    canPullMore: boolean
    requestInProgress: boolean
}

const notificationStoreOptions: StoreOptions<NotificationStoreState> = {
    state: {
        allNotifications: [],
        canPullMore: true,
        requestInProgress : false
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
        markAllAsRead(state) {
            state.allNotifications.forEach((ele : NotificationWrapper) => {
                ele.Read = true
            })
        },
        stopAllowingPull(state) {
            state.canPullMore = false
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
        }
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
