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

        <v-dialog
            v-model="showRequestSoc"
            persistent
            max-width="40%"
        >
        </v-dialog>

        <v-select
            v-model="deploymentType"
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

                            <v-list-item @click="showRequestSoc = true">
                                <v-list-item-title>Request Missing</v-list-item-title>
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
            </v-form>

            <v-divider></v-divider>
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
         deepCopyFullDeployment } from '../../ts/deployments'
import DocFileTable from './DocFileTable.vue'
import UploadDocumentationForm from '../components/dashboard/UploadDocumentationForm.vue'
import DocSearcherForm from './DocSearcherForm.vue'
import { ControlDocumentationFile } from '../../ts/controls'

const VueProps = Vue.extend({
    props: {
        value: Object as () => FullDeployment,
    }
})

const KSelfHosted : number = 0
const KVendorHosted : number = 1
const KNoHost : number = -1

@Component({
    components: {
        DocFileTable,
        UploadDocumentationForm,
        DocSearcherForm
    }
})
export default class DeploymentEditor extends VueProps {
    deploymentType: number = KNoHost
    canEdit: boolean = false
    formValid: boolean = false
    editableDeployment: FullDeployment = {} as FullDeployment

    showUploadSoc: boolean = false
    showAddSoc: boolean = false

    showRequestSoc: boolean = false

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

    resetDeploymentType() {
        if (!!this.value.SelfDeployment) {
            this.deploymentType = KSelfHosted
        } else if (!!this.value.VendorDeployment) {
            this.deploymentType = KVendorHosted
        } else {
            this.deploymentType = KNoHost
        }
    }

    mounted() {
        this.resetDeploymentType()
        this.resetEditCopyFromProps()
    }

    cancel() {
        this.canEdit = false
        this.resetDeploymentType()
        this.resetEditCopyFromProps()
    }

    save() {
        this.$emit('input', JSON.parse(JSON.stringify(this.editableDeployment)))
        Vue.nextTick(this.cancel)
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
}

</script>
