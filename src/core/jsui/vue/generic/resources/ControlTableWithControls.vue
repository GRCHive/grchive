<template>
    <div>
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Controls
                </v-list-item-title>
            </v-list-item-content>
            <v-list-item-action>
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>
            <v-spacer></v-spacer>
            <v-list-item-action v-if="!disableNew">
                <v-dialog v-model="showHideCreateNewControl" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            New
                        </v-btn>
                    </template>
                    <create-new-control-form
                        :node-id="-1"
                        @do-save="saveNewControl"
                        @do-cancel="cancelNewControl">
                    </create-new-control-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>

        <advanced-control-filters
            v-model="controlFilter"
        >
        </advanced-control-filters>

        <v-divider></v-divider>

        <control-table
            :value="value"
            :resources="filteredControls"
            :search="filterText"
            :use-crud-delete="!disableDelete"
            :confirm-delete="!disableDelete"
            @delete="deleteSelectedControl"
            @input="modifySelected"
            :selectable="enableSelect"
            :multi="enableSelect"
        >
        </control-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { getAllControls, TAllControlOutput } from '../../../ts/api/apiControls'
import { deleteControls, TDeleteControlOutput } from '../../../ts/api/apiControls'
import { contactUsUrl } from '../../../ts/url'
import CreateNewControlForm from '../../components/dashboard/CreateNewControlForm.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import ControlTable from '../ControlTable.vue'
import AdvancedControlFilters from '../filters/AdvancedControlFilters.vue'
import { ControlFilterData, NullControlFilterData } from '../../../ts/controls'

const Props = Vue.extend({
    props: {
        value: {
            type: Array,
            default: () => [],
        },
        exclude: {
            type: Array,
            default: () => [],
        },
        disableNew: {
            type: Boolean,
            default: false,
        },
        disableDelete: {
            type: Boolean,
            default: false,
        },
        enableSelect: {
            type: Boolean,
            default: false,
        }
    }
})

@Component({
    components: {
        CreateNewControlForm,
        ControlTable,
        AdvancedControlFilters
    }
})
export default class ControlTableWithControls extends Props {
    allControls: ProcessFlowControl[] = []
    filterText : string = ""
    showHideCreateNewControl: boolean = false
    controlFilter: ControlFilterData = JSON.parse(JSON.stringify(NullControlFilterData))

    get excludeSet() : Set<number> {
        return new Set<number>(this.exclude.map((ele : any) => ele.Id))
    }

    get filteredControls() : ProcessFlowControl[] {
        return this.allControls.filter((ele : ProcessFlowControl) => !this.excludeSet.has(ele.Id))
    }

    @Watch('controlFilter', {deep:true})
    refreshControls() {
        getAllControls({
            orgName: PageParamsStore.state.organization!.OktaGroupName,
            filter: this.controlFilter,
        }).then((resp : TAllControlOutput) => {
            this.allControls = resp.data
        }).catch((err) => {
            //@ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    saveNewControl(control : ProcessFlowControl) {
        this.allControls.unshift(control)
        this.showHideCreateNewControl = false
    }

    cancelNewControl() {
        this.showHideCreateNewControl = false
    }

    deleteSelectedControl(control : ProcessFlowControl, global : boolean) {
        deleteControls({
            nodeId: -1,
            riskIds: [-1],
            controlIds: [control.Id],
            global: true,
        }).then((resp : TDeleteControlOutput) => {
            this.allControls.splice(
                this.allControls.findIndex((ele : ProcessFlowControl) =>
                    ele.Id == control.Id),
                1)
        }).catch((err) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    modifySelected(vals : ProcessFlowControl[]) {
        this.$emit('input', vals)
    }

    mounted() {
        this.refreshControls()
    }
}

</script>
