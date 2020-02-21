<template>
    <div>
        <v-dialog persistent max-width="40%" v-model="showHideExisting">
            <v-card>
                <v-card-title>
                    Link Controls
                </v-card-title>
                <v-divider></v-divider>

                <control-table
                    v-model="selectedControls"
                    :resources="unlinkedControls"
                    selectable
                    multi
                >
                </control-table>

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

        <v-list two-line>
            <v-list-group
                v-for="risk in currentFlowRisks"
                :key="`risk-${risk.Id}`"
            >
                <template v-slot:activator>
                    <v-list-item-action>
                        <v-btn icon @mousedown.stop @click.stop="showExistingControlDialog(risk.Id)">
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </v-list-item-action>

                    <v-list-item-content>
                        <v-list-item-title>{{ risk.Name }}</v-list-item-title>
                        <v-list-item-subtitle>{{ risk.Description }}</v-list-item-subtitle>
                    </v-list-item-content>
                </template>
            </v-list-group>
        </v-list>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import VueSetup from '../../../ts/vueSetup'
import ControlTable from '../../generic/ControlTable.vue'
import { addExistingControls } from '../../../ts/api/apiControls'
import { contactUsUrl } from '../../../ts/url'

@Component({
    components: {
        ControlTable
    }
})
export default class ProcessFlowRiskControlLinkage extends Vue {
    existingControlRiskId: number = -1
    showHideExisting : boolean = false

    selectedControls : ProcessFlowControl[] = []

    get unlinkedControls() : ProcessFlowControl[] {
        if (this.existingControlRiskId == -1) {
            return []
        }

        let controlsForRisk : ProcessFlowControl[] = VueSetup.store.getters.controlsForRisk(this.existingControlRiskId)
        let usedSet : Set<number> = new Set(controlsForRisk.map((ele : ProcessFlowControl) => ele.Id))
        return VueSetup.store.getters.controlList.filter((ele : ProcessFlowControl) => !usedSet.has(ele.Id))
    }

    get currentFlowRisks() : ProcessFlowRisk[] {
        return VueSetup.store.getters.riskList
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
}

</script>
