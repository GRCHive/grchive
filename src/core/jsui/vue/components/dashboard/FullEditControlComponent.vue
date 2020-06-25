<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <v-dialog v-model="showNewFolder" persistent max-width="40%">
            <create-new-folder-form
                :control-id="fullControlData.Control.Id"
                @do-cancel="showNewFolder = false"
                @do-save="saveFolder"
                v-if="!!fullControlData"
            >
            </create-new-folder-form>
        </v-dialog>

        <v-dialog v-model="showEditFolder" persistent max-width="40%">
            <create-new-folder-form
                edit-mode
                dialog-mode
                :control-id="fullControlData.Control.Id"
                :reference-folder="editFolder"
                @do-cancel="showEditFolder = false"
                @do-save="onEditFolder"
                v-if="!!fullControlData && !!editFolder"
            >
            </create-new-folder-form>
        </v-dialog>

        <v-dialog v-model="showDeleteFolder" persistent max-width="40%">
            <generic-delete-confirmation-form
                item-name="folders"
                :items-to-delete="[deleteFolder.Name]"
                :use-global-deletion="false"
                @do-cancel="showDeleteFolder = false"
                @do-delete="onDeleteFolder"
                v-if="!!deleteFolder"
            >
            </generic-delete-confirmation-form>
        </v-dialog>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Control: {{ fullControlData.Control.Name }}
                        <span class="subtitle-1" v-if="fullControlData.Control.Name != fullControlData.Control.Identifier">
                            ({{ fullControlData.Control.Identifier }})
                        </span>
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ fullControlData.Control.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-spacer></v-spacer>
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
                            item-name="control"
                            :items-to-delete="[`${fullControlData.Control.Name} (${fullControlData.Control.Identifier})`]"
                            :use-global-deletion="false"
                            @do-cancel="showHideDelete = false"
                            @do-delete="onDelete">
                        </generic-delete-confirmation-form>
                    </v-dialog>
                </v-list-item-action>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-tabs>
                    <v-tab>Overview</v-tab>
                    <v-tab-item>
                        <v-row>
                            <v-col cols="5">
                                <create-new-control-form 
                                     :node-id="-1"
                                     :risk-id="-1"
                                     :edit-mode="true"
                                     :control="fullControlData.Control"
                                     :staged-edits="true"
                                     @do-save="onEditControl"
                                >
                                </create-new-control-form>
                            </v-col>

                            <v-col cols="7">
                                <v-card class="mb-4" v-if="!!relevantFolders">
                                    <v-card-title>
                                        <span class="mr-2">
                                            Documentation
                                        </span>
                                        <v-spacer></v-spacer>

                                        <v-menu bottom left offset-y>
                                            <template v-slot:activator="{on}">
                                                <v-btn
                                                    icon
                                                    color="primary"
                                                    v-on="on"
                                                >
                                                    <v-icon>mdi-plus</v-icon>
                                                </v-btn>
                                            </template>
                                            <v-list dense>
                                                <v-list-item @click="showNewFolder = true">
                                                    <v-list-item-title>
                                                        New Folder
                                                    </v-list-item-title>
                                                </v-list-item>
                                            </v-list>
                                        </v-menu>
                                    </v-card-title>
                                    <v-divider></v-divider>
                                    <v-tabs v-model="currentFolderIdx">
                                        <template v-for="folder in relevantFolders">
                                            <v-tab :key="`tab-${folder.Id}`">
                                                {{ folder.Name }}
                                                <v-spacer></v-spacer>
                                                <v-menu bottom left offset-y>
                                                    <template v-slot:activator="{on}">
                                                        <v-btn icon v-on="on" @mousedown.stop @click.stop>
                                                            <v-icon small>
                                                                mdi-dots-vertical
                                                            </v-icon>
                                                        </v-btn>
                                                   </template>
                                                   <v-list dense>
                                                        <v-list-item @click="startEditFolder(folder)">
                                                            <v-list-item-title>
                                                                Edit
                                                            </v-list-item-title>
                                                        </v-list-item>

                                                        <v-list-item @click="startDeleteFolder(folder)">
                                                            <v-list-item-title>
                                                                Delete
                                                            </v-list-item-title>
                                                        </v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-tab>
                                            <v-tab-item :key="`item-${folder.Id}`">
                                                <doc-file-manager
                                                    :folder="folder"
                                                    disable-sample
                                                >
                                                </doc-file-manager>
                                            </v-tab-item>
                                        </template>
                                    </v-tabs>
                                </v-card>

                                <v-card v-if="!!relevantRequests" class="mb-4">
                                    <v-card-title>
                                        <span class="mr-2">
                                            Requests
                                        </span>
                                        <v-spacer></v-spacer>

                                        <v-dialog v-model="showHideRequest" persistent max-width="40%">
                                            <template v-slot:activator="{ on }">
                                                <v-btn color="warning" icon v-on="on">
                                                    <v-icon>mdi-plus</v-icon>
                                                </v-btn>
                                            </template>

                                            <create-new-request-form
                                                :reference-control="fullControlData.Control"
                                                @do-cancel="showHideRequest = false"
                                                @do-save="newRequest"
                                            >
                                            </create-new-request-form>
                                        </v-dialog>
                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <doc-request-table :resources="relevantRequests">
                                    </doc-request-table>
                                </v-card>

                                <v-card class="mb-4">
                                    <v-card-title>
                                        Related Resources
                                    </v-card-title>
                                    <v-divider></v-divider>

                                    <v-tabs>
                                        <v-tab>Process Flows</v-tab>
                                        <v-tab-item>
                                            <process-flow-table
                                                :resources="fullControlData.Flows"
                                            >
                                            </process-flow-table>
                                        </v-tab-item>

                                        <v-tab>Risks</v-tab>
                                        <v-tab-item>
                                            <risk-table
                                                :resources="fullControlData.Risks"
                                            >
                                            </risk-table>
                                        </v-tab-item>

                                        <v-tab :disabled="!relevantSystems">Systems</v-tab>
                                        <v-tab-item>
                                            <systems-table
                                                :resources="relevantSystems"
                                                v-if="!!relevantSystems"
                                            >
                                            </systems-table>
                                        </v-tab-item>

                                        <v-tab :disabled="!relevantAccounts">Accounts</v-tab>
                                        <v-tab-item>
                                            <general-ledger-accounts-table
                                                :resources="relevantAccounts"
                                                v-if="!!relevantAccounts"
                                            >
                                            </general-ledger-accounts-table>
                                        </v-tab-item>
                                    </v-tabs>
                                </v-card>
                            </v-col>
                        </v-row>
                    </v-tab-item>

                    <v-tab>Audit Trail</v-tab>
                    <v-tab-item>
                        <audit-trail-viewer
                            :resource-type="['process_flow_controls']"
                            :resource-id="[`${fullControlData.Control.Id}`]"
                            no-header
                        >
                        </audit-trail-viewer>
                    </v-tab-item>
                </v-tabs>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import CreateNewControlForm from './CreateNewControlForm.vue'
