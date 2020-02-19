<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content id="reqTitle">
                    <v-list-item-title class="title">
                        SQL Query Request: {{ currentRequest.Name }}
                    </v-list-item-title>

                    <v-list-item-subtitle v-if="!!currentMetadata">
                        Relevant Query: {{ currentMetadata.Name }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-list-item-action>
                    <v-dialog v-model="showHideDelete"
                              persistent
                              max-width="40%"
                    >
                        <template v-slot:activator="{ on }">
                            <v-btn color="error" v-on="on">
                                Delete
                            </v-btn>
                        </template>

                        <generic-delete-confirmation-form
                            item-name="SQL query requests"
                            :items-to-delete="[currentRequest.Name]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete">
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item-action>

                <v-spacer></v-spacer>

                <v-list-item-action v-if="!currentApproval">
                    <v-dialog persistent max-width="40%" v-model="showHideDenyReason">
                        <template v-slot:activator="{on}">
                            <v-btn 
                                color="error"
                                v-on="on"
                            >
                                Deny
                            </v-btn>
                        </template>

                        <v-card>
                            <v-card-title>
                                Denial Reason
                            </v-card-title>
                            <v-divider></v-divider>

                            <div class="ma-4">
                                <v-textarea v-model="reason"
                                            label="Reason"
                                            filled
                                            hide-details
                                ></v-textarea> 
                            </div>

                            <v-card-actions>
                                <v-btn
                                    color="error"
                                    @click="showHideDenyReason = false"
                                >
                                    Cancel
                                </v-btn>
                                <v-spacer></v-spacer>
                                <v-btn
                                    color="success"
                                    @click="approveDeny(false, reason)"
                                >
                                    Save
                                </v-btn>
                            </v-card-actions>
                        </v-card>


                    </v-dialog>
                </v-list-item-action>

                <v-list-item-action v-if="!currentApproval">
                    <v-btn 
                        color="success"
                        @click="approveDeny(true)"
                    >
                        Approve
                    </v-btn>
                </v-list-item-action>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="6">
                        <create-new-sql-request-form
                            :reference-request="currentRequest"
                            :force-query-id="currentRequest.QueryId"
                            edit-mode
                            @do-save="onEditRequest"
                            class="mb-4"
                        >
                        </create-new-sql-request-form>

                        <v-card>
                            <v-card-title>
                                Query
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-row v-if="!currentQuery" align="center" justify="center">
                                <v-progress-circular indeterminate size="64"></v-progress-circular>
                            </v-row>

                            <sql-text-area
                                :value="currentQuery.Query"
                                readonly
                                v-else
                            >
                            </sql-text-area>
                        </v-card>

                    </v-col>

                    <v-col cols="6">
                        <v-card class="mb-4">
                            <v-card-title>
                                Approval Status
                            </v-card-title>
                            <v-divider></v-divider>

                            <div class="px-4">
                                <p class="ma-0 py-4">
                                    <span class="font-weight-bold">
                                        Status:
                                    </span>

                                    <span v-if="!currentApproval">
                                        Pending

                                        <v-icon
                                            small
                                            color="warning"
                                        >
                                            mdi-help-circle
                                        </v-icon>
                                    </span>

                                    <span v-else-if="currentApproval.Response">
                                        Approved

                                        <v-icon
                                            small
                                            color="success"
                                        >
                                            mdi-check
                                        </v-icon>
                                    </span>

                                    <span v-else>
                                        Denied

                                        <v-icon
                                            small
                                            color="error"
                                        >
                                            mdi-cancel
                                        </v-icon>
                                    </span>
                                </p>

                                <div v-if="!!currentApproval">
                                    <user-search-form-component
                                        label="Responder"
                                        v-bind:user="responderUser"
                                        readonly
                                    ></user-search-form-component>

                                    <p>
                                        <span class="font-weight-bold">
                                            Responded At:
                                        </span>
                                        {{ responseTime }}
                                    </p>

                                    <p v-if="!currentApproval.Response">
                                        <span class="font-weight-bold">
                                            Reason:
                                        </span>
                                        <pre class="pb-4">{{ currentApproval.Reason }}</pre>
                                    </p>
                                </div>
                            </div>
                        </v-card>

                        <v-card>
                            <v-card-title>
                                Comments
                            </v-card-title>
                            <v-divider></v-divider>

                            <comment-manager
                                :params="commentParams"
                            ></comment-manager>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgDocRequestsUrl } from '../../../ts/url'
import { standardFormatTime } from '../../../ts/time'
import { DbSqlQueryRequest, DbSqlQueryRequestApproval, DbSqlQuery, DbSqlQueryMetadata } from '../../../ts/sql'
import { 
    getSqlRequest, TGetSqlRequestOutput,
    modifyStatusSqlRequest, TModifyStatusSqlRequestOutput,
    deleteSqlRequest,
} from '../../../ts/api/apiSqlRequests'
import { getSqlQuery, TGetSqlQueryOutput } from '../../../ts/api/apiSqlQueries'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import CommentManager from '../../generic/CommentManager.vue'
import CreateNewSqlRequestForm from './CreateNewSqlRequestForm.vue'
import SqlTextArea from '../../generic/SqlTextArea.vue'
import MetadataStore from '../../../ts/metadata'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'

@Component({
    components: {
        GenericDeleteConfirmationForm,
        CommentManager,
        CreateNewSqlRequestForm,
        UserSearchFormComponent,
        SqlTextArea
    }
})
export default class FullEditSqlRequestComponent extends Vue {
    showHideDelete: boolean = false
    showHideDenyReason : boolean = false

    reason : string = ""

    currentRequest: DbSqlQueryRequest | null = null
    currentApproval : DbSqlQueryRequestApproval | null = null
    currentQuery : DbSqlQuery | null = null
    currentMetadata : DbSqlQueryMetadata | null = null

    get ready() : boolean {
        return !!this.currentRequest
    }

    get responderUser() : User | null {
        if (!this.currentApproval) {
            return null
        }
        return MetadataStore.getters.getUser(this.currentApproval.ResponsderUserId)
    }

    get responseTime() : string {
        if (!this.currentApproval) {
            return ""
        }

        return standardFormatTime(this.currentApproval.ResponseTime) 
    }

    get commentParams() : Object {
        return {
            sqlRequestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }
    }

    refreshQuery() {
        if (!this.currentRequest) {
            return
        }

        getSqlQuery({
            metadataId: -1,
            orgId: PageParamsStore.state.organization!.Id,
            queryId: this.currentRequest.QueryId,
        }).then((resp : TGetSqlQueryOutput) => {
            this.currentQuery = resp.data.Queries[0]
            this.currentMetadata = resp.data.Metadata
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

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getSqlRequest({
            requestId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSqlRequestOutput) => {
            this.currentRequest = resp.data.Request
            this.currentApproval = resp.data.Approval
            this.refreshQuery()
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

    onDelete() {
        deleteSqlRequest({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            window.location.replace(createOrgDocRequestsUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onEditRequest(req : DbSqlQueryRequest) {
        this.currentRequest = req
    }

    approveDeny(status : boolean, reason : string = "") {
        modifyStatusSqlRequest({
            requestId: this.currentRequest!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            approve: status,
            reason: reason,
        }).then((resp : TModifyStatusSqlRequestOutput) => {
            this.reason = ""
            this.showHideDenyReason = false
            this.currentApproval = resp.data
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

<style scoped>

#reqTitle {
    flex: none !important;
}

</style>
