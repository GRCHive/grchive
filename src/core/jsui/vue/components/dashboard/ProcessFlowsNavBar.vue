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
import { contactUsUrl } from '../../../ts/url'
import { deleteProcessFlow, TDeleteProcessFlowInput, TDeleteProcessFlowOutput } from '../../../ts/api/apiProcessFlow'
import { getAllProcessFlow, TGetAllProcessFlowInput, TGetAllProcessFlowOutput } from '../../../ts/api/apiProcessFlow'
import VueSetup from '../../../ts/vueSetup'
import RenderLayout from '../../../ts/render/renderLayout'
import { getCurrentCSRF } from '../../../ts/csrf'

interface NavLinks {
    icon : string
    url : string
    path : string
    title : string
    disabled: boolean
    action: {
        icon: string,
        fn: () => void
    }
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
                    title: data.Name,
                    action: {
                        icon: "mdi-delete",
                        fn: () => {
                            this.onDeleteProcessFlow(data)
                        }
                    }
                })
            }
            return navLinks
        }
    },
    methods: {
        onDeleteProcessFlow(processFlow : ProcessFlowBasicData) {
            deleteProcessFlow(<TDeleteProcessFlowInput>{
                csrf: getCurrentCSRF(),
                flowId: processFlow.Id
            }).then((resp : TDeleteProcessFlowOutput) => {
                VueSetup.store.commit('deleteProcessFlow', processFlow.Id)
                VueSetup.currentRouter.replace('/')
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
        doRefresh() {
            // In the case where the URL does not specify a process flow, take the currently
            // selected process flow for the nav bar model and change the route path to point to it.
            // In the case where the URL does specify a process flow, make sure the refresh
            // points to it at the end.
            let currentProcessFlowId : number = -1

            if (!!VueSetup.currentRouter.currentRoute.params.flowId) {
                currentProcessFlowId = parseInt(VueSetup.currentRouter.currentRoute.params.flowId)
            }

            getAllProcessFlow(<TGetAllProcessFlowInput>{
                csrf: getCurrentCSRF(),
                requested: currentProcessFlowId,
                //@ts-ignore
                organization: this.$root.orgGroupId
            }).then((resp : ResponseData) => {
                if (resp.data.Flows.length == 0) {
                    // Don't make changes here...it'll cause an exception.
                    return
                }
                VueSetup.store.commit('setProcessFlowBasicData', resp.data.Flows)
                if (currentProcessFlowId != -1) {
                    //@ts-ignore
                    VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                        index: resp.data.RequestedIndex,
                        csrf: getCurrentCSRF()
                    })
                } else {
                    // If there's no path parameter for the route then we should manually replace the route on the router ourselves.
                    VueSetup.currentRouter.replace(this.navLinks[VueSetup.store.state.currentProcessFlowIndex].path)

                    // Force a refresh of the current index so that we also force
                    // the full data to be pulled as well.
                    VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                        index: VueSetup.store.state.currentProcessFlowIndex,
                        csrf: getCurrentCSRF()
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
        save(name : String, id : Number) {
            this.showAddFlow = false;

            // Refresh the list of process flows.
            // Then point ourselves to the most recently created process flow.
            getAllProcessFlow(<TGetAllProcessFlowInput>{
                csrf: getCurrentCSRF(),
                requested: id,
                //@ts-ignore
                organization: this.$root.orgGroupId
            }).then((resp : ResponseData) => {
                VueSetup.store.commit('setProcessFlowBasicData', resp.data.Flows)
                VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                    index: resp.data.RequestedIndex,
                    csrf: getCurrentCSRF()
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
            if (idx != VueSetup.store.state.currentProcessFlowIndex) {
                // Need to explicitly clear render layout to prevent rendering when we
                // switch layouts - not sure how much I like having this here though.
                RenderLayout.store.commit('resetNodeLayout')
                VueSetup.store.dispatch('requestSetCurrentProcessFlowIndex', {
                    index: idx,
                    csrf: getCurrentCSRF()
                })
            }
        }
    },
    mounted() {
        this.doRefresh()
    }
})

</script>
