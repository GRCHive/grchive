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
    notificationsPulled: boolean
    allNotifications : NotificationWrapper[]
}

const notificationStoreOptions: StoreOptions<NotificationStoreState> = {
    state: {
        notificationsPulled: false,
        allNotifications: [],
    },
    mutations: {
        startNotificationPull(state) {
            state.notificationsPulled = true
        },
        setNotifications(state, data) {
            state.allNotifications = data
        },
        markAllAsRead(state) {
            state.allNotifications.forEach((ele : NotificationWrapper) => {
                ele.Read = true
            })
        }
    },
    actions: {
        pullNotifications(context) {
            context.commit('startNotificationPull')
            allNotifications({
                userId: PageParamsStore.state.user!.Id,
            }).then((resp : TAllNotificationOutput) => {
                context.commit('setNotifications', resp.data)
            }).catch((err : any) => {
                console.log(err)
            })
        }
    },
    getters: {
        hasUnreadNotifications(state) : boolean {
            return state.allNotifications.some((ele : NotificationWrapper) => !ele.Read)
        }
    }
}

let store = new Vuex.Store<NotificationStoreState>(notificationStoreOptions)
export default store
