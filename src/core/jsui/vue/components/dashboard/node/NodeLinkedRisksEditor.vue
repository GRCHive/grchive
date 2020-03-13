<template>
    <div>
        <v-dialog persistent max-width="40%" v-model="showHideNew">
            <create-new-risk-form
                :node-id="currentNode.Id"
                @do-cancel="showHideNew = false"
                @do-save="saveRisk"
            >
            </create-new-risk-form>
        </v-dialog>

        <v-dialog persistent max-width="40%" v-model="showHideExisting">
            <v-card>
                <v-card-title>
                    Link Risks
                </v-card-title>
                <v-divider></v-divider>

                <risk-table-with-controls
                    class="ma-4"
                    v-model="selectedRisks"
                    :exclude="nodeRisks"
                    disable-new
                    disable-delete
                    enable-select
                >
                </risk-table-with-controls>

                <v-card-actions>
                    <v-btn
                        color="error"
                        @click="cancelLink"
                    >
                        Cancel
                    </v-btn>
                    <v-spacer></v-spacer>
                    <v-btn
                        color="success"
                        @click="saveLink"
                        :disabled="selectedRisks.length == 0"
                    >
                        Link
                    </v-btn>

                </v-card-actions>

            </v-card>
        </v-dialog>

        <v-list dense class="pa-0">
            <v-list-item class="pa-0">
                <v-list-item-action class="ma-0">
                    <v-btn icon @click="minimize = !minimize">
                        <v-icon small>
                            {{ !minimize ? "mdi-window-minimize" : "mdi-arrow-expand-all" }}
                        </v-icon>
                    </v-btn>
                </v-list-item-action>

                <v-subheader class="flex-grow-1 pr-0">
                    LINKED RISKS
                </v-subheader>

                <v-list-item-action class="ma-0">
                    <v-menu bottom left offset-y>
                        <template v-slot:activator="{ on }">
                            <v-btn
                                icon
                                v-on="on"
                            >
                                <v-icon small>
                                    mdi-plus
                                </v-icon>
                            </v-btn>
                        </template>

                        <v-list dense>
                            <v-list-item @click="showHideNew = true">
                                <v-list-item-title>
                                    New
                                </v-list-item-title>
                            </v-list-item>

                            <v-list-item @click="showHideExisting = true">
                                <v-list-item-title>
                                    Add Existing
                                </v-list-item-title>
                            </v-list-item>
                        </v-list>

                    </v-menu>
                </v-list-item-action>
            </v-list-item>
        </v-list>

        <risk-table
            :resources="nodeRisks"
            use-crud-delete
            @delete="deleteLinkedRisk"
            v-if="!minimize"
        >
        </risk-table>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import VueSetup from '../../../../ts/vueSetup'
import RiskTable from '../../../generic/RiskTable.vue'
import RiskTableWithControls from '../../../generic/resources/RiskTableWithControls.vue'
import CreateNewRiskForm from '../CreateNewRiskForm.vue'
import {
    addExistingRisk,
    deleteRisk
} from '../../../../ts/api/apiRisks'
import { contactUsUrl } from '../../../../ts/url'

@Component({
    components: {
        RiskTable,
        RiskTableWithControls,
        CreateNewRiskForm
    }
})
export default class NodeLinkedRisksEditor extends Vue {
    showHideNew : boolean = false
    showHideExisting : boolean = false
    minimize : boolean = false

    selectedRisks : ProcessFlowRisk[] = []

    get currentNode() : ProcessFlowNode {
        return VueSetup.store.getters.currentNodeInfo
    }

    get nodeRisks() : ProcessFlowRisk[] {
        return VueSetup.store.getters.risksForNode(this.currentNode.Id)
    }

    deleteLinkedRisk(risk : ProcessFlowRisk, global : boolean) {
        deleteRisk({
            nodeId: this.currentNode.Id,
            riskIds: [risk.Id],
            global: false,
        }).then(() => {
            VueSetup.store.dispatch('deleteBatchRisks', {
                nodeId: this.currentNode.Id,
                riskIds: [risk.Id],
                global: global,
            })
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    saveRisk(risk : ProcessFlowRisk) {
        VueSetup.store.commit('setRisk', risk)
        VueSetup.store.commit('addRisksToNode', {
            nodeId: this.currentNode.Id,
            riskIds: [risk.Id],
        })
        this.showHideNew = false
    }

    cancelLink() {
        this.showHideExisting = false
        this.selectedRisks = []
    }

    saveLink() {
        addExistingRisk({
            nodeId: this.currentNode.Id,
            riskIds: this.selectedRisks.map((ele : ProcessFlowRisk) => ele.Id),
        }).then(() => {
            this.selectedRisks.forEach(this.saveRisk)
            this.showHideExisting = false
            this.selectedRisks = []
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }
}
</script>
