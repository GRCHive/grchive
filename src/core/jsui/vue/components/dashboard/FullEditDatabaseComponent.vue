<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-dialog v-model="showHideDeleteConnection"
                      persistent
                      max-width="40%"
            >
                <generic-delete-confirmation-form
                    item-name="database connections"
                    :items-to-delete="[currentDb.Name]"
                    :use-global-deletion="false"
                    @do-cancel="showHideDeleteConnection = false"
                    @do-delete="onDeleteDbConn">
                </generic-delete-confirmation-form>
            </v-dialog>

            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Database: {{ currentDb.Name }}
                    </v-list-item-title>

                    <v-list-item-subtitle>
                        {{ fullTypeString }}
                    </v-list-item-subtitle>
                </v-list-item-content>

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
                        item-name="databases"
                        :items-to-delete="[currentDb.Name]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row>
                            <v-col cols="6">
                                <create-new-database-form
                                    :edit-mode="true"
                                    :reference-db="currentDb"
                                    @do-save="onEdit">
                                </create-new-database-form>
                            </v-col>

                            <v-col cols="6">
                                <v-card class="mb-4">
                                    <v-card-title>
                                        Relevant Systems
                                        <v-spacer></v-spacer>

                                        <v-dialog persistent
                                                  max-width="40%"
                                                  v-model="showHideLinkSystem">
                                            <template v-slot:activator="{ on }">
                                                <v-btn color="primary" icon v-on="on">
                                                    <v-icon>mdi-plus</v-icon>
                                                </v-btn>
                                            </template>

                                            <v-card>
                                                <v-card-title>
                                                    Link Systems
                                                </v-card-title>
                                                <v-divider></v-divider>

                                                <systems-table :resources="allSystems"
                                                               selectable
                                                               multi
                                                               v-model="systemsToLink"
                                                ></systems-table>

                                                <v-card-actions>
                                                    <v-btn color="error" @click="showHideLinkSystem = false">
                                                        Cancel
                                                    </v-btn>
                                                    <v-spacer></v-spacer>
                                                    <v-btn color="success" @click="linkSystems">
                                                        Link
                                                    </v-btn>
                                                </v-card-actions>
                                            </v-card>
                                        </v-dialog>
                                    </v-card-title>
                                    <v-divider></v-divider>
                                    <systems-table
                                        :resources="relatedSystems"
                                        use-crud-delete
                                        @delete="onDeleteSysLink"
                                    ></systems-table>
                                </v-card>

                                <v-card class="mb-4">
                                    <v-card-title>
                                        Connection Info

                                        <v-spacer></v-spacer>
                                        <v-btn color="error"
                                               outlined
                                               fab
                                               small
                                               v-if="hasDb"
                                               @click="showHideDeleteConnection=true"
                                            >
                                            <v-icon>mdi-delete</v-icon>
                                        </v-btn>

                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <v-row align="center" justify="center" v-if="!hasDb">
                                        <v-dialog v-model="showHideNewConn"
                                                  persistent
                                                  max-width="40%"
                                        >
                                            <template v-slot:activator="{ on }">
                                                <v-btn color="primary"
                                                       v-on="on"
                                                       fab
                                                       outlined
                                                       x-large
                                                       class="my-6"
                                                       :disabled="!canConnectToDb">
                                                    <v-icon>mdi-plus</v-icon>
                                                </v-btn>
                                            </template>

                                            <create-new-db-connection-form
                                                :db-id="currentDb.Id"
                                                @do-cancel="showHideNewConn = false"
                                                @do-save="onNewDbConn">
                                            </create-new-db-connection-form>
                                        </v-dialog>
                                    </v-row>

                                    <database-connection-read-only-component v-else
                                        :connection="dbConn">
                                    </database-connection-read-only-component>
                                </v-card>

                                <v-card>
                                    <v-card-title>
                                        Deployment
                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <div v-if="!deployment">
                                        <v-row align="center" justify="center">
                                            <v-btn color="primary"
                                                   fab
                                                   outlined
                                                   x-large
                                                   class="my-6"
                                                   @click="addDeployment">
                                                <v-icon>mdi-plus</v-icon>
                                            </v-btn>
                                        </v-row>
                                    </div>

                                    <div v-else>
                                        <deployment-editor
                                            class="pa-4"
                                            v-model="deployment"
                                            :db-id="currentDb.Id"
                                        >
                                        </deployment-editor>
                                    </div>
                                </v-card>

                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab :disabled="!canConnectToDb">SQL</v-tab>
                    <v-tab-item>
                        <database-sql-editor
                            :db-id="currentDb.Id"
                        >
                        </database-sql-editor>
                    </v-tab-item>

                    <v-tab>Requests</v-tab>
                    <v-tab-item>
                        <dashboard-sql-request-list
                            class="px-2"
                            :db-id="currentDb.Id"
                        >
                        </dashboard-sql-request-list>
                    </v-tab-item>
                </v-tabs>

            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { getDatabase, TGetDatabaseOutputs } from '../../../ts/api/apiDatabases'
