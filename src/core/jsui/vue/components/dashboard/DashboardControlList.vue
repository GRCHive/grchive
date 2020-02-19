<template>
    <section class="ma-4">
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
            <v-list-item-action>
                <v-dialog v-model="showHideCreateNewControl" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            New
                        </v-btn>
                    </template>
                    <create-new-control-form
                        :node-id="-1"
                        :risk-id="-1"
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
            :resources="allControls"
            :search="filterText"
            use-crud-delete
            confirm-delete
            @delete="deleteControl"
        >
        </control-table>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import { getAllControls, TAllControlInput, TAllControlOutput} from '../../../ts/api/apiControls'
import { deleteControls, TDeleteControlInput, TDeleteControlOutput} from '../../../ts/api/apiControls'
import { contactUsUrl } from '../../../ts/url'
import CreateNewControlForm from './CreateNewControlForm.vue'
import { PageParamsStore } from '../../../ts/pageParams'
import ControlTable from '../../generic/ControlTable.vue'
import AdvancedControlFilters from '../../generic/filters/AdvancedControlFilters.vue'
import { ControlFilterData, NullControlFilterData } from '../../../ts/controls'

export default Vue.extend({
    data : () => ({
        allControls: [] as ProcessFlowControl[],
        filterText : "",
        showHideCreateNewControl : false,
        controlFilter: NullControlFilterData,
    }),
    components: {
        CreateNewControlForm,
        ControlTable,
        AdvancedControlFilters
    },
    methods: {
        refreshControls() {
            getAllControls(<TAllControlInput>{
                orgName: PageParamsStore.state.organization!.OktaGroupName,
                filter: this.controlFilter,
            }).then((resp : TAllControlOutput) => {
                this.allControls = resp.data
            }).catch((err : any) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        deleteControl(control : ProcessFlowControl, global : boolean) {
            deleteControls(<TDeleteControlInput>{
                nodeId: -1,
                riskIds: [-1],
                controlIds: [control.Id],
                global: global
            }).then(() => {
                this.allControls.splice(
                    this.allControls.findIndex((ele : ProcessFlowControl) =>
                        ele.Id == control.Id),
                    1)
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        saveNewControl(control : ProcessFlowControl) {
            this.allControls.unshift(control)
            this.showHideCreateNewControl = false
        },
        cancelNewControl() {
            this.showHideCreateNewControl = false
        },

    },
    watch: {
        controlFilter: {
            deep: true,

            handler() {
                this.refreshControls()
            }
        }
    },
    mounted() {
        this.refreshControls()
    }
})
</script>

<style scoped>

.headerItem {
    max-height: 30px !important;
    min-height: 30px !important;
}

</style>
