<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Process Flows
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
                <v-dialog v-model="showHideCreateNewFlow" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            New
                        </v-btn>
                    </template>
                    <create-new-process-flow-form
                        @do-save="onCreateNewFlow"
                        @do-cancel="showHideCreateNewFlow = false"
                    >
                    </create-new-process-flow-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <process-flow-table
            :resources="allFlows"
            :search="filterText"
            use-crud-delete
            confirm-delete
            @delete="deleteFlow"
        >
        </process-flow-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import CreateNewProcessFlowForm from './CreateNewProcessFlowForm.vue'
import { getAllProcessFlow, TGetAllProcessFlowInput, TGetAllProcessFlowOutput } from '../../../ts/api/apiProcessFlow'
import { deleteProcessFlow, TDeleteProcessFlowInput, TDeleteProcessFlowOutput } from '../../../ts/api/apiProcessFlow'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import ProcessFlowTable from '../../generic/ProcessFlowTable.vue'

export default Vue.extend({
    data : () => ({
        showHideCreateNewFlow: false,
        allFlows: [] as ProcessFlowBasicData[],
        filterText: "",
    }),
    components: {
        CreateNewProcessFlowForm,
        ProcessFlowTable
    },
    methods: {
        onCreateNewFlow(data : ProcessFlowBasicData) {
            this.showHideCreateNewFlow = false
            this.allFlows.unshift(data)
        },
        refreshFlows() {
            getAllProcessFlow(<TGetAllProcessFlowInput>{
                requested: -1,
                organization: PageParamsStore.state.organization!.OktaGroupName
            }).then((resp : TGetAllProcessFlowOutput) => {
                this.allFlows = resp.data.Flows
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        deleteFlow(flow : ProcessFlowBasicData) {
            deleteProcessFlow(<TDeleteProcessFlowInput>{
                flowId: flow.Id,
            }).then((resp : TDeleteProcessFlowOutput) => {
                this.allFlows.splice(
                    this.allFlows.findIndex((ele : ProcessFlowBasicData) =>
                        ele.Id == flow.Id),
                    1)
            }).catch((err) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong, please try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
    },
    mounted() {
        this.refreshFlows()
    }
})

</script>

<style scoped>

.headerItem {
    max-height: 30px !important;
    min-height: 30px !important;
}

</style>
