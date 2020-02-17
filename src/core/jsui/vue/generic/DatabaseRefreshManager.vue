<template>
    <div>
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-list-item>
                <v-list-item-content>
                    <v-select
                        :value="selectedRefresh"
                        @input="initialSelectRefresh"
                        label="Versions"
                        filled
                        :items="refreshItems"
                        hide-details
                        dense
                    >
                    </v-select>
                </v-list-item-content>

                <v-list-item-action>
                    <v-dialog
                        v-model="showHideDelete"
                        persistent
                        max-width="40%"
                    >
                        <template v-slot:activator="{on}">
                            <v-btn
                                color="error"
                                icon
                                v-on="on"
                                :disabled="!selectedRefresh"
                            >
                                <v-icon>
                                    mdi-delete
                                </v-icon>
                            </v-btn>
                        </template>

                        <generic-delete-confirmation-form
                            item-name="refreshes"
                            :items-to-delete="[currentRefreshName]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="deleteCurrentRefresh">
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item-action>

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
                :fn-mode="fnMode"
            >
            </database-refresh-viewer>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { PageParamsStore } from '../../ts/pageParams'
import {
    allSqlRefresh, TAllSqlRefreshOutput,
    newSqlRefresh, TNewSqlRefreshOutput,
    getSqlRefresh, TGetSqlRefreshOutput,
    deleteSqlRefresh,
} from '../../ts/api/apiSqlRefresh'
import { contactUsUrl } from '../../ts/url'
import {
    DbRefresh,
    dbRefreshIdentifier,
    DbSchema
} from '../../ts/sql'
import { standardFormatTime } from '../../ts/time'
import { DatabaseStore, getStoreForDatabase } from '../../ts/vuex/databaseStore'
import DatabaseRefreshViewer from './DatabaseRefreshViewer.vue'
import GenericDeleteConfirmationForm from '../components/dashboard/GenericDeleteConfirmationForm.vue'

const Props = Vue.extend({
    props: {
        dbId: Number,
        fnMode: {
            type: Boolean,
            default: false,
        }
    },
})

@Component({
    components: {
        DatabaseRefreshViewer,
        GenericDeleteConfirmationForm
    }
})
export default class DatabaseRefreshManager extends Props {
    store : DatabaseStore = getStoreForDatabase(this.dbId)
    showHideDelete : boolean = false

    get selectedRefresh() : DbRefresh | null {
        return this.store.state.selectedRefresh
    }

    get refreshInProgress() : boolean {
        return this.store.state.isPollingRefresh
    }

    get currentRefreshName() : string {
        return dbRefreshIdentifier(this.selectedRefresh!)
    }

    get refreshItems() : any[] {
        if (!this.store.state.allRefreshes) {
            return []
        }
        return this.store.state.allRefreshes.map((ele : DbRefresh, idx : number) => ({
            text: `#${this.store.state.allRefreshes!.length - idx} :: ${dbRefreshIdentifier(ele)}`,
            value: ele,
        }))
    }

    get isLoading() : boolean {
        return (this.store.state.allRefreshes == null)
    }

    initialSelectRefresh(ref : DbRefresh) {
        this.store.dispatch('requestSetNewRefresh', ref)
    }

    requestRefreshSchema() {
        newSqlRefresh({
            dbId: this.dbId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TNewSqlRefreshOutput) => {
            this.store.commit('addRefresh', resp.data)
            this.store.dispatch('requestSetNewRefresh', resp.data)
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

    deleteCurrentRefresh() {
        if (!this.store.state.selectedRefresh) {
            return
        }

        let refresh : DbRefresh = this.store.state.selectedRefresh!
        deleteSqlRefresh({
            refreshId: refresh.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.store.dispatch('deleteRefresh', refresh)
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
}

</script>
