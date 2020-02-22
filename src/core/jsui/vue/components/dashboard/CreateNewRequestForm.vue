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

        <v-select
            v-model="currentLinkage"
            :items="requestLinkageItems"
            label="Link To"
            :rules="[rules.nonZero]"
            :readonly="!canEdit || !!referenceCat || !!referenceControl || catId != -1"
            filled
        >
        </v-select>

        <div v-if="currentLinkage == 1">
            <document-category-search-form-component
                v-model="realCat"
                :available-cats="availableCats"
                :load-cats="loadCats"
                :rules="[rules.required]"
                :readonly="editMode"
                v-if="catId == -1 || !!referenceCat"
            ></document-category-search-form-component>
        </div>

        <div v-if="currentLinkage == 2">
            <control-search-form-component
                v-model="linkControl"
                :rules="[rules.required]"
                :readonly="editMode || !!referenceControl"
            >
            </control-search-form-component>
        </div>
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
import { DocumentRequest, RequestLinkageMode, requestLinkageItems } from '../../../ts/docRequests'
import DocumentCategorySearchFormComponent from '../../generic/DocumentCategorySearchFormComponent.vue'
import ControlSearchFormComponent from '../../generic/ControlSearchFormComponent.vue'

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
    },
    components: {
        DocumentCategorySearchFormComponent,
        ControlSearchFormComponent
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

    currentLinkage : RequestLinkageMode = RequestLinkageMode.None
    requestLinkageItems: any[] = requestLinkageItems

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
            vendorProductId: this.vendorProductId
        }

        if (this.currentLinkage == RequestLinkageMode.DocCat) {
            params.catId = this.realCatId
        } else if (this.currentLinkage == RequestLinkageMode.Controls) {
            params.controlId = this.linkControl!.Id
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
            vendorProductId: this.vendorProductId
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
    clearForm() {
        if (!!this.referenceReq) {
            this.name = this.referenceReq.Name
            this.description = this.referenceReq.Description
            this.realCat = this.referenceCat
        } else {
            this.name = ""
            this.description = ""
            if (this.catId == -1) {
                this.realCat = null
            }
            this.linkControl = null
        }

        if (!!this.referenceCat || this.catId != -1) {
            this.currentLinkage = RequestLinkageMode.DocCat
        } else if (!!this.referenceControl) {
            this.currentLinkage = RequestLinkageMode.Controls
            this.linkControl = this.referenceControl
        }
    }
}

</script>
