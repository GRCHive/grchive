<template>
    <div>
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>
        
        <div v-else>
            <v-row>
                <v-col cols="6">
                    <v-list-item>
                        <v-list-item-content id="querySelector">
                            <v-select
                                label="Queries"
                                filled
                                hide-details
                                dense
                                :items="metadataItems"
                                :value="currentMetadata"
                                @input="selectMetadata"
                            >
                            </v-select>
                        </v-list-item-content>

                        <v-list-item-content class="ml-2" id="versionSelector">
                            <v-select
                                label="Versions"
                                filled
                                hide-details
                                dense
                                :items="queryItems"
                                :value="currentVersion"
                                @input="selectVersion"
                            >
                            </v-select>
                        </v-list-item-content>

                        <v-list-item-action>
                            <v-btn
                                color="error"
                                icon
                                @click="deleteQuery"
                                :disabled="!canDelete"
                            >
                                <v-icon>
                                    mdi-delete
                                </v-icon>
                            </v-btn>
                        </v-list-item-action>

                        <v-spacer></v-spacer>

                        <v-list-item-action>
                            <v-dialog
                                v-model="showHideNewQuery"
                                persistent
                                max-width="40%"
                            >
                                <template v-slot:activator="{ on }">
                                    <v-btn
                                        color="primary"
                                        icon
                                        v-on="on"
                                    >
                                        <v-icon>
                                            mdi-plus
                                        </v-icon>
                                    </v-btn>
                                </template>

                                <create-new-sql-query-form
                                    :db-id="dbId"
                                    @do-cancel="showHideNewQuery = false"
                                    @do-save="onSaveNewQuery"
                                >
                                </create-new-sql-query-form>
                            </v-dialog>
                        </v-list-item-action>
                    </v-list-item>

                    <div class="px-4" v-if="!!currentMetadata">
                        <create-new-sql-query-form
                            :db-id="dbId"
                            :enable-query="false"
                            :reference-metadata="currentMetadata"
                            edit-mode
                            @do-save="onEditQueryMetadata"
                        >
                        </create-new-sql-query-form>

                        <v-card class="mt-4" v-if="!!currentVersion">
                            <v-card-title class="pl-3">
                                Version Information 
                            </v-card-title>
                            <v-divider></v-divider>

                            <v-form class="ma-4">
                                <user-search-form-component
                                    :user="queryOwner"
                                    label="Upload User"
                                    readonly
                                ></user-search-form-component>

                                <v-text-field
                                    :value="queryUploadTime"
                                    label="Upload Time"
                                    readonly
                                >
                                </v-text-field>
                            </v-form>
                        </v-card>
                    </div>
                </v-col>

                <v-col cols="6">
                    <div v-if="!!currentVersion">
                        <v-list-item class="pa-0">
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
                    </div>
                </v-col>
            </v-row>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import CreateNewSqlQueryForm from '../components/dashboard/CreateNewSqlQueryForm.vue'
import UserSearchFormComponent from './UserSearchFormComponent.vue'
import { DbSqlQueryMetadata , DbSqlQuery } from '../../ts/sql'
import {
    TAllSqlQueryOutput, allSqlQuery,
    TGetSqlQueryOutput, getSqlQuery,
    TUpdateSqlQueryOutput, updateSqlQuery,
    deleteSqlQuery
} from '../../ts/api/apiSqlQueries'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import MetadataStore from '../../ts/metadata'
import SqlTextArea from './SqlTextArea.vue'
import { standardFormatTime } from '../../ts/time'

const Props = Vue.extend({
    props: {
        dbId: Number,
    }
})

@Component({
    components: {
        CreateNewSqlQueryForm,
        UserSearchFormComponent,
        SqlTextArea
    }
})
export default class DatabaseQueryManager extends Props {
    allMetadata : DbSqlQueryMetadata[] | null = null
    allVersions : DbSqlQuery[] | null = null

    currentMetadata : DbSqlQueryMetadata | null = null
    currentVersion : DbSqlQuery | null = null

    editableQuery : string = ""
    canEditQuery : boolean = false
    queryKey : number = 0
    queryRunning : boolean = false

    showHideNewQuery : boolean = false

    get metadataItems() : any[] {
        if (!this.allMetadata) {
            return []
        }
        return this.allMetadata.map((ele : DbSqlQueryMetadata) => ({
            text: ele.Name,
            value: ele,
        }))
    }

