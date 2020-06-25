<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Documentation Request
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required]"
                      :readonly="!canEdit"
        ></v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :readonly="!canEdit"
        ></v-textarea> 

        <user-search-form-component
            class="mt-4"
            label="Assignee"
            :user.sync="assignee"
            :readonly="!canEdit"
        >
        </user-search-form-component>

        <date-time-picker-form-component
            v-model="dueDate"
            label="Due Date"
            :readonly="!canEdit"
            clearable
        >
        </date-time-picker-form-component>

        <template v-if="vendorProductId == -1">
            <document-category-search-form-component
                v-model="realCat"
                :available-cats="availableCats"
                :load-cats="loadCats"
                :rules="[rules.required]"
                :readonly="editMode || catId != -1 || !!referenceCat"
            ></document-category-search-form-component>

            <control-search-form-component
                v-model="linkControl"
                :readonly="editMode || !!referenceControl"
                v-if="!editMode || !!linkControl"
            >
            </control-search-form-component>

            <control-folder-search-form-component
                v-if="!!linkControl"
                :control-id="linkControl.Id"
                v-model="currentFolder"
                :readonly="editMode || !!referenceFolder"
                :rules="[rules.required]"
            >
            </control-folder-search-form-component>
        </template>
    </v-form>

    <v-card-actions>
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
            :disabled="!formValid"
            v-if="canEdit"
        >
            Save
        </v-btn>

        <v-btn
            color="success"
            @click="canEdit = true"
            v-else
        >
            Edit
        </v-btn>

    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { Watch } from 'vue-property-decorator'
import * as rules from '../../../ts/formRules'
import { newDocRequest, TNewDocRequestOutput, TNewDocRequestInput } from '../../../ts/api/apiDocRequests'
import { updateDocRequest, TUpdateDocRequestOutput, TUpdateDocRequestInput } from '../../../ts/api/apiDocRequests'
import { ControlDocumentationCategory } from '../../../ts/controls'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { DocumentRequest } from '../../../ts/docRequests'
import { FileFolder } from '../../../ts/folders'
import MetadataStore from '../../../ts/metadata'
import DocumentCategorySearchFormComponent from '../../generic/DocumentCategorySearchFormComponent.vue'
import ControlSearchFormComponent from '../../generic/ControlSearchFormComponent.vue'
import UserSearchFormComponent from '../../generic/UserSearchFormComponent.vue'
import DateTimePickerFormComponent from '../../generic/DateTimePickerFormComponent.vue'
import ControlFolderSearchFormComponent from '../../generic/ControlFolderSearchFormComponent.vue'

const Props = Vue.extend({
    props: {
        catId: {
            type: Number,
            default: -1
        },
        availableCats: {
            type: Array,
            default: () => []
        },
        loadCats: {
            type: Boolean,
            default: false
        },
        vendorProductId: {
            type: Number,
            default: -1
        },
        editMode: {
            type: Boolean,
            default: false
        },
        referenceReq: {
            type: Object,
            default: () => null as DocumentRequest | null
        },
        referenceCat: {
            type: Object,
            default: () => null as ControlDocumentationCategory | null
        },
        referenceControl: {
            type: Object,
            default: () => null as ProcessFlowControl | null
        },
        referenceFolder: {
            type: Object,
            default: () => null as FileFolder | null
        },
    },
    components: {
        DocumentCategorySearchFormComponent,
        ControlSearchFormComponent,
        UserSearchFormComponent,
        DateTimePickerFormComponent,
        ControlFolderSearchFormComponent,
    }
})

@Component
export default class CreateNewRequestForm extends Props {
    formValid: boolean = false
    rules = rules
    name: string = ""
    description: string = ""
    realCat: ControlDocumentationCategory | null = null
    linkControl : ProcessFlowControl | null = null
    canEdit: boolean = false
    assignee : User | null = null
    dueDate : Date | null = null
    currentFolder: FileFolder | null = null

    get realCatId() : number {
        if (this.catId == -1) {
            return this.realCat!.Id
        }
        return this.catId
    }

    cancel() {
        this.$emit('do-cancel')
        if (this.editMode) {
            this.canEdit = false
            this.clearForm()
        }
    }

    onSuccess(resp : TNewDocRequestOutput | TUpdateDocRequestOutput) {
        this.$emit('do-save', resp.data.Request)
        if (this.editMode) {
            this.canEdit = false
        }
    }

    onError() {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops. Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true);
    }

    doSave() {
        let params : TNewDocRequestInput = {
            name: this.name,
            description: this.description,
            orgId: PageParamsStore.state.organization!.Id,
            requestedUserId: PageParamsStore.state.user!.Id,
            vendorProductId: this.vendorProductId,
            assigneeUserId: !!this.assignee ? this.assignee.Id : null,
            dueDate: this.dueDate,
            catId: this.realCatId,
        }

        if (!!this.linkControl) {
            params.controlId = this.linkControl.Id
        }

        if (!!this.currentFolder) {
            params.folderId = this.currentFolder.Id
        }

        newDocRequest(
            params
        ).then((resp : TNewDocRequestOutput) => {
            this.onSuccess(resp)
            this.clearForm()
        }).catch((err : any) => {
            this.onError()
        })
    }

    doEdit() {
        let params : TUpdateDocRequestInput = {
            requestId: this.referenceReq!.Id,
            name: this.name,
            description: this.description,
            orgId: PageParamsStore.state.organization!.Id,
            requestedUserId: PageParamsStore.state.user!.Id,
            vendorProductId: this.vendorProductId,
            assigneeUserId: !!this.assignee ? this.assignee.Id : null,
            dueDate: this.dueDate,
            catId: this.realCatId,
        }

        // Don't let people update cat or controls. NO need 
        // to put these parameters here.

        updateDocRequest(
            params
        ).then((resp : TNewDocRequestOutput) => {
            this.onSuccess(resp)
        }).catch((err : any) => {
            this.onError()
        })
    }

    save() {
        if (this.editMode) {
            this.doEdit()
        } else {
            this.doSave()
        }
    }

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }

    @Watch('referenceControl')
    @Watch('referenceCat')
    @Watch('referenceFolder')
    clearForm() {
        if (!!this.referenceReq) {
            this.name = this.referenceReq.Name
            this.description = this.referenceReq.Description
            this.assignee = MetadataStore.getters.getUser(this.referenceReq.AssigneeUserId)
            this.dueDate = this.referenceReq.DueDate
        } else {
            this.name = ""
            this.description = ""
            if (this.catId == -1) {
                this.realCat = null
            }
            this.linkControl = null
            this.currentFolder = null
            this.assignee = null
            this.dueDate = null
        }

        if (!!this.referenceCat) {
            this.realCat = this.referenceCat
        }

        if (!!this.referenceControl) {
            this.linkControl = this.referenceControl
        }

        if (!!this.referenceFolder) {
            this.currentFolder = this.referenceFolder
        }
    }
}

</script>
