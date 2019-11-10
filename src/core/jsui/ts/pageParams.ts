import Vuex, { StoreOptions } from 'vuex'

export interface PageParamsStoreState {
    organization: {
        Url: string
        Name: string
        OktaGroupName: string
        Id: number
    } | null
    user: {
        Id: number
        FirstName: string
        LastName: string
        Email: string
        Auth: boolean
        Verified: boolean
    } | null
    site: {
        CompanyName: string
        Domain: string
        Host: string
    } | null,
    auth: {
        OktaServer: string
        OktaClientId: string
        OktaRedirectUri: string
        OktaScope: string
    } | null
}

const storeOptions: StoreOptions<PageParamsStoreState> = {
    state: {
        organization: null,
        user: null,
        site: null,
        auth: null
    },
    mutations: {
        replaceState(state, data : PageParamsStoreState) {
            state.organization = data.organization
            state.user = data.user
            state.site = data.site
            state.auth = data.auth
        }
    }
}

export let PageParamsStore = new Vuex.Store<PageParamsStoreState>(storeOptions)
