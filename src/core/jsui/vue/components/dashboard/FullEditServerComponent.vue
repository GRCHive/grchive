<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-dialog persistent
                      max-width="40%"
                      v-model="showHideLinkSystems"
            >
                <v-card>
                    <v-card-title>
                        Link Systems
                    </v-card-title>

                    <system-table-with-controls
                        class="ma-4"
                        v-model="systemsToLink"
                        disable-new
                        disable-delete
                        enable-select
                        :exclude="relevantSystems"
                        :deployment-type="KSelfHosted"
                    >
                    </system-table-with-controls>

                    <v-card-actions>
                        <v-btn color="error" @click="showHideLinkSystems = false">
                            Cancel
                        </v-btn>
                        <v-spacer></v-spacer>
                        <v-btn color="success" @click="linkToSystems">
                            Link
                        </v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>

            <v-dialog persistent
                      max-width="40%"
                      v-model="showHideLinkDatabases"
            >
                <v-card>
                    <v-card-title>
                        Link Databases
                    </v-card-title>

                    <db-table-with-controls
                        class="ma-4"
                        v-model="databasesToLink"
                        disable-new
                        disable-delete
                        enable-select
                        :exclude="relevantDbs"
                        :deployment-type="KSelfHosted"
                    >
                    </db-table-with-controls>

                    <v-card-actions>
                        <v-btn color="error" @click="showHideLinkDatabases = false">
                            Cancel
                        </v-btn>
                        <v-spacer></v-spacer>
                        <v-btn color="success" @click="linkToDatabases">
                            Link
                        </v-btn>
                    </v-card-actions>
                </v-card>
            </v-dialog>

            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Server: {{ currentServer.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ currentServer.Description }}
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
                        item-name="servers"
                        :items-to-delete="[currentServer.Name]"
                        :use-global-deletion="false"
                        @do-cancel="showHideDelete = false"
                        @do-delete="onDelete">
                    </generic-delete-confirmation-form>
                </v-dialog>

            </v-list-item>
            <v-divider></v-divider>

            <v-tabs>
                <v-tab>Overview</v-tab>
                <v-tab-item>
                    <v-container fluid>
                        <v-row>
                            <v-col cols="8">
                                <create-new-server-form
                                    edit-mode
                                    :reference-server="currentServer"
                                    @do-save="onEdit"
                                >
                                </create-new-server-form>
                            </v-col>

                            <v-col cols="4">
                                <v-card>
                                    <v-card-title>
                                        Linked Deployments

                                        <v-spacer></v-spacer>

                                        <v-menu offset-y>
                                            <template v-slot:activator="{ on }">
                                                <v-btn color="primary" icon v-on="on">
                                                    <v-icon>mdi-plus</v-icon>
                                                </v-btn>
                                            </template>

                                            <v-list dense>
                                                <v-list-item @click="showHideLinkSystems = true">
                                                    <v-list-item-title>Add Systems</v-list-item-title>
                                                </v-list-item>
                                                <v-list-item @click="showHideLinkDatabases = true">
                                                    <v-list-item-title>Add Databases</v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>

                                    </v-card-title>

                                    <v-tabs
                                        v-model="relevantTab"
                                    >
                                        <v-tab>Systems</v-tab>
                                        <v-tab>Databases</v-tab>
                                    </v-tabs>

                                    <v-tabs-items
                                        v-model="relevantTab"
                                    >
                                        <v-tab-item>
                                            <systems-table
                                                :resources="relevantSystems"
                                                use-crud-delete
                                                @delete="deleteSystemLink"
                                            >
                                            </systems-table>
                                        </v-tab-item>

                                        <v-tab-item>
                                            <db-table
                                                :resources="relevantDbs"
                                                use-crud-delete
                                                @delete="deleteDbLink"
                                            >
                                            </db-table>
                                        </v-tab-item>
                                    </v-tabs-items>
                                </v-card>
                            </v-col>
                        </v-row>
                    </v-container>
                </v-tab-item>

                <v-tab>Audit Trail</v-tab>
                <v-tab-item>
                    <audit-trail-viewer
                        :resource-type="['infrastructure_servers']"
                        :resource-id="[`${currentServer.Id}`]"
                        no-header
                    >
                    </audit-trail-viewer>
                </v-tab-item>
            </v-tabs>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import CreateNewServerForm from './CreateNewServerForm.vue'
