<template>
    <section>
        <v-dialog v-model="showHideNewControl" persistent max-width="40%">
            <create-new-control-form
                :node-id="currentNode.Id"
                :risk-id="currentRelevantRiskId"
                @do-cancel="onCancelNewControl"
                @do-save="onSaveNewControl">
            </create-new-control-form>
        </v-dialog>

        <v-dialog v-model="showHideExistingControl" persistent max-width="40%">
        </v-dialog>

        <v-list two-line>
            <v-list-group v-for="(item, index) in risksForNode"
                          :key="index"
                          class="pa-1"
                          :value="true">
                <template v-slot:activator>
                    <v-list-item-action>
                        <v-menu offset-y>
                            <template v-slot:activator="{ on }">
                                <v-btn icon v-on="on" @mousedown.stop @click.stop>
                                    <v-icon>mdi-plus</v-icon>
                                </v-btn>
                            </template>
                            <v-list dense>
                                <v-list-item @click="showNewControlDialog(item.Id)">
                                    New Control
                                </v-list-item>
                                <v-list-item @click="showExistingControlDialog(item.Id)">
                                    Existing Control
                                </v-list-item>
                            </v-list>
                        </v-menu>
                    </v-list-item-action>
                    
                    <v-list-item-content>
                        <v-list-item-title>
                            Risk:&nbsp;{{ item.Name }}
                        </v-list-item-title>

                        <v-list-item-subtitle>
                            {{ item.Description }}
                        </v-list-item-subtitle>
                    </v-list-item-content>
                </template>
            </v-list-group>
        </v-list>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from "../../../ts/vueSetup"
import CreateNewControlForm from './CreateNewControlForm.vue'

export default Vue.extend({
    data : () => ({
        showHideNewControl : false,
        showHideExistingControl : false,
        currentRelevantRiskId : -1
    }),
    components : {
        CreateNewControlForm
    },
    computed: {
        currentNode() : ProcessFlowNode {
            return VueSetup.store.getters.currentNodeInfo
        },
        risksForNode() : ProcessFlowRisk[] {
            return VueSetup.store.getters.risksForNode(this.currentNode.Id)
        },
    },
    methods : {
        showNewControlDialog(riskId : number) {
            this.currentRelevantRiskId = riskId
            this.showHideNewControl = true
        },
        showExistingControlDialog(riskId : number) {
            this.currentRelevantRiskId = riskId
            this.showHideExistingControl = true
        },
        onCancelNewControl() {
            this.showHideNewControl = false
        },
        onSaveNewControl() {
            this.showHideNewControl = false
        }
    }
})

</script>
