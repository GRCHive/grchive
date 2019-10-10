<template>
    <section id="flowRenderer" class="ma-0" :style="contentStyle">
        <section class="max-height" v-if="hasProcessFlowToRender">
            <h1>Has Process Flow</h1>
        </section>

        <section class="max-height" v-else>
            <v-row class="max-height" align="center" width="100%">
                <v-col>
                    <v-row justify="center">
                        <p class="display-1">This process flow is empty!</p>
                    </v-row>
                    <v-row justify="center">
                        <p class="body-1">Get started by clicking the "Add Node" button.</p>
                    </v-row>
                    <v-row justify="center">
                        <v-btn icon @click="refreshProcessFlow"
                                    :disabled="processFlowLoading"
                                    :loading="processFlowLoading"
                        >
                            <v-icon x-large>mdi-refresh</v-icon>
                        </v-btn>
                    </v-row>
                </v-col>
            </v-row>
        </section>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup'
import { isObjectEmpty } from '../../../ts/object'

export default Vue.extend({
    props: {
        contentMaxHeightClip: Number,
        contentMaxWidthClip: Number
    },
    computed: {
        hasProcessFlowToRender() : boolean {
            return isObjectEmpty(VueSetup.store.state.currentProcessFlowFullData)
        },
        contentStyle() {
            return {
                "height": "100%",
                "max-height": `calc(100% - ${this.contentMaxHeightClip.toString()}px)`,
                "width": "100%",
                "max-width": `calc(100% - ${this.contentMaxWidthClip.toString()}px)`
            }
        },
        processFlowLoading() : boolean {
            return VueSetup.store.getters.isFullRequestInProgress
        }
    },
    methods: {
        refreshProcessFlow() {
            //@ts-ignore
            VueSetup.store.dispatch('refreshCurrentProcessFlowFullData', this.$root.csrf)
        }
    }
})
</script>
