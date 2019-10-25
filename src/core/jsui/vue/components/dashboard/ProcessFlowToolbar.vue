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
                    <v-list-item dense>
                        <v-list-item-action class="ma-0">
                            <v-btn icon @click.stop="decreaseZoom">
                                <v-icon small>mdi-minus</v-icon>
                            </v-btn>
                        </v-list-item-action>
                        <v-list-item-title class="text-center">
                            Zoom: {{ zoomPercentage }}%
                        </v-list-item-title>
                        <v-list-item-action class="ma-0">
                            <v-btn icon @click.stop="increaseZoom">
                                <v-icon small>mdi-plus</v-icon>
                            </v-btn>
                        </v-list-item-action>
                    </v-list-item>
                    <v-divider></v-divider>
                    <v-list-item dense @click="resetZoom">
                        <v-list-item-title>
                            Reset Zoom
                        </v-list-item-title>
                    </v-list-item>
                    <v-divider></v-divider>
                    <v-list-item dense @click="resetView">
                        <v-list-item-title>
                            Reset View
                        </v-list-item-title>
                    </v-list-item>
                    <v-divider></v-divider>
                    <v-list-item dense @click="fitToGraph">
                        <v-list-item-title>
                            Fit to Graph
                        </v-list-item-title>
                    </v-list-item>
                    <v-divider></v-divider>
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
import RenderLayout from '../../../ts/render/renderLayout'
import LocalSettings from '../../../ts/localSettings'
import { contactUsUrl, newProcessFlowNodeAPIUrl } from '../../../ts/url'
import { postFormUrlEncoded } from '../../../ts/http'

export default Vue.extend({
    computed: {
        rawTypeOptions() : ProcessFlowNodeType[] {
            return MetadataStore.state.nodeTypes
        },
        zoomPercentage() : number {
            let zoom : number = LocalSettings.state.viewBoxZoom
            return Math.round(zoom * 100.0)
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
        handleScroll(e : WheelEvent) {
            if (e.deltaY < 0.0) {
                this.increaseZoom()
            } else {
                this.decreaseZoom()
            }
        },
        resetZoom() {
            LocalSettings.commit('setViewBoxZoom', 1.0)
        },
        resetView() {
            LocalSettings.commit('setViewBoxTransform', { tx: 0, ty : 0 })
            LocalSettings.commit('setViewBoxZoom', 1.0)
        },
        increaseZoom() {
            if (!RenderLayout.store.state.ready) {
                return
            }

            LocalSettings.commit('setViewBoxZoom', LocalSettings.state.viewBoxZoom + 0.05)
        },
        decreaseZoom() {
            if (!RenderLayout.store.state.ready) {
                return
            }

            LocalSettings.commit('setViewBoxZoom', LocalSettings.state.viewBoxZoom - 0.05)
        },
        fitToGraph() {
            let graphBbox = (<SVGSVGElement>RenderLayout.store.state.fullGraph).getBBox()
            LocalSettings.commit('setViewBoxTransform', { tx: graphBbox.x, ty : graphBbox.y })

            let zoomX = RenderLayout.store.state.rendererRect.width / graphBbox.width
            let zoomY = RenderLayout.store.state.rendererRect.height / graphBbox.height
            LocalSettings.commit('setViewBoxZoom', Math.min(zoomX, zoomY))
        }
    },
    mounted() {
        document.addEventListener('keydown', this.handleHotkeys)
        document.addEventListener('wheel', this.handleScroll)
    }
})

</script>

<style scoped>

.v-menu__content {
    border-radius: 0px !important;
}

</style>