import { deleteDatabase, TDeleteDatabaseOutputs } from '../../../ts/api/apiDatabases'
import { deleteDatabaseConnection, TDeleteDbConnOutputs } from '../../../ts/api/apiDatabases'
import { deleteDbSysLink } from '../../../ts/api/apiSystems'
import { linkSystemsToDatabase } from '../../../ts/api/apiDatabases'
import { PageParamsStore } from '../../../ts/pageParams'
import { Database, DatabaseConnection, getDbTypeAsString, isDatabaseSupported } from '../../../ts/databases'
import CreateNewDatabaseForm from './CreateNewDatabaseForm.vue'
import { contactUsUrl, createOrgDatabaseUrl } from '../../../ts/url'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import MetadataStore from '../../../ts/metadata'
import CreateNewDbConnectionForm from './CreateNewDbConnectionForm.vue'
import { Watch } from 'vue-property-decorator'
import DatabaseConnectionReadOnlyComponent from '../../generic/DatabaseConnectionReadOnlyComponent.vue'
import SystemsTable from '../../generic/SystemsTable.vue'
import { System } from '../../../ts/systems'
import { FullDeployment } from '../../../ts/deployments'
import DeploymentEditor from '../../generic/DeploymentEditor.vue'
import DatabaseSqlEditor from '../../generic/DatabaseSqlEditor.vue'
import DashboardSqlRequestList from './DashboardSqlRequestList.vue'
import { newDeployment, TNewDeploymentOutput } from '../../../ts/api/apiDeployments'

@Component({
    components: {
        CreateNewDatabaseForm,
        GenericDeleteConfirmationForm,
        CreateNewDbConnectionForm,
        DatabaseConnectionReadOnlyComponent,
        SystemsTable,
        DeploymentEditor,
        DatabaseSqlEditor,
        DashboardSqlRequestList
    }
})
export default class FullEditDatabaseComponent extends Vue {
    currentDb: Database | null = null
    dbConn: DatabaseConnection | null = null
    relatedSystems : System[] = []
    allSystems: System[] = []

    systemsToLink : System[] = []
    deployment: FullDeployment | null = null

    showHideDelete: boolean = false
    showHideDeleteConnection: boolean = false
    showHideNewConn : boolean = false
    showHideLinkSystem: boolean = false

    addDeployment() {
        newDeployment({
            orgId: PageParamsStore.state.organization!.Id,
            systemId: null,
            dbId: this.currentDb!.Id,
        }).then((resp : TNewDeploymentOutput) => {
            this.deployment = resp.data
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

    get hasDb() : boolean {
        return !!this.dbConn
    }

    get canConnectToDb() : boolean {
        return isDatabaseSupported(this.currentDb!)
    }

    get ready() : boolean { 
        return !!this.currentDb && MetadataStore.state.dbTypesInitialized
    }

    get fullTypeString() : string {
        return `${getDbTypeAsString(this.currentDb!)} ${this.currentDb!.Version}`
    }

    refreshDbData() {
        let data = window.location.pathname.split('/')
        let dbId = Number(data[data.length - 1])

        getDatabase({
            dbId: dbId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetDatabaseOutputs) => {
            this.currentDb = resp.data.Database
            this.dbConn = resp.data.Connection
            this.allSystems = resp.data.AllSystems
            this.deployment = resp.data.Deployment

            let idSet = new Set(resp.data.RelevantSystemIds)
            this.relatedSystems = resp.data.AllSystems.filter((ele : System) => idSet.has(ele.Id))
        }).catch((err : any) => {
            window.location.replace('/404')
        })
    }

    mounted() {
        this.refreshDbData()
    }

    onDelete() {
        deleteDatabase({
            dbId: this.currentDb!.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TDeleteDatabaseOutputs) => {
            window.location.replace(createOrgDatabaseUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onEdit(db : Database) {
        this.currentDb!.Name = db.Name
        this.currentDb!.TypeId = db.TypeId
        this.currentDb!.OtherType = db.OtherType
        this.currentDb!.Version = db.Version
    }

    onNewDbConn(conn : DatabaseConnection) {
        this.dbConn = conn
        this.showHideNewConn = false
    }

    onDeleteDbConn() {
        deleteDatabaseConnection({
            connId: this.dbConn!.Id,
            dbId: this.currentDb!.Id,
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TDeleteDbConnOutputs) => {
            this.dbConn = null
            this.showHideDeleteConnection = false
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

    linkSystems() {
        linkSystemsToDatabase({
            dbId: this.currentDb!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            sysIds: this.systemsToLink.map((ele : System) => ele.Id)
        }).then(() => {
            let idSet = new Set([...this.systemsToLink, ...this.relatedSystems].map((ele: System) => ele.Id))
            this.relatedSystems = this.allSystems.filter((ele : System) => idSet.has(ele.Id))
            this.systemsToLink = []
            this.showHideLinkSystem = false
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

    onDeleteSysLink(sys : System) {
        deleteDbSysLink({
            sysId: sys.Id,
            orgId: PageParamsStore.state.organization!.Id,
            dbId: this.currentDb!.Id,
        }).then(() => {
            this.relatedSystems = this.relatedSystems.filter((ele : System) => ele.Id != sys.Id)
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
