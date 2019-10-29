<template>
    <section class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Risks
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
                <v-btn class="primary">
                    Add
                </v-btn>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>
        <v-card
            v-for="(item, index) in filteredRisks"
            :key="index"
            class="my-2"
        >
            <v-list-item two-line :href="generateRiskUrl(item.Id)">
                <v-list-item-content>
                    <v-list-item-title v-html="highlightText(item.Name)">
                    </v-list-item-title>
                    <v-list-item-subtitle v-html="highlightText(item.Description)">
                    </v-list-item-subtitle>
                </v-list-item-content>
                <v-spacer></v-spacer>
                <v-list-item-action>
                    <v-btn icon @click.stop @mousedown.stop>
                        <v-icon>mdi-delete</v-icon>
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
        </v-card>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import { getAllRisks, TAllRiskInput, TAllRiskOutput } from '../../../ts/api/apiRisks'
import { contactUsUrl, createRiskUrl } from '../../../ts/url'
import { replaceWithMark, sanitizeTextForHTML } from '../../../ts/text'

export default Vue.extend({
    data : () => ({
        allRisks: [] as ProcessFlowRisk[],
        filterText : ""
    }),
    computed: {
        filter() : (a : ProcessFlowRisk) => boolean {
            const filterText = this.filterText.trim()
            return (ele : ProcessFlowRisk) : boolean => {
                return ele.Name.toLocaleLowerCase().includes(filterText.toLocaleLowerCase()) ||
                    ele.Description.toLocaleLowerCase().includes(filterText.toLocaleLowerCase())
            }
        },
        filteredRisks() : ProcessFlowRisk[] {
            return this.allRisks.filter(this.filter)
        }
    },
    methods: {
        highlightText(input : string) : string {
            const safeInput = sanitizeTextForHTML(input)
            const useFilter = this.filterText.trim()
            if (useFilter.length == 0) {
                return safeInput
            }
            return replaceWithMark(
                safeInput,
                sanitizeTextForHTML(useFilter))
        },
        generateRiskUrl(riskId : number) : string {
            //@ts-ignore
            return createRiskUrl(this.$root.orgGroupId, riskId)
        },
        refreshRisks() {
            getAllRisks(<TAllRiskInput>{
                //@ts-ignore
                csrf: this.$root.csrf
            }).then((resp : TAllRiskOutput) => {
                this.allRisks = resp.data
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
    },
    mounted() {
        this.refreshRisks()
    }
})
</script>
