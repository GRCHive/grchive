<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Control: {{ fullControlData.Control.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ fullControlData.Control.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="8">
                        <create-new-control-form ref="editControl"
                                                 :node-id="-1"
                                                 :risk-id="-1"
                                                 :edit-mode="true"
                                                 :control="fullControlData.Control"
                                                 :staged-edits="true"
                                                 @do-save="onEditControl">
                        </create-new-control-form>
                    </v-col>

                    <v-col cols="4">
                        <v-card class="mb-4">
                            <v-card-title>
                                <span class="mr-2">
                                    Documentation
                                </span>
                                <v-spacer></v-spacer>
                            </v-card-title>
                            <v-divider></v-divider>
                        </v-card>

                        <v-card class="mb-4">
                            <v-card-title>
                                Related Nodes
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-list two-line>
                                <v-list-item v-for="(item, index) in fullControlData.Nodes"
                                             :key="index"
                                >
                                    <v-list-item-content>
                                        <v-list-item-title>
                                            {{ item.Name }}
                                        </v-list-item-title>

                                        <v-list-item-subtitle>
                                            {{ item.Description }}
                                        </v-list-item-subtitle>
                                    </v-list-item-content>
                                </v-list-item>
                            </v-list>
                        </v-card>

                        <v-card>
                            <v-card-title>
                                Related Risks
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-list two-line>
                                <v-list-item v-for="(item, index) in fullControlData.Risks"
                                             :key="index"
                                             :href="generateRiskUrl(item.Id)"
                                >
                                    <v-list-item-content>
                                        <v-list-item-title>
                                            {{ item.Name }}
                                        </v-list-item-title>

                                        <v-list-item-subtitle>
                                            {{ item.Description }}
                                        </v-list-item-subtitle>
                                    </v-list-item-content>
                                </v-list-item>
                            </v-list>
                        </v-card>

                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import CreateNewControlForm from './CreateNewControlForm.vue'
import { FullControlData } from '../../../ts/controls'
import { getSingleControl, TSingleControlInput, TSingleControlOutput } from '../../../ts/api/apiControls'
import { createRiskUrl, contactUsUrl } from '../../../ts/url'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import { deleteControlDocCat, TDeleteControlDocCatInput, TDeleteControlDocCatOutput } from '../../../ts/api/apiControlDocumentation'
import { PageParamsStore } from '../../../ts/pageParams'

export default Vue.extend({
    data: () => ({
        expandDescription: false,
        ready: false,
        fullControlData: Object() as FullControlData,
        showHideNewCat: false,
        showHideEditCat : false,
        showHideDeleteCat : false,
    }),
    computed: {
    },
    methods: {
        refreshData() {
            let data = window.location.pathname.split('/')
            let controlId = Number(data[data.length - 1])

            getSingleControl(<TSingleControlInput>{
                controlId: controlId
            }).then((resp : TSingleControlOutput) => {
                this.fullControlData = resp.data
                this.ready = true

                Vue.nextTick(() => {
                    //@ts-ignore
                    this.$refs.editControl.clearForm()
                })
            }).catch((err : any) => {
                window.location.replace('/404')
            })
        },
        onEditControl(control : ProcessFlowControl) {
            this.fullControlData.Control.Name = control.Name
            this.fullControlData.Control.Description = control.Description
            this.fullControlData.Control.ControlTypeId = control.ControlTypeId
            this.fullControlData.Control.FrequencyType = control.FrequencyType
            this.fullControlData.Control.FrequencyInterval = control.FrequencyInterval
            this.fullControlData.Control.OwnerId = control.OwnerId

            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editControl.clearForm()
            })
        },
        generateRiskUrl(riskId : number) : string {
            return createRiskUrl(
                PageParamsStore.state.organization!.OktaGroupName,
                riskId)
        },
    },
    components: {
        CreateNewControlForm,
        GenericDeleteConfirmationForm,
    },
    mounted() {
        this.refreshData()
    }
})

</script>

<style scoped>

>>>.v-select__selections input {
    display: none !important;
}

>>>.v-tabs .v-slide-group {
    display: none !important;
}

</style>
