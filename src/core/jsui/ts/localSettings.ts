import Vuex, { StoreOptions } from 'vuex'

interface LocalSettingsStoreState {
    miniNavBar: boolean
    showHideAttributeEditor: boolean
    viewBoxTransform: TransformData
    viewBoxZoom: number
}

const MiniNavBarLocalStorageKey : string = "miniNavBar"
const ShowHideAttributeEditorLocalStorageKey : string = "showHideAttributeEditor"
const ViewBoxTransformLocalStorageKey : string = "viewBoxTransform"
const ViewBoxZoomLocalStorageKey : string = "viewBoxZoom"

const localSettingStore: StoreOptions<LocalSettingsStoreState> = {
    state: {
        miniNavBar : false,
        showHideAttributeEditor: false,
        viewBoxTransform: {
            tx: 0,
            ty: 0
        },
        viewBoxZoom: 1.0
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
        },
        setViewBoxTransform(state, val) {
            state.viewBoxTransform = val
            window.localStorage.setItem(ViewBoxTransformLocalStorageKey, JSON.stringify(val))
        },
        setViewBoxZoom(state, val) {
            state.viewBoxZoom = val
            window.localStorage.setItem(ViewBoxZoomLocalStorageKey, val)
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

let viewBoxTransform = window.localStorage.getItem(ViewBoxTransformLocalStorageKey)
if (viewBoxTransform != null) {
    store.commit('setViewBoxTransform', JSON.parse(viewBoxTransform))
}

let viewBoxZoom = window.localStorage.getItem(ViewBoxZoomLocalStorageKey)
if (viewBoxZoom != null) {
    store.commit('setViewBoxZoom', Number(viewBoxZoom))
}

export default store