import { Server } from '../../../ts/infrastructure'
import { getServer, TGetServerOutput } from '../../../ts/api/apiServers'
import { getAllSystems, TAllSystemsOutputs } from '../../../ts/api/apiSystems'
import { allDatabase, TAllDatabaseOutputs } from '../../../ts/api/apiDatabases'
import { unlinkDeploymentFromServer, linkDeploymentToServer } from '../../../ts/api/apiDeployments'
import { deleteServer } from '../../../ts/api/apiServers'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgServersUrl } from '../../../ts/url'
import { System } from '../../../ts/systems'
import { Database, NullDatabaseFilterData } from '../../../ts/databases'
import { KSelfHosted } from '../../../ts/deployments'
import SystemsTable from '../../generic/SystemsTable.vue'
import SystemTableWithControls from '../../generic/resources/SystemTableWithControls.vue'
import DbTable from '../../generic/DbTable.vue'
import DbTableWithControls from '../../generic/resources/DbTableWithControls.vue'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'

@Component({
    components: {
        GenericDeleteConfirmationForm,
        CreateNewServerForm,
        SystemsTable,
        SystemTableWithControls,
        DbTable,
        DbTableWithControls,
        AuditTrailViewer
    }
})
export default class FullEditServerComponent extends Vue {
    KSelfHosted: number = KSelfHosted

    currentServer: Server = {} as Server
    relevantSystems : System[] = []
    relevantDbs : Database[] = []

    ready : boolean = false
    expandDescription : boolean = false
    showHideDelete: boolean = false
    relevantTab : number | null = null

    showHideLinkSystems : boolean = false
    showHideLinkDatabases : boolean = false

    systemsToLink : System[] = []
    databasesToLink : Database[] = []

    refreshData() {
        let data = window.location.pathname.split('/')
        let resourceId = Number(data[data.length - 1])

        getServer({
            serverId: resourceId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TGetServerOutput) => {
            this.currentServer = resp.data.Server
            this.relevantSystems = resp.data.RelevantSystems
            this.relevantDbs = resp.data.RelevantDbs
            this.ready = true
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
        deleteServer({
            serverId: this.currentServer.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            window.location.replace(createOrgServersUrl(PageParamsStore.state.organization!.OktaGroupName))
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

    onEdit(s : Server) {
        this.currentServer = s
    }

    linkToSystems() {
        linkDeploymentToServer({
            systemId: this.systemsToLink.map((ele : System) => ele.Id),
            serverId: this.currentServer.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.relevantSystems.push(...this.systemsToLink)
            this.systemsToLink = []
            this.showHideLinkSystems = false
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

    linkToDatabases() {
        linkDeploymentToServer({
            dbId: this.databasesToLink.map((ele : Database) => ele.Id),
            serverId: this.currentServer.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.relevantDbs.push(...this.databasesToLink)
            this.databasesToLink = []
            this.showHideLinkDatabases = false
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

    deleteSystemLink(sys : System) {
        unlinkDeploymentFromServer({
            systemId: sys.Id,
            serverId: this.currentServer.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.relevantSystems.splice(
                this.relevantSystems.findIndex(
                    (ele : System) => ele.Id == sys.Id
                ),
                1
            )
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

    deleteDbLink(db : Database) {
        unlinkDeploymentFromServer({
            dbId: db.Id,
            serverId: this.currentServer.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then(() => {
            this.relevantDbs.splice(
                this.relevantDbs.findIndex(
                    (ele : Database) => ele.Id == db.Id
                ),
                1
            )
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