import { FullControlData, ControlDocumentationFile } from '../../../ts/controls'
import { getSingleControl, TSingleControlInput, TSingleControlOutput } from '../../../ts/api/apiControls'
import { createRiskUrl, contactUsUrl, createOrgControlsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { System } from '../../../ts/systems'
import SystemsTable from '../../generic/SystemsTable.vue'
import RiskTable from '../../generic/RiskTable.vue'
import ProcessFlowTable from '../../generic/ProcessFlowTable.vue'
import DocFileManager from '../../generic/DocFileManager.vue'
import DocRequestTable from '../../generic/DocRequestTable.vue'
import {
    TAllControlSystemLinkOutput, allControlSystemLink
} from '../../../ts/api/apiControlSystemLinks'
import GeneralLedgerAccountsTable from '../../generic/GeneralLedgerAccountsTable.vue'
import { allControlGLLink, TAllControlGLLinkOutput } from '../../../ts/api/apiControlGLLinks'
import { GeneralLedger, GeneralLedgerAccount } from '../../../ts/generalLedger'
import { allControlFolderLink, TAllControlFolderLinkOutput } from '../../../ts/api/apiControlFolderLinks'
import {
    allFolderFileLink, TAllFolderFileLinkOutput ,
    newFolderFileLink,
    deleteFolderFileLink
} from '../../../ts/api/apiFolderFileLinks'
import {
    allDocRequestControlLink, TAllDocRequestControlLinksOutput
} from '../../../ts/api/apiDocRequestControlLinks'
import { deleteFolder } from '../../../ts/api/apiFolders'
import { FileFolder } from '../../../ts/folders'
import { DocumentRequest } from '../../../ts/docRequests'

import CreateNewFolderForm from './CreateNewFolderForm.vue'
import CreateNewRequestForm from './CreateNewRequestForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import AuditTrailViewer from '../../generic/AuditTrailViewer.vue'
import { deleteControls, TDeleteControlOutput} from '../../../ts/api/apiControls'

@Component({
    components: {
        CreateNewControlForm,
        SystemsTable,
        RiskTable,
        ProcessFlowTable,
        GeneralLedgerAccountsTable,
        DocFileManager,
        DocRequestTable,
        CreateNewFolderForm,
        CreateNewRequestForm,
        GenericDeleteConfirmationForm,
        AuditTrailViewer,
    }
})
export default class FullEditControlComponent extends Vue {
    expandDescription: boolean = false
    fullControlData: FullControlData | null = null
    relevantSystems: System[] | null = null
    relevantGL: GeneralLedger | null =  null
    relevantFolders: FileFolder[] | null =  null
    relevantRequests: DocumentRequest[] | null = null
    folderToFiles: Record<number, ControlDocumentationFile[]> = Object()
    currentFolderIdx: number = 0

    showNewFolder: boolean = false
    showEditFolder: boolean = false
    showDeleteFolder: boolean = false
    showHideRequest: boolean = false
    showHideDelete: boolean = false

    editFolder : FileFolder | null = null
    deleteFolder : FileFolder | null = null

    get ready() : boolean {
        return this.fullControlData != null
    }

    get relevantAccounts() : GeneralLedgerAccount[] | null {
        if (!this.relevantGL) {
            return null
        }
        return this.relevantGL.listAccounts
    }

    refreshSystemLink() {
        allControlSystemLink({
            controlId: this.fullControlData!.Control.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllControlSystemLinkOutput) => {
            this.relevantSystems = resp.data.Systems!
        })
    }

    refreshGLLink() {
        allControlGLLink({
            controlId: this.fullControlData!.Control.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllControlGLLinkOutput) => {
            this.relevantGL = new GeneralLedger()
            this.relevantGL.rebuildGL(resp.data.Categories!, resp.data.Accounts!)
        })
    }

    refreshFileFolders() {
        allControlFolderLink({
            controlId: this.fullControlData!.Control.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllControlFolderLinkOutput) => {
            this.relevantFolders = resp.data.Folders!
        })
    }

    refreshRequests() {
        allDocRequestControlLink({
            controlId: this.fullControlData!.Control.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllDocRequestControlLinksOutput) => {
            this.relevantRequests = resp.data.Requests!
        })
    }

    refreshData() {
        let data = window.location.pathname.split('/')
        let controlId = Number(data[data.length - 1])

        getSingleControl(<TSingleControlInput>{
            controlId: controlId,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TSingleControlOutput) => {
            this.fullControlData = resp.data
            this.refreshSystemLink()
            this.refreshGLLink()
            this.refreshFileFolders()
            this.refreshRequests()
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

    onEditControl(control : ProcessFlowControl) {
        Vue.set(this.fullControlData!, 'Control', control)
    }

    generateRiskUrl(riskId : number) : string {
        return createRiskUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            riskId)
    }

    mounted() {
        this.refreshData()
    }

    saveFolder(folder : FileFolder) {
        if (!this.relevantFolders) {
            return
        }
        this.relevantFolders!.push(folder)
        this.showNewFolder = false
    }

    onEditFolder(folder : FileFolder) {
        if (!this.relevantFolders) {
            return
        }

        let idx : number = this.relevantFolders.findIndex(
            (ele : FileFolder) => ele.Id == folder.Id)

        if (idx == -1) {
            return
        }

        Vue.set(this.relevantFolders, idx, folder)
        this.showEditFolder = false
    }

    onDeleteFolder() {
        if (!this.deleteFolder) {
            return
        }

        let folder : FileFolder = this.deleteFolder

        deleteFolder({
            orgId: PageParamsStore.state.organization!.Id,
            folderId: folder.Id,
        }).then(() => {
            let idx : number = this.relevantFolders!.findIndex(
                (ele : FileFolder) => ele.Id == folder.Id)

            if (idx == -1) {
                return
            }

            this.relevantFolders!.splice(idx, 1)
            this.showDeleteFolder = false
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

    startDeleteFolder(folder : FileFolder) {
        this.deleteFolder = folder
        this.showDeleteFolder = true
    }

    startEditFolder(folder : FileFolder) {
        this.editFolder = folder
        this.showEditFolder = true
    }

    newRequest(req : DocumentRequest) {
        this.relevantRequests!.unshift(req)
        this.showHideRequest = false
    }

    onDelete() {
        deleteControls({
            nodeId: -1,
            riskIds: [-1],
            controlIds: [this.fullControlData!.Control.Id],
            global: true
        }).then(() => {
            window.location.replace(createOrgControlsUrl(PageParamsStore.state.organization!.OktaGroupName))
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
}

</script>
