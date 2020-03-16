<template>
    <v-navigation-drawer
        absolute
        right
        :style="clipStyle"
        ref="attrNavDrawer"
        :value="showHide"
        :width="400"
        disable-resize-watcher
        mobile-break-point="-1"
    >
        <v-tabs v-model="tab"
                grow>
            <v-tab v-if="nodeSelected">Node</v-tab>
            <v-tab>Process Flow</v-tab>
        </v-tabs>
        <div class="ma-1" style="max-height: calc(100% - 48px);">
            <v-tabs-items v-model="tab">
                <v-tab-item v-if="nodeSelected">
                    <process-flow-node-attribute-editor></process-flow-node-attribute-editor>
                </v-tab-item>

                <v-tab-item>
                    <process-flow-flow-attribute-editor></process-flow-flow-attribute-editor>
                </v-tab-item>
            </v-tabs-items>
        </div>
    </v-navigation-drawer>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import VueSetup from '../../../ts/vueSetup' 
import ProcessFlowNodeAttributeEditor from './ProcessFlowNodeAttributeEditor.vue'
import ProcessFlowFlowAttributeEditor from './ProcessFlowFlowAttributeEditor.vue'

const Props = Vue.extend({
    props: {
        customClipHeight : Number,
        showHide : Boolean
    },
})

@Component({
    components: {
        ProcessFlowNodeAttributeEditor,
        ProcessFlowFlowAttributeEditor
    },
})
export default class ProcessFlowAttributeEditor extends Props {
    tab : number =  0
    get clipStyle() : any {
        return {
            "height":  "100vh !important",
            "max-height": "calc(100% - " + this.customClipHeight.toString()  + "px) !important",
            "top" : this.customClipHeight.toString() + "px"
        }
    }

    get nodeSelected() : boolean {
        return VueSetup.store.getters.isNodeSelected
    } 

    @Watch('nodeSelected')
    resetTabs() {
        Vue.nextTick(() => {
            this.tab = 1
        })
    }
}

</script>
