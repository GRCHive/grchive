import Vue from 'vue'
import { Store, StoreOptions } from 'vuex'
import { DbRefresh } from '../sql'
import { PageParamsStore } from '../pageParams'
import {
    allSqlRefresh, TAllSqlRefreshOutput,
    getSqlRefresh, TGetSqlRefreshOutput,
} from '../api/apiSqlRefresh'

interface DatabaseStoreState {
    allRefreshes: DbRefresh[] | null
    allRefreshInProgress: boolean 

    selectedRefresh : DbRefresh | null
    isPollingRefresh : boolean
}

const refreshIntervalSeconds : number = 15
const databaseStore : StoreOptions<DatabaseStoreState> = {
    state: {
        allRefreshes: null,
        allRefreshInProgress : false,
        selectedRefresh: null,
        isPollingRefresh: false
    },
    mutations: {
        startAllRefresh(state : DatabaseStoreState) {
            state.allRefreshInProgress = true
        },
        setAllRefresh(state : DatabaseStoreState, data : DbRefresh[] | null) {
            state.allRefreshes = data
            state.allRefreshInProgress = false
        },
        replaceRefreshIndex(state : DatabaseStoreState, { data, idx }) {
            Vue.set(state.allRefreshes!, idx, data)
        },
        removeRefreshIndex( state : DatabaseStoreState, idx : number ) {
            state.allRefreshes!.splice(idx , 1)
        },
        addRefresh(state : DatabaseStoreState, refresh : DbRefresh) {
            state.allRefreshes!.unshift(refresh)
        },
        startPollingRefresh(state : DatabaseStoreState) {
            state.isPollingRefresh = true
        },
        stopPollingRefresh(state : DatabaseStoreState) {
            state.isPollingRefresh = false
        },
        selectNewRefresh(state : DatabaseStoreState, ref : DbRefresh | null) {
            state.selectedRefresh = ref
            state.isPollingRefresh = false
        },
    },
    actions: {
        startRefreshPoll(context) {
            if (!context.state.selectedRefresh) {
                return
            }

            context.commit('startPollingRefresh')
            let refreshId : number = context.state.selectedRefresh!.Id
            // Silently ignore any polling errors.
            let intervalId = setInterval(() => {
                let idx = context.state.allRefreshes!.findIndex((ele : DbRefresh) => ele.Id == refreshId)
                if (idx == -1) {
                    clearInterval(intervalId)
                    return
                }

                getSqlRefresh({
                    refreshId : refreshId,
                    orgId: PageParamsStore.state.organization!.Id,
                }).then((resp : TGetSqlRefreshOutput) => {
                    if (!resp.data.RefreshFinishTime) {
                        return
                    }

                    let idx = context.state.allRefreshes!.findIndex((ele : DbRefresh) => ele.Id == refreshId)
                    if (idx == -1) {
                        clearInterval(intervalId)
                        return
                    }

                    context.commit('replaceRefreshIndex', { data: resp.data, idx : idx})
                    if (!!context.state.selectedRefresh && context.state.selectedRefresh.Id == refreshId) {
                        context.commit('selectNewRefresh', resp.data)
                    }

                    clearInterval(intervalId)
                })
            }, refreshIntervalSeconds * 1000)

        },
        requestSetNewRefresh(context, ref : DbRefresh | null) {
            context.commit('selectNewRefresh', ref)
            if (!!context.state.selectedRefresh && !context.state.selectedRefresh.RefreshFinishTime) {
                context.dispatch('startRefreshPoll')
            }
        },
        getRefreshesForDb(context, id : number) {
            if (context.state.allRefreshInProgress) {
                return
            }
            context.commit('startAllRefresh')

            allSqlRefresh({
                dbId: id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TAllSqlRefreshOutput) => {
                context.commit('setAllRefresh', resp.data)
                if (resp.data.length > 0) {
                    context.dispatch('requestSetNewRefresh', resp.data[0])
                }
            })
        },
        deleteRefresh(context, ref : DbRefresh) {
            let idx = context.state.allRefreshes!.findIndex((ele : DbRefresh) => ele.Id == ref.Id)
            if (idx == -1) {
                return
            }

            if (!!context.state.selectedRefresh && ref.Id == context.state.selectedRefresh.Id &&
                    !context.state.selectedRefresh.RefreshFinishTime) {
                context.commit('stopPollingRefresh')
            }

            context.commit('removeRefreshIndex', idx)
            context.dispatch('requestSetNewRefresh', 
                context.state.allRefreshes!.length > 0 ? 
                    context.state.allRefreshes![0] :
                    null)
        }
    }
}

let allDatabaseStores : Record<number, Store<DatabaseStoreState>> = Object()

export type DatabaseStore = Store<DatabaseStoreState>
export function getStoreForDatabase(id : number) : DatabaseStore {
    if (!(id in allDatabaseStores)) {
        let store = new Store<DatabaseStoreState>(databaseStore)
        store.dispatch('getRefreshesForDb', id)
        allDatabaseStores[id] = store
    }

    return allDatabaseStores[id]
}
