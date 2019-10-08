<template>
    <generic-nav-bar
        :selectedPage="selectedNavIndex"
        :navLinks="navLinks"
        :mini="false"
        :style="xOffsetStyle"
        primary-color="accent"
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

interface NavLinks {
    icon : string
    url : string
    title : string
    disabled: boolean
}

interface ProcessFlowData {
    Name : string
    Description: string
    CreationTime: Date
    LastUpdatedTime: Date
}

interface ResponseData {
    data: {
        Flows: ProcessFlowData[],
        RequestedIndex: number
    }
}

export default Vue.extend({
    components: {
        GenericNavBar,
        CreateNewProcessFlowForm
    },
    data : () => ({
        navLinks: [] as NavLinks[],
        showAddFlow: false,
        selectedNavIndex : 0
    }),
    computed: {
        xOffsetStyle() {
            return {
                // TODO: This needs to change if the width of the initial nav changes
                //@ts-ignore
                "transform": "translateX(" + VueSetup.store.state.primaryNavBarWidth + "px)"
            }
        }
    },
    methods: {
        createNavLinksFromProcessFlowData(data : ProcessFlowData[]) {
            this.navLinks = [] as NavLinks[]
            for (let i = 0; i < data.length; ++i) {
                this.navLinks.push({
                    icon: '',    
                    disabled: false,
                    url: '',
                    title: data[i].Name
                })
            }
        },
        doRefresh() {
            this.refresh(-1).then((resp : ResponseData) => {
                this.createNavLinksFromProcessFlowData(resp.data.Flows)
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
                this.createNavLinksFromProcessFlowData(resp.data.Flows)
                this.selectedNavIndex = resp.data.RequestedIndex;
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }
    },
    mounted() {
        this.doRefresh()
    }
})

</script>
