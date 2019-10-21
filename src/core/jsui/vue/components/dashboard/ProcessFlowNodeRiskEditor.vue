<template>
    <section>
        <v-list-item class="pa-1">
            <v-list-item-action class="ma-1">
                <v-btn color="primary">
                    Add Existing
                </v-btn>
            </v-list-item-action>
            <v-list-item-action class="ma-1">
                <v-dialog v-model="showHideCreateNewRisk" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="success" v-on="on">
                            Create New
                        </v-btn>
                    </template>
                    <create-new-risk-form
                        :node-id="currentNode.Id"
                        @do-save="saveNewRisk"
                        @do-cancel="cancelNewRisk">
                    </create-new-risk-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from '../../../ts/vueSetup' 
import CreateNewRiskForm from './CreateNewRiskForm.vue'

export default Vue.extend({
    data : () => ({
        showHideCreateNewRisk : false
    }),
    components : {
        CreateNewRiskForm
    },
    computed : {
        currentNode() : ProcessFlowNode {
            return VueSetup.store.getters.currentNodeInfo
        },
    },
    methods : {
        saveNewRisk() {
            this.showHideCreateNewRisk = false
        },
        cancelNewRisk() {
            this.showHideCreateNewRisk = false
        }
    }
})
</script>
