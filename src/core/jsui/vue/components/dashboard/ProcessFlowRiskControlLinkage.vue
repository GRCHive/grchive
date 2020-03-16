<template>
    <div>
        <v-dialog persistent max-width="40%" v-model="showHideExisting">
            <v-card>
                <v-card-title>
                    Link Controls
                </v-card-title>
                <v-divider></v-divider>

                <control-table-with-controls
                    v-if="showHideExisting"
                    class="ma-4"
                    v-model="selectedControls"
                    :exclude="controlsForRisk(existingControlRiskId)"
                    disable-new
                    disable-delete
                    enable-select
                >
                </control-table-with-controls>

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
                        :disabled="selectedControls.length == 0"
                    >
                        Link
                    </v-btn>

                </v-card-actions>
            </v-card>
        </v-dialog>

        <template v-for="risk in currentFlowRisks">
            <v-list-item
                :key="`risk-${risk.Id}`"
                @click="clickRiskBar(risk.Id)"
            >
                <v-list-item-action>
                    <v-btn icon @mousedown.stop @click.stop="showExistingControlDialog(risk.Id)">
                        <v-icon>mdi-plus</v-icon>
                    </v-btn>
                </v-list-item-action>

                <v-list-item-content>
                    <v-list-item-title>{{ risk.Identifier }}</v-list-item-title>
                    <v-list-item-subtitle>{{ risk.Name }}</v-list-item-subtitle>
                </v-list-item-content>

                <v-spacer></v-spacer>

                <v-list-item-action>
                    <v-btn icon>
                        <v-icon>{{ riskExpansion[risk.Id] ? "mdi-chevron-up" : "mdi-chevron-down" }}</v-icon>
                    </v-btn>
                </v-list-item-action>
            </v-list-item>

            <control-table
                :resources="controlsForRisk(risk.Id)"
                :key="`controls-${risk.Id}`"
                mini
                use-crud-delete
                @delete="unlink(risk.Id, arguments[0].Id)"
                v-if="riskExpansion[risk.Id]"
            >
            </control-table>
        </template>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import VueSetup from '../../../ts/vueSetup'
import ControlTable from '../../generic/ControlTable.vue'
import ControlTableWithControls from '../../generic/resources/ControlTableWithControls.vue'
import { addExistingControls, deleteControls } from '../../../ts/api/apiControls'
import { contactUsUrl } from '../../../ts/url'

@Component({
    components: {
        ControlTable,
        ControlTableWithControls
    }
})
export default class ProcessFlowRiskControlLinkage extends Vue {
    existingControlRiskId: number = -1
    showHideExisting : boolean = false

    selectedControls : ProcessFlowControl[] = []
    riskExpansion : Record<number, boolean> = Object()

    clickRiskBar(riskId : number) {
        Vue.set(this.riskExpansion, riskId, !this.riskExpansion[riskId])
    }

    get currentFlowRisks() : ProcessFlowRisk[] {
        return VueSetup.store.getters.linkedRiskList
    }

    @Watch('currentFlowRisks')
    resetRiskExpansion() {
        this.riskExpansion = this.currentFlowRisks.map(() => false)
    }

    get controlsForRisk() : (a : number) => ProcessFlowControl[] {
        return (riskId : number) : ProcessFlowControl[] => {
            return VueSetup.store.getters.controlsForRisk(riskId)
        }
    }

    showExistingControlDialog(riskId: number) {
        this.existingControlRiskId = riskId
        this.showHideExisting = true
        this.selectedControls = []
    }

    cancelLink() {
        this.existingControlRiskId = -1
        this.showHideExisting = false
        this.selectedControls = []
    }

    saveLink() {
        addExistingControls({
            nodeId: -1,
            riskId: this.existingControlRiskId,
            controlIds: this.selectedControls.map((ele : ProcessFlowControl) => ele.Id)
        }).then(() => {
            this.selectedControls.forEach((ele : ProcessFlowControl) => {
                VueSetup.store.commit('addControlToRisk', {
                    controlId: ele.Id,
                    riskId: this.existingControlRiskId,
                })
            })
            this.existingControlRiskId = -1
            this.showHideExisting = false
            this.selectedControls = []
        }).catch((err : any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong, please reload the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    unlink(riskId: number, controlId : number) {
        deleteControls({
            nodeId: -1,
            riskIds: [riskId],
            controlIds: [controlId],
            global: false
        }).then(() => {
            VueSetup.store.dispatch('deleteBatchControls', {
                nodeId: -1,
                riskIds: [riskId],
                controlIds: [controlId],
                global: false
            })
        }).catch((err : any) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong, please reload the page and try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);

        })
    }
}

</script>
