<template>
    <div>
        <v-dialog persistent max-width="40%" v-model="showHideNew">
            <create-new-control-form
                :node-id="currentNode.Id"
                :risk-id="-1"
                @do-cancel="showHideNew = false"
                @do-save="saveControl"
            >
            </create-new-control-form>
        </v-dialog>

        <v-dialog persistent max-width="40%" v-model="showHideExisting">
            <v-card>
                <v-card-title>
                    Link Controls
                </v-card-title>
                <v-divider></v-divider>

                <control-table-with-controls
                    class="ma-4"
                    v-model="selectedControls"
                    :exclude="nodeControls"
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
                    LINKED CONTROLS
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

        <control-table
            :resources="nodeControls"
            use-crud-delete
            mini
            @delete="deleteLinkedControl"
            v-if="!minimize"
        >
        </control-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import VueSetup from '../../../../ts/vueSetup'
import ControlTable from '../../../generic/ControlTable.vue'
import ControlTableWithControls from '../../../generic/resources/ControlTableWithControls.vue'
import CreateNewControlForm from '../CreateNewControlForm.vue'
import {
    addExistingControls,
    deleteControls
} from '../../../../ts/api/apiControls'
import { contactUsUrl } from '../../../../ts/url'

@Component({
    components: {
        ControlTable,
        CreateNewControlForm,
        ControlTableWithControls
    }
})
export default class NodeLinkedControlsEditor extends Vue {
    showHideNew : boolean = false
    showHideExisting : boolean = false
    minimize: boolean = false

    selectedControls : ProcessFlowControl[] = []

    get currentNode() : ProcessFlowNode {
        return VueSetup.store.getters.currentNodeInfo
    }

    get nodeControls() : ProcessFlowControl[] {
        return VueSetup.store.getters.controlsForNode(this.currentNode.Id)
    }

    get unlinkedControls() : ProcessFlowControl[] {
        let usedSet : Set<number> = new Set(this.nodeControls.map((ele : ProcessFlowControl) => ele.Id))
        return VueSetup.store.getters.controlList.filter((ele : ProcessFlowControl) => !usedSet.has(ele.Id))
    }

    deleteLinkedControl(control : ProcessFlowControl, global : boolean) {
        deleteControls({
            nodeId: this.currentNode.Id,
            riskIds: [-1],
            controlIds: [control.Id],
            global: false,
        }).then(() => {
            VueSetup.store.dispatch('deleteBatchControls', {
                nodeId: this.currentNode.Id,
                controlIds: [control.Id],
                riskIds: [],
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

    saveControl(control : ProcessFlowControl) {
        VueSetup.store.commit('setControl', {control})
        VueSetup.store.commit('addControlToNode', {
            controlId: control.Id,
            nodeId: this.currentNode.Id
        })
        this.showHideNew = false
    }

    cancelLink() {
        this.showHideExisting = false
        this.selectedControls = []
    }

    saveLink() {
        addExistingControls({
            nodeId: this.currentNode.Id,
            riskId: -1,
            controlIds: this.selectedControls.map((ele : ProcessFlowControl) => ele.Id)
        }).then(() => {
            this.selectedControls.forEach(this.saveControl)
            this.showHideExisting = false
            this.selectedControls = []
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
