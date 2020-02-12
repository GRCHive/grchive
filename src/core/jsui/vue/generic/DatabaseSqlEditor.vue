<template>
    <v-container fluid>
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-tabs vertical>
                <v-tab>Schemas</v-tab>
                <v-tab-item>
                    <v-list-item>
                        <v-list-item-content>
                            <v-select
                                v-model="selectedRefresh"
                                label="Versions"
                                filled
                                :items="refreshItems"
                                hide-details
                                dense
                            >
                            </v-select>
                        </v-list-item-content>

                        <v-spacer></v-spacer>

                        <v-list-item-action>
                            <v-btn
                                color="primary"
                                icon
                                :loading="refreshInProgress"
                                @click="requestRefreshSchema"
                            >
                                <v-icon>
                                    mdi-refresh
                                </v-icon>
                            </v-btn>
                        </v-list-item-action>
                    </v-list-item>

                    <database-refresh-viewer
                        class="px-4"
                        :refresh="selectedRefresh"
                    >
                    </database-refresh-viewer>
                </v-tab-item>

                <v-tab>Queries</v-tab>
                <v-tab-item>
                </v-tab-item>
            </v-tabs>

            <div id="queryOutput">
            </div>
        </div>
    </v-container>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { PageParamsStore } from '../../ts/pageParams'
import {
    allSqlRefresh, TAllSqlRefreshOutput,
    newSqlRefresh, TNewSqlRefreshOutput,
    getSqlRefresh, TGetSqlRefreshOutput,
} from '../../ts/api/apiSqlRefresh'
import { contactUsUrl } from '../../ts/url'
import {
    DbRefresh,
    dbRefreshIdentifier,
    DbSchema
} from '../../ts/sql'
import { standardFormatTime } from '../../ts/time'
import DatabaseRefreshViewer from './DatabaseRefreshViewer.vue'

const Props = Vue.extend({
    props: {
        dbId: Number,
    }
})

const refreshIntervalSeconds : number = 15

@Component({
    components: {
        DatabaseRefreshViewer,
    }
})
export default class DatabaseSqlEditor extends Props {
    schemaRefreshes : DbRefresh[] | null = null
    selectedRefresh : DbRefresh | null = null
    refreshInProgress : boolean = false

    get refreshItems() : any[] {
        if (!this.schemaRefreshes) {
            return []
        }
        return this.schemaRefreshes.map((ele : DbRefresh, idx : number) => ({
            text: `#${this.schemaRefreshes!.length - idx} :: ${dbRefreshIdentifier(ele)}`,
            value: ele,
        }))
    }

    get isLoading() : boolean {
        return (this.schemaRefreshes == null)
    }

    initialSelectRefresh(ref : DbRefresh) {
        this.selectedRefresh = ref
        if (!this.selectedRefresh.RefreshFinishTime) {
            this.startRefreshPoll(this.selectedRefresh.Id)
        }
    }

    refreshData() {
        allSqlRefresh({
            dbId: this.dbId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllSqlRefreshOutput) => {
            this.schemaRefreshes = resp.data
            if (this.schemaRefreshes!.length > 0) {
                this.initialSelectRefresh(this.schemaRefreshes![0])
            }
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

    requestRefreshSchema() {
        this.refreshInProgress = true
        newSqlRefresh({
            dbId: this.dbId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TNewSqlRefreshOutput) => {
            this.schemaRefreshes!.unshift(resp.data)
            this.initialSelectRefresh(resp.data)
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

    startRefreshPoll(refreshId : number) {
        this.refreshInProgress = true

        // Silently ignore any polling errors.
        let intervalId = setInterval(() => {
            getSqlRefresh({
                refreshId : refreshId,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TNewSqlRefreshOutput) => {
                if (!resp.data.RefreshFinishTime) {
                    return
                }

                this.refreshInProgress = false
                let idx = this.schemaRefreshes!.findIndex((ele : DbRefresh) => ele.Id == refreshId)
                if (idx == -1) {
                    return
                }

                this.schemaRefreshes![idx] = resp.data
                if (!!this.selectedRefresh && this.selectedRefresh.Id == refreshId) {
                    this.selectedRefresh = resp.data
                }

            }).catch((err : any) => {})
        }, refreshIntervalSeconds * 1000)
    }

    mounted() {
        this.refreshData()
    }
}

</script>
