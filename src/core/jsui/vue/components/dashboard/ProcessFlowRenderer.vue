<template>
    <section id="flowRenderer" class="ma-0" :style="contentStyle">
        <section class="max-height" v-if="hasProcessFlowToRender">
            <h1>Has Process Flow</h1>
        </section>

        <section class="max-height" v-else>
            <v-row class="max-height" align="center" justify="center">
                <h1>No Content Here</h1>
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
        }
    }
})
</script>
