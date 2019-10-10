<template>
    <v-toolbar flat height="30px">
        <v-toolbar-items>
            <v-divider vertical></v-divider>
            <v-btn text color="accent">
                File
                <v-icon small color="accent">mdi-chevron-down</v-icon>
            </v-btn>
            <v-divider vertical></v-divider>
            <v-btn text color="accent">
                View
                <v-icon small color="accent">mdi-chevron-down</v-icon>
            </v-btn>
            <v-divider vertical></v-divider>
        </v-toolbar-items>
        <div class="flex-grow-1"></div>
        <v-toolbar-items>
            <v-divider vertical></v-divider>
            <v-menu offset-y>
                <template v-slot:activator="{ on }">
                    <v-btn text color="accent" v-on="on">
                        Add Node
                        <v-icon small color="accent">mdi-plus</v-icon>
                    </v-btn>
                </template>
                <v-list dense>
                    <v-list-item v-for="(item, index) in rawTypeOptions"
                                 :key="index"
                                 @click="createNewNode($event, item.Id)"
                                 dense
                    >
                        <v-list-item-title>
                            {{ item.Name }}
                        </v-list-item-title>
                    </v-list-item>

                </v-list>
            </v-menu>
            <v-divider vertical></v-divider>
        </v-toolbar-items>
    </v-toolbar>
</template>

<script lang="ts">

interface ResponseData {
    data : ProcessFlowNodeType[]
}

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import axios from 'axios'
import * as qs from 'query-string'
import { contactUsUrl, getAllProcessFlowNodeTypesAPIUrl, newProcessFlowNodeAPIUrl } from '../../../ts/url'
import { postFormUrlEncoded } from '../../../ts/http'

export default Vue.extend({
    data : () => ({
        rawTypeOptions: [] as ProcessFlowNodeType[]
    }),
    methods: {
        createNewNode(_ : MouseEvent, nodeTypeId : number) {
            // Create a new node of the given type.
            postFormUrlEncoded(newProcessFlowNodeAPIUrl, {
                nodeTypeId,
                flowId: VueSetup.store.getters.currentProcessFlowBasicData.Id,
                //@ts-ignore
                csrf: this.$root.csrf
            }).then((resp) => {
                // TODO: Make this more efficient and just do a local adjustment of the data?
                //       That'd require some more syncing stuff...which is fancier.
                // Force a refresh of the data for the currently selected process flow.
                //@ts-ignore
                VueSetup.store.dispatch('refreshCurrentProcessFlowFullData', this.$root.csrf)
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please reload the page and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
    },
    mounted() {
        // Send out an AJAX request to get all available node types so that we don't have to
        // manually keep this in sync with the available types on the server.
        let passedData = Object()
        //@ts-ignore
        passedData['csrf'] = this.$root.csrf

        axios.get(getAllProcessFlowNodeTypesAPIUrl+ '?' + qs.stringify(passedData)).then((resp : ResponseData) => {
            this.rawTypeOptions = resp.data
        }).catch((err) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong, please reload the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }
})

</script>

<style scoped>

.v-menu__content {
    border-radius: 0px !important;
}

</style>
