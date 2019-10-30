<template>
    <section class="ma-4">
        <!--
        <v-dialog v-model="showHideDeleteRisk" persistent max-width="40%">
            <generic-delete-confirmation-form
                item-name="risks"
                :items-to-delete="currentRisksToDelete"
                v-on:do-cancel="showHideDeleteRisk = false"
                v-on:do-delete="deleteSelectedRisks"
                :use-global-deletion="true"
                :force-global-deletion="true">
            </generic-delete-confirmation-form>
        </v-dialog>
        -->

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
                <!--
                <v-dialog v-model="showHideCreateNewRisk" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn class="primary" v-on="on">
                            Create New
                        </v-btn>
                    </template>
                    <create-new-risk-form
                        :node-id="-1"
                        @do-save="saveNewRisk"
                        @do-cancel="cancelNewRisk">
                    </create-new-risk-form>
                </v-dialog>
                -->
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>
        <v-card
            v-for="(item, index) in filteredRisks"
            :key="index"
            class="my-2"
        >
            <v-list-item two-line @click="goToControl(item.Id)">
                <v-list-item-content>
                    <v-list-item-title v-html="highlightText(item.Name)">
                    </v-list-item-title>
                    <v-list-item-subtitle v-html="highlightText(item.Description)">
                    </v-list-item-subtitle>
                </v-list-item-content>
                <v-spacer></v-spacer>
                <v-list-item-action>
                    <v-btn icon @click.stop="doDeleteControl(item)" @mousedown.stop @mouseup.stop>
                        <v-icon>mdi-delete</v-icon>
                    </v-btn>
                </v-list-item-action>
            </v-list-item>
        </v-card>
    </section>
</template>

<script lang="ts">

import Vue from 'vue'
import { getAllControls, TAllControlInput, TAllControlOutput} from '../../../ts/api/apiControls'
import { replaceWithMark, sanitizeTextForHTML } from '../../../ts/text'
import { contactUsUrl } from '../../../ts/url'

export default Vue.extend({
    data : () => ({
        allControls: [] as ProcessFlowControl[],
        filterText : "",
    }),
    components: {
    },
    computed: {
        filter() : (a : ProcessFlowControl) => boolean {
            const filterText = this.filterText.trim()
            return (ele : ProcessFlowControl) : boolean => {
                return ele.Name.toLocaleLowerCase().includes(filterText.toLocaleLowerCase()) ||
                    ele.Description.toLocaleLowerCase().includes(filterText.toLocaleLowerCase())
            }
        },
        filteredRisks() : ProcessFlowRisk[] {
            return this.allControls.filter(this.filter)
        },
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
        refreshControls() {
            getAllControls(<TAllControlInput>{
                //@ts-ignore
                csrf: this.$root.csrf
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
        goToControl(controlId : number) {
        },
        doDeleteControl(control : ProcessFlowControl) {
        },
    },
    mounted() {
        this.refreshControls()
    }
})
</script>
