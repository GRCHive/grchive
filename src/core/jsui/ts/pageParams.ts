import Vuex, { StoreOptions } from 'vuex'

export interface PageParamsStoreState {
    organization: {
        Url: string
        Name: string
        OktaGroupName: string
        Id: number
    } | null
    user: {
        FirstName: string
        LastName: string
        Email: string
        Auth: boolean
    } | null
    site: {
        CompanyName: string
        Domain: string
        Host: string
    } | null
}

const storeOptions: StoreOptions<PageParamsStoreState> = {
    state: {
        organization: null,
        user: null,
        site: null,
    },
    mutations: {
        replaceState(state, data : PageParamsStoreState) {
            state.organization = data.organization
            state.user = data.user
            state.site = data.site
        }
    }
}

export let PageParamsStore = new Vuex.Store<PageParamsStoreState>(storeOptions)
