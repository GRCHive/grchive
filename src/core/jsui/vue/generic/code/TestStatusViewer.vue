<template>
    <div>
        <v-overlay :value="exporting">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>


        <div v-if="isLoading">
            <v-progress-circular
                indeterminate
                size="16"
            ></v-progress-circular>
            <span>Loading...</span>
        </div>

        <div v-else>
            <v-icon
                :color="iconColor"
                small
            >
                {{ iconMdi }}
            </v-icon>

            <span>{{ summary.SuccessfulTests }} / {{ summary.TotalTests }}</span>

            <v-btn
                @click="exportToExcel"
                class="ml-2"
                color="primary"
                v-if="showExport"
                small
            >
                Export
            </v-btn>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    getCodeRunTest, TGetCodeRunTestOutput,
    exportTests, TExportTestOutput
} from '../../../ts/api/apiTests'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { CodeRunTestSummary } from '../../../ts/tests'
import { saveAs } from 'file-saver'

const Props = Vue.extend({
    props: {
        runId: Number,
        showExport: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class TestStatusViewer extends Props {
    summary : CodeRunTestSummary | null = null
    exporting : boolean = false

    get isLoading() : boolean {
        return !this.summary
    }

    get isSuccess() : boolean {
        return this.summary!.SuccessfulTests == this.summary!.TotalTests
    }

    get iconColor() : string {
        if (this.isSuccess) {
            return 'success'
        } else {
            return 'error'
        }
    }

    get iconMdi() : string {
        if (this.isSuccess) {
            return 'mdi-check-circle'
        } else {
            return 'mdi-close-circle'
        }
    }

    exportToExcel() {
        this.exporting = true
        exportTests({
            orgId: PageParamsStore.state.organization!.Id,
            runId: this.runId,
        }).then((resp : TExportTestOutput) => {
            saveAs(resp.data, "tests.xlsx")
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => {
            this.exporting = false
        })
    }

    refreshData() {
        getCodeRunTest({
            orgId: PageParamsStore.state.organization!.Id,
            runId: this.runId,
            summary: true,
        }).then((resp : TGetCodeRunTestOutput) => {
            this.summary = resp.data
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

    mounted() {
        this.refreshData()
    }
}

</script>
