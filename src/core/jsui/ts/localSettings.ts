import Vuex, { StoreOptions } from 'vuex'

interface LocalSettingsStoreState {
    miniNavBar: boolean
    showHideAttributeEditor: boolean
}

const MiniNavBarLocalStorageKey : string = "miniNavBar"
const ShowHideAttributeEditorLocalStorageKey : string = "showHideAttributeEditor"

const localSettingStore: StoreOptions<LocalSettingsStoreState> = {
    state: {
        miniNavBar : false,
        showHideAttributeEditor: false
    },
    mutations: {
        setMiniNavBar(state, val) {
            state.miniNavBar = val
            window.localStorage.setItem(MiniNavBarLocalStorageKey,
                val ? "true" : "false")
        },
        setShowHideAttributeEditor(state, val) {
            state.showHideAttributeEditor = val
            window.localStorage.setItem(ShowHideAttributeEditorLocalStorageKey,
                val ? "true" : "false")
        }
    },
}

let store = new Vuex.Store<LocalSettingsStoreState>(localSettingStore)

// Initialize store from local storage
let miniNavBar = window.localStorage.getItem(MiniNavBarLocalStorageKey)
if (miniNavBar != null) {
    store.commit('setMiniNavBar', miniNavBar == "true")
}

let showHideAttrEditor = window.localStorage.getItem(ShowHideAttributeEditorLocalStorageKey)
if (showHideAttrEditor != null) {
    store.commit('setShowHideAttributeEditor', showHideAttrEditor == "true")
}

export default store
