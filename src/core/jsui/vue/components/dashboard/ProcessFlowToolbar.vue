<template>
    <v-toolbar flat height="30px">
        <v-toolbar-items>
            <v-divider vertical></v-divider>
            <v-menu offset-y>
                <template v-slot:activator="{ on }">
                    <v-btn text color="accent" v-on="on">
                        View
                        <v-icon small color="accent">mdi-chevron-down</v-icon>
                    </v-btn>
                </template>

                <v-list dense>
                    <v-list-item dense @click="resetView">
                        <v-list-item-title>
                            Reset View
                        </v-list-item-title>
                    </v-list-item>
                </v-list>
            </v-menu>
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
import MetadataStore from '../../../ts/metadata'
import LocalSettings from '../../../ts/localSettings'
import { contactUsUrl, newProcessFlowNodeAPIUrl } from '../../../ts/url'
import { postFormUrlEncoded } from '../../../ts/http'

export default Vue.extend({
    computed: {
        rawTypeOptions() : ProcessFlowNodeType[] {
            return MetadataStore.state.nodeTypes
        }
    },
    methods: {
        createNewNode(_ : MouseEvent, nodeTypeId : number) {
            // Create a new node of the given type.
            postFormUrlEncoded(newProcessFlowNodeAPIUrl, {
                typeId: nodeTypeId,
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
        handleHotkeys(e : KeyboardEvent) {
            if (e.code == "Delete") {
                VueSetup.store.dispatch('requestDeletionOfSelection', {
                    //@ts-ignore
                    csrf: this.$root.csrf
                })
                e.stopPropagation()
            }
        },
        resetView() {
            LocalSettings.commit('setViewBoxTransform', { tx: 0, ty : 0 })
        }
    },
    mounted() {
        document.addEventListener('keydown', this.handleHotkeys)
    }
})

</script>

<style scoped>

.v-menu__content {
    border-radius: 0px !important;
}

</style>