    get queryItems() : any[] {
        if (!this.allVersions) {
            return []
        }
        return this.allVersions.map((ele : DbSqlQuery) => ({
            text: `v${ele.Version}`,
            value: ele,
        }))
    }

    get canDelete() : boolean {
        if (!this.allMetadata) {
            return false
        }
        return this.allMetadata!.length > 0 && !!this.currentMetadata
    }

    get isLoading() : boolean {
        return !this.allMetadata || !this.allVersions
    }

    get queryUploadTime() : string {
        if (!this.currentVersion) {
            return "N/A"
        }
        return standardFormatTime(this.currentVersion.UploadTime)
    }

    get queryOwner() : User | null {
        if (!this.currentVersion) {
            return null
        }
        return MetadataStore.getters.getUser(this.currentVersion!.UploadUserId)
    }

    selectVersion(version : DbSqlQuery) {
        this.currentVersion = version
        this.cancelEditQuery()
    }

    refreshVersions() {
        getSqlQuery({
            metadataId: this.currentMetadata!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetSqlQueryOutput) => {
            this.allVersions = resp.data
            if (this.allVersions!.length > 0) {
                this.selectVersion(this.allVersions![0])
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

    selectMetadata(meta : DbSqlQueryMetadata | null, refresh: boolean = true) {
        this.currentMetadata = meta

        if (!!this.currentMetadata) {
            if (refresh) {
                this.refreshVersions()
            }
        } else {
            this.allVersions = []
            this.currentVersion = null
        }
    }

    refreshMetadata() {
        allSqlQuery({
            dbId: this.dbId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllSqlQueryOutput) => {
            this.allMetadata = resp.data
            if (this.allMetadata!.length > 0) {
                this.selectMetadata(this.allMetadata![0])
            } else {
                this.selectMetadata(null)
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

    mounted() {
        this.refreshMetadata()
    }

    onSaveNewQuery(metadata: DbSqlQueryMetadata, query : DbSqlQuery) {
        this.showHideNewQuery = false
        Vue.nextTick(() => {
            this.allMetadata!.unshift(metadata)
            this.allVersions!.unshift(query)

            this.selectMetadata(this.allMetadata![0], false)
            this.selectVersion(this.allVersions![0])
        })
    }

    onEditQueryMetadata(metadata: DbSqlQueryMetadata) {
        let idx = this.allMetadata!.findIndex((ele : DbSqlQueryMetadata) => ele.Id == metadata.Id)
        if (idx == -1) {
            return
        }
        Vue.set(this.allMetadata!, idx, metadata)
        if (this.currentMetadata!.Id == metadata.Id) {
            this.currentMetadata = metadata
        }
    }

    cancelEditQuery() {
        this.editableQuery = this.currentVersion!.Query
        this.queryKey += 1
        this.canEditQuery = false
    }

    saveEditQuery() {
        updateSqlQuery({
            orgId: PageParamsStore.state.organization!.Id,
            metadataId: this.currentMetadata!.Id,
            query: {
                query: this.editableQuery,
                uploadUserId: PageParamsStore.state.user!.Id,
            }
        }).then((resp : TUpdateSqlQueryOutput) => {
            if (resp.data.Query) {
                this.allVersions!.unshift(resp.data.Query)
                this.selectVersion(this.allVersions![0])
                this.canEditQuery = false
            } else {
                throw "Did not receive a response query."
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

    deleteQuery() {
        let id = this.currentMetadata!.Id

        deleteSqlQuery({
            orgId: PageParamsStore.state.organization!.Id,
            metadataId: id,
        }).then(() => {
            let idx = this.allMetadata!.findIndex((ele : DbSqlQueryMetadata) => ele.Id == id)
            if (idx == -1) {
                return
            }

            this.allMetadata!.splice(idx, 1)
            if (this.allMetadata!.length > 0) {
                this.selectMetadata(this.allMetadata![0])
            } else {
                this.selectMetadata(null)
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

    resetQuery() {
        this.queryKey += 1
    }

    runQuery() {
        this.queryRunning = true
    }
}

</script>

<style scoped>

#querySelector {
    flex: 3 1;
}

#versionSelector {
    flex: 1 1;
}

</style>
