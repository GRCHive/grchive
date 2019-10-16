<template>
    <generic-nav-bar
        :selectedPage="selectedNavIndex"
        :navLinks="navLinks"
        :mini="false"
        :style="xOffsetStyle"
        primary-color="accent"
        @item-change="onItemClick"
    >
        <v-list-item>
            <v-list-item-title>
                Process Flows
            </v-list-item-title>
            <v-list-item-action>
                <v-btn icon class="mx-3" @click="doRefresh">
                    <v-icon>mdi-refresh</v-icon>
                </v-btn>
            </v-list-item-action>
            <v-list-item-action>
                <v-dialog v-model="showAddFlow" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn icon class="mx-3" v-on="on">
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>
                    <create-new-process-flow-form
                        v-on:do-save="save"
                        v-on:do-cancel="showAddFlow = false"
                    >
                    </create-new-process-flow-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>
    </generic-nav-bar>
</template>

<script lang="ts">

import Vue from 'vue'
import GenericNavBar from '../GenericNavBar.vue'
import CreateNewProcessFlowForm from './CreateNewProcessFlowForm.vue'
import axios from 'axios'
import * as qs from 'query-string'
import { contactUsUrl, getAllProcessFlowAPIUrl} from '../../../ts/url'
import VueSetup from '../../../ts/vueSetup'
import RenderLayout from '../../../ts/render/renderLayout'

interface NavLinks {
    icon : string
    url : string
    path : string
    title : string
    disabled: boolean
}

interface ResponseData {
    data: {
        Flows: ProcessFlowBasicData[],
        RequestedIndex: number
    }
}

export default Vue.extend({
    components: {
        GenericNavBar,
        CreateNewProcessFlowForm
    },
    data : () => ({
        showAddFlow: false
    }),
    computed: {
        xOffsetStyle() {
            return {
                // TODO: This needs to change if the width of the initial nav changes
                //@ts-ignore
                "transform": "translateX(" + VueSetup.store.state.primaryNavBarWidth + "px)"
            }
        },
        selectedNavIndex() : Number {
            return VueSetup.store.state.currentProcessFlowIndex
        },
        allBasicData() : ProcessFlowBasicData[] {
            return VueSetup.store.state.allProcessFlowBasicData
        },
        navLinks() : NavLinks[] {
            let navLinks : NavLinks[] = []
            for (let data of this.allBasicData) {
                navLinks.push({
                    icon: '',    
                    disabled: false,
                    url: '',
                    path: data.Id.toString(),
                    title: data.Name
                })
            }
            return navLinks
        }
    },
    methods: {
        doRefresh() {
            // In the case where the URL does not specify a process flow, take the currently
            // selected process flow for the nav bar model and change the route path to point to it.
            // In the case where the URL does specify a process flow, make sure the refresh
            // points to it at the end.
            let currentProcessFlowId : number = -1

            if (!!VueSetup.currentRouter.currentRoute.params.flowId) {
                currentProcessFlowId = parseInt(VueSetup.currentRouter.currentRoute.params.flowId)
            }

            this.refresh(currentProcessFlowId).then((resp : ResponseData) => {
                if (resp.data.Flows.length == 0) {
                    // Don't make changes here...it'll cause an exception.
                    return
                }
                VueSetup.store.commit('setProcessFlowBasicData', resp.data.Flows)
                if (currentProcessFlowId != -1) {
                    //@ts-ignore
                    VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                        index: resp.data.RequestedIndex,
                        //@ts-ignore
                        csrf: this.$root.csrf
                    })
                } else {
                    // If there's no path parameter for the route then we should manually replace the route on the router ourselves.
                    VueSetup.currentRouter.replace(this.navLinks[VueSetup.store.state.currentProcessFlowIndex].path)

                    // Force a refresh of the current index so that we also force
                    // the full data to be pulled as well.
                    VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                        index: VueSetup.store.state.currentProcessFlowIndex,
                        //@ts-ignore
                        csrf: this.$root.csrf
                    })
                }
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        refresh(id : Number) : Promise<ResponseData> {
            // id = -1 means that we don't need the index of the passed in 
            // process flow id in the results (save us some javascript computation time).
            let passedData = Object()
            if (id >= 0) {
                passedData['requested'] = id
            }
            //@ts-ignore
            passedData['csrf'] = this.$root.csrf
            //@ts-ignore
            passedData['organization'] = this.$root.orgGroupId

            return new Promise<ResponseData>(function(resolve, reject){
                axios.get(getAllProcessFlowAPIUrl + '?' + qs.stringify(passedData)).then((resp : ResponseData) => {
                    // 
                    resolve(resp)
                }).catch((err) => {
                    //
                    reject(err)
                })
            })
        },
        save(name : String, id : Number) {
            this.showAddFlow = false;

            // Refresh the list of process flows.
            // Then point ourselves to the most recently created process flow.
            this.refresh(id).then((resp : ResponseData) => {
                VueSetup.store.commit('setProcessFlowBasicData', resp.data.Flows)
                VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                    index: resp.data.RequestedIndex,
                    //@ts-ignore
                    csrf: this.$root.csrf
                })
                VueSetup.currentRouter.replace(this.navLinks[VueSetup.store.state.currentProcessFlowIndex].path)
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        onItemClick(_ : MouseEvent, idx : number) {
            // Need to explicitly clear render layout to prevent rendering when we
            // switch layouts - not sure how much I like having this here though.
            RenderLayout.store.commit('resetNodeLayout')
            VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                index: idx,
                //@ts-ignore
                csrf: this.$root.csrf
            })
        }
    },
    mounted() {
        this.doRefresh()
    }
})

</script>
