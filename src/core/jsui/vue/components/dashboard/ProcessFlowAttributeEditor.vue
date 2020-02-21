<template>
    <v-navigation-drawer absolute right :style="clipStyle" ref="attrNavDrawer" :value="showHide" :width="400">
        <v-tabs v-model="tab"
                grow>
            <v-tab>Node</v-tab>
            <v-tab>Link</v-tab>
        </v-tabs>
        <section v-if="enabled" class="ma-1" style="max-height: calc(100% - 48px);">
            <v-tabs-items v-model="tab">
                <v-tab-item>
                    <process-flow-node-attribute-editor></process-flow-node-attribute-editor>
                </v-tab-item>

                <v-tab-item>
                    <process-flow-risk-control-linkage></process-flow-risk-control-linkage>
                </v-tab-item>
            </v-tabs-items>
        </section>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import ProcessFlowNodeAttributeEditor from './ProcessFlowNodeAttributeEditor.vue'
import ProcessFlowRiskControlLinkage from './ProcessFlowRiskControlLinkage.vue'

export default Vue.extend({
    props: {
        customClipHeight : Number,
        showHide : Boolean
    },
    data : () => ({
        tab : 0
    }),
    components: {
        ProcessFlowNodeAttributeEditor,
        ProcessFlowRiskControlLinkage
    },
    computed: {
        clipStyle() : any {
            return {
                "height":  "100vh !important",
                "max-height": "calc(100% - " + this.customClipHeight.toString()  + "px) !important",
                "top" : this.customClipHeight.toString() + "px"
            }
        },
        enabled() : boolean {
            return VueSetup.store.getters.isNodeSelected
        },
    },
})

</script>
