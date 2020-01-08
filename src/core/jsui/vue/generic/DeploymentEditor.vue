<template>
    <div>
        <v-dialog
            v-model="showUploadSoc"
            persistent
            max-width="40%"
        >
            <upload-documentation-form
                :cat-id="-1"
                @do-cancel="showUploadSoc = false"
                @do-save="onSaveUpload">
            </upload-documentation-form>
        </v-dialog>

        <v-dialog
            v-model="showAddSoc"
            persistent
            max-width="40%"
            v-if="!!editableDeployment.VendorDeployment"
        >
            <doc-searcher-form
                :exclude-files="editableDeployment.VendorDeployment.SocFiles"
                @do-select="onSelectSOC"
                @do-cancel="showAddSoc = false">
            </doc-searcher-form>
        </v-dialog>

        <v-select
            v-model="editableDeployment.DeploymentType"
            filled
            label="Type"
            :items="typeItems"
            hide-details
            :disabled="!canEdit"
            @change="changeDeploymentType"
        >
        </v-select>

        <v-divider class="my-2"></v-divider>

        <!-- self hosted -->
        <div v-if="deploymentType == 0">
            <v-form
                v-model="formValid"
            >
            </v-form>

            <v-divider></v-divider>
        </div>

        <!-- vendor hosted -->
        <div v-if="deploymentType == 1">
            <v-form
                v-model="formValid"
            >
                <v-text-field
                    v-model="editableDeployment.VendorDeployment.VendorName"
                    label="Vendor Name"
                    filled
                    :disabled="!canEdit"
                    hide-details
                    class="mb-2"
                >
                </v-text-field>

                <v-text-field
                    v-model="editableDeployment.VendorDeployment.VendorProduct"
                    label="Vendor Product"
                    filled
                    :disabled="!canEdit"
                    hide-details
                    class="mb-2"
                >
                </v-text-field>
            </v-form>

            <v-divider class="mb-2"></v-divider>

            <v-row class="title ml-2">
                SOC Reports

                <v-spacer></v-spacer>

                <v-menu offset-y>
                    <template v-slot:activator="{ on }">
                        <v-btn
                            color="primary"
                            icon
                            v-on="on"
                            :disabled="!canEdit"
                        >
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>

                    <v-list dense>
                        <v-list-item @click="showUploadSoc = true">
                            <v-list-item-title>Upload New</v-list-item-title>
                        </v-list-item>

                        <v-list-item @click="showAddSoc = true">
                            <v-list-item-title>Add Existing</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </v-row>

            <doc-file-table
                :resources="editableDeployment.VendorDeployment.SocFiles"
                :use-crud-delete="canEdit"
                confirm-delete
                @delete="unlinkSOCReport"
            >
            </doc-file-table>

            <v-divider class="mb-2"></v-divider>

            <v-row class="title ml-2">
                SOC Report Requests

                <v-spacer></v-spacer>

                <v-dialog v-model="showRequestSoc" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="warning" icon v-on="on">
                            <v-icon>mdi-plus</v-icon>
                        </v-btn>
                    </template>

                    <create-new-request-form
                        load-cats
                        :soc-request-deploy-id="editableDeployment.Id"
                        @do-cancel="showRequestSoc = false"
                        @do-save="onRequestSOC">
                    </create-new-request-form>
                </v-dialog>

            </v-row>

            <doc-request-table :resources="socRequests">
            </doc-request-table>

        </div>

        <v-list-item class="pa-0">
            <v-btn
                color="error"
                @click="cancel"
                v-if="canEdit"
            >
                Cancel
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn
                color="success"
                @click="save"
                :disabled="!canSubmit"
                v-if="canEdit"
            >
                Save
            </v-btn>

            <v-btn
                color="primary"
                @click="canEdit = true"
                v-if="!canEdit"
            >
                Edit
            </v-btn>
        </v-list-item>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { FullDeployment,    
         SelfDeployment,
         VendorDeployment,
         createEmptyVendorDeployment,
         createEmptySelfDeployment,
         deepCopyFullDeployment,
         KSelfHosted,
         KVendorHosted,
         KNoHost
} from '../../ts/deployments'
import {
    TUpdateDeploymentOutput,
    updateDeployment
} from '../../ts/api/apiDeployments'
import {
    TGetAllDocumentRequestOutput,
    getAllDocRequests
} from '../../ts/api/apiDocRequests'
import DocFileTable from './DocFileTable.vue'
import UploadDocumentationForm from '../components/dashboard/UploadDocumentationForm.vue'
import DocSearcherForm from './DocSearcherForm.vue'
import CreateNewRequestForm from '../components/dashboard/CreateNewRequestForm.vue'
import { ControlDocumentationFile } from '../../ts/controls'
import { DocumentRequest } from '../../ts/docRequests'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import DocRequestTable from './DocRequestTable.vue'

