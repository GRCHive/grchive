<template>
    <div v-if="!!query">
        <v-dialog persistent max-width="40%" v-model="showNewRequest">
            <create-new-sql-request-form
                :force-query-id="query.Id"
                @do-cancel="showNewRequest = false"
                @do-save="onNewRequest"
            >
            </create-new-sql-request-form>
        </v-dialog>

        <v-list-item class="pa-0">
            <v-list-item-content>
                <v-text-field
                    v-model="runCode"
                    label="Run Code"
                    dense
                    hide-details
                    clearable
                    outlined
                >
                </v-text-field>
            </v-list-item-content>

            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-btn
                    color="warning"
                    icon
                    x-small
                    @click="resetQuery"
                >
                    <v-icon>mdi-bug</v-icon>
                </v-btn>
            </v-list-item-action>

            <v-list-item-action>
                <v-btn
                    color="primary"
                    icon
                    x-small
                    @click="runQuery"
                    :loading="queryRunning"
                >
                    <v-icon>mdi-play</v-icon>
                </v-btn>
            </v-list-item-action>
        </v-list-item>

        <sql-text-area
            v-model="editableQuery"
            :readonly="!canEditQuery"
            :key="`MANAGER-${queryKey}`"
        >
        </sql-text-area>

        <v-list-item class="pa-0">
            <v-list-item-action>
                <v-btn
                    color="error"
                    @click="cancelEditQuery"
                    v-if="canEditQuery"
                >
                    Cancel
                </v-btn>
            </v-list-item-action>
            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-btn
                    color="success"
                    @click="saveEditQuery"
                    v-if="canEditQuery"
                >
                    Save
                </v-btn>
            </v-list-item-action>

            <v-list-item-action>
                <v-btn
                    color="success"
                    @click="canEditQuery = true"
                    v-if="!canEditQuery"
                >
                    Edit
                </v-btn>
            </v-list-item-action>
        </v-list-item>

        <v-divider></v-divider>

        <sql-result-display
            v-if="!!currentResult"
            :result="currentResult"
        >
        </sql-result-display>
    </div>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { DbSqlQuery, DbSqlQueryRequest } from '../../ts/sql'
import {
    TRunSqlQueryOutput, runSqlQuery,
} from '../../ts/api/apiSqlQueries'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import SqlTextArea from './SqlTextArea.vue'
import SqlResultDisplay from './SqlResultDisplay.vue'
import CreateNewSqlRequestForm from '../components/dashboard/CreateNewSqlRequestForm.vue'

const Props = Vue.extend({
    props : {
        query: {
            type: Object,
            default: () => null as DbSqlQuery | null
        }
    }
})

@Component({
    components: {
        SqlTextArea,
        SqlResultDisplay,
        CreateNewSqlRequestForm,
    }
})
export default class DatabaseQueryRunner extends Props {
    editableQuery : string = ""
    runCode : string = ""

    showNewRequest: boolean = false

    canEditQuery : boolean = false
    queryKey : number = 0
    queryRunning : boolean = false
    versionIdToResult : Record<number, TRunSqlQueryOutput> = Object()

    get currentResult() : TRunSqlQueryOutput | null {
        if (!this.query) {
            return null
        }

        if (!(this.query.Id in this.versionIdToResult)) {
            return null
        }

        return this.versionIdToResult[this.query.Id]
    }

    cancelEditQuery() {
        this.editableQuery = this.query!.Query
        this.queryKey += 1
        this.canEditQuery = false
    }

    saveEditQuery() {
        this.$emit('on-new-version', this.editableQuery)
        Vue.nextTick(() => {
            this.cancelEditQuery()
        })
    }

    resetQuery() {
        this.queryKey += 1
    }

    doExecute(runCode : string) {
        this.queryRunning = true
        runSqlQuery({
            queryId: this.query!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            runCode: runCode,
        }).then((resp : TRunSqlQueryOutput) => {
            this.queryRunning = false
            Vue.set(this.versionIdToResult, this.query!.Id, resp)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
            this.queryRunning = false
        })
    }

    startRequest() {
        this.showNewRequest = true
    }

    onNewRequest(req : DbSqlQueryRequest) {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Your request has been saved.",
            false,
            "Contact Us",
            contactUsUrl,
            false);
        this.showNewRequest = false
    }

    runQuery() {
        let runCode = !!this.runCode ? this.runCode.trim() : ""
        if (runCode == "") {
            this.startRequest()
        } else {
            this.doExecute(runCode)
        }
    }

    @Watch('query')
    refreshQuery() {
        this.cancelEditQuery()
    }

    mounted() {
        this.refreshQuery()
    }
}


</script>
