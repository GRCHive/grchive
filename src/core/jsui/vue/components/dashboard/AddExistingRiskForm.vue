<template>

<v-card>
    <v-card-title>
        Add Existing Risk
    </v-card-title>
    <v-divider></v-divider>

    <section class="ma-2">
        <v-list two-line>
            <v-list-item-group multiple v-model="selectedRisks">
                <v-list-item v-for="(item, index) in allRisks"
                             :key="index" class="pa-1" :value="item">
                    <template v-slot:default="{active, toggle}"
                              v-if="!preselectedSet.has(item.Id)">
                        <v-list-item-action class="ma-1">
                            <v-checkbox :input-value="active"
                                        @true-value="item"
                                        @click="toggle">
                            </v-checkbox>
                        </v-list-item-action>

                        <v-list-item-content>
                            <v-list-item-title>
                                {{ item.Name }}
                            </v-list-item-title>

                            <v-list-item-subtitle>
                                {{ item.Description }}
                            </v-list-item-subtitle>
                        </v-list-item-content>
                    </template>
                </v-list-item>
            </v-list-item-group>
        </v-list>
    </section>
    <v-divider></v-divider>
    <v-card-actions>
        <v-btn color="error" @click="onCancel">
            Cancel
        </v-btn>
        <div class="flex-grow-1"></div>
        <v-btn color="success" @click="onSelect">
            Select
        </v-btn>
    </v-card-actions>

</v-card>
    
</template>

<script lang="ts">

import Vue from 'vue'
import VueSetup from "../../../ts/vueSetup"
import { contactUsUrl } from "../../../ts/url"
import { newRisk } from "../../../ts/api/apiRisks"

export default Vue.extend({
    props: {
        preselectedRisks: Array
    },
    data: () => ({
        selectedRisks : [] as ProcessFlowRisk[]
    }),
    computed: {
        allRisks() : ProcessFlowRisk[] {
            let risks : ProcessFlowRisk[] = []
            for (let riskId of VueSetup.store.state.currentProcessFlowFullData.RiskKeys) {
                risks.push(VueSetup.store.state.currentProcessFlowFullData.Risks[riskId])
            }
            return risks
        },
        preselectedSet() : Set<number> {
            let newSet = new Set<number>()
            for (let r of this.preselectedRisks) {
                newSet.add((<ProcessFlowRisk>r).Id)
            }
            return newSet
        }
    },
    methods: {
        onCancel() {
            this.$emit('do-cancel')
        },
        onSelect() {
            this.$emit('do-select', this.selectedRisks)
        }
    }
})

</script>
