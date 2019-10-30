<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Risk: {{ fullRiskData.Risk.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ fullRiskData.Risk.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>
            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="8">
                        <create-new-risk-form ref="editRisk"
                                              :node-id="-1"
                                              :edit-mode="true"
                                              :default-name="fullRiskData.Risk.Name"
                                              :default-description="fullRiskData.Risk.Description"
                                              :risk-id="fullRiskData.Risk.Id"
                                              :staged-edits="true"
                                              @do-save="onEditRisk">
                        </create-new-risk-form>
                    </v-col>

                    <v-col cols="4">
                        <v-card class="mb-4">
                            <v-card-title>
                                Related Nodes
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-list two-line>
                                <v-list-item v-for="(item, index) in fullRiskData.Nodes"
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
                                Related Controls
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-list two-line>
                                <v-list-item v-for="(item, index) in fullRiskData.Controls"
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
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import { FullRiskData } from '../../../ts/risks'
import { getSingleRisk, TSingleRiskInput, TSingleRiskOutput} from '../../../ts/api/apiRisks'
import CreateNewRiskForm from './CreateNewRiskForm.vue'

export default Vue.extend({
    data: () => ({
        expandDescription: false,
        ready: false,
        fullRiskData: Object() as FullRiskData
    }),
    methods: {
        onEditRisk(risk : ProcessFlowRisk) {
            //@ts-ignore
            this.$root.fullRiskData.Risk.Name = risk.Name
            //@ts-ignore
            this.$root.fullRiskData.Risk.Description = risk.Description

            Vue.nextTick(() => {
                //@ts-ignore
                this.$refs.editRisk.clearForm()
            })
        },
        refreshRiskData() {
            let data = window.location.pathname.split('/')
            let riskId = Number(data[data.length - 1])

            getSingleRisk(<TSingleRiskInput>{
                //@ts-ignore
                csrf: this.$root.csrf,
                riskId: riskId,
            }).then((resp : TSingleRiskOutput) => {
                this.fullRiskData = resp.data
                this.ready = true

                Vue.nextTick(() => {
                    //@ts-ignore
                    this.$refs.editRisk.clearForm()
                })
            }).catch((err : any) => {
                window.location.replace('/404')
            })
        }
    },
    components: {
        CreateNewRiskForm
    },
    mounted() {
        this.refreshRiskData()
    }
})

</script>