const VueProps = Vue.extend({
    props: {
        value: Object as () => FullDeployment,
    }
})

@Component({
    components: {
        DocFileTable,
        UploadDocumentationForm,
        DocSearcherForm,
        CreateNewRequestForm,
        DocRequestTable
    }
})
export default class DeploymentEditor extends VueProps {
    canEdit: boolean = false
    formValid: boolean = false
    editableDeployment: FullDeployment = {} as FullDeployment
    socRequests: DocumentRequest[] = []

    showUploadSoc: boolean = false
    showAddSoc: boolean = false

    showRequestSoc: boolean = false

    get deploymentType() : number {
        return this.editableDeployment.DeploymentType
    }

    get canSubmit() : boolean {
        if (this.deploymentType == KSelfHosted || this.deploymentType == KVendorHosted) {
            return this.formValid
        }
        return true
    }

    get typeItems() : any[] {
        return [
            {
                text: 'Self Hosted',
                value: KSelfHosted
            },
            {
                text: 'Vendor Hosted',
                value: KVendorHosted
            },
            {
                text: 'N/A',
                value: KNoHost,
            },
        ]
    }

    resetEditCopyFromProps() {
        this.editableDeployment = deepCopyFullDeployment(this.value)
    }

    mounted() {
        this.resetEditCopyFromProps()
        this.reloadSocRequests()
    }

    cancel() {
        this.canEdit = false
        this.resetEditCopyFromProps()
    }

    save() {
        updateDeployment({
            deployment: this.editableDeployment
        }).then((resp : TUpdateDeploymentOutput) => {
            this.$emit('input', resp.data)
            Vue.nextTick(this.cancel)
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

    changeDeploymentType(val : number) {
        if (this.deploymentType == KSelfHosted) {
            this.editableDeployment.VendorDeployment = null
            this.editableDeployment.SelfDeployment = Vue.observable(createEmptySelfDeployment())
        } else if (this.deploymentType == KVendorHosted) {
            this.editableDeployment.VendorDeployment = Vue.observable(createEmptyVendorDeployment())
            this.editableDeployment.SelfDeployment = null
        } else {
            this.editableDeployment.VendorDeployment = null
            this.editableDeployment.SelfDeployment = null
        }
    }

    onSaveUpload(file : ControlDocumentationFile) {
        this.showUploadSoc = false
        this.editableDeployment.VendorDeployment!.SocFiles.push(file)
    }

    unlinkSOCReport(val : ControlDocumentationFile) {
        this.editableDeployment.VendorDeployment!.SocFiles.splice(
            this.editableDeployment.VendorDeployment!.SocFiles.findIndex(
                (ele : ControlDocumentationFile) => ele.Id == val.Id
            ),
            1
        )
    }

    onSelectSOC(files : ControlDocumentationFile[]) {
        this.showAddSoc = false
        this.editableDeployment.VendorDeployment!.SocFiles.push(...files)
    }

    onRequestSOC(req : DocumentRequest) {
        this.showRequestSoc = false
        this.socRequests.push(req)
    }

    reloadSocRequests() {
        getAllDocRequests({
            orgId: PageParamsStore.state.organization!.Id,
            deployId: this.value.Id,
        }).then((resp : TGetAllDocumentRequestOutput) => {
            this.socRequests = resp.data
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
