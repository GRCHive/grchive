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
                      :disabled="!canEdit"
        ></v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :disabled="!canEdit"
        ></v-textarea> 
        
        <document-category-search-form-component
            v-if="catId == -1"
            v-model="realCat"
            :available-cats="availableCats"
            :load-cats="loadCats"
            :rules="[rules.required]"
            :disabled="!canEdit"
        ></document-category-search-form-component>
    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
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
import * as rules from '../../../ts/formRules'
import { newDocRequest, TNewDocRequestOutput } from '../../../ts/api/apiDocRequests'
import { updateDocRequest, TUpdateDocRequestOutput } from '../../../ts/api/apiDocRequests'
import { ControlDocumentationCategory } from '../../../ts/controls'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import { DocumentRequest } from '../../../ts/docRequests'
import DocumentCategorySearchFormComponent from '../../generic/DocumentCategorySearchFormComponent.vue'

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
        socRequestDeployId: {
            type: Number,
            default: -1
        },
        editMode: {
            type: Boolean,
            default: false
        },
        referenceReq: {
            type: Object as () => DocumentRequest | null,
            default: null
        }
    },
    components: {
        DocumentCategorySearchFormComponent
    }
})

@Component
export default class CreateNewRequestForm extends Props {
    formValid: boolean = false
    rules = rules
    name: string = ""
    description: string = ""
    realCat: ControlDocumentationCategory | null = null
    canEdit: boolean = false

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
        }
    }

    onSuccess(resp : TNewDocRequestOutput | TUpdateDocRequestOutput) {
        this.$emit('do-save', resp.data)
        if (this.editMode) {
            this.canEdit = false
        }
    }

    onError() {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops. Something went wrong. Try again.",
            false,
            "",
            contactUsUrl,
            true);
    }

    doSave() {
        newDocRequest({
            name: this.name,
            description: this.description,
            catId: this.realCatId,
            orgId: PageParamsStore.state.organization!.Id,
            requestedUserId: PageParamsStore.state.user!.Id,
            socRequestDeployId: this.socRequestDeployId
        }).then((resp : TNewDocRequestOutput) => {
            this.onSuccess(resp)
        }).catch((err : any) => {
            this.doSave()
        })
    }

    doEdit() {
        updateDocRequest({
            requestId: this.referenceReq!.Id,
            name: this.name,
            description: this.description,
            catId: this.realCatId!,
            orgId: PageParamsStore.state.organization!.Id,
            requestedUserId: PageParamsStore.state.user!.Id,
            socRequestDeployId: this.socRequestDeployId
        }).then((resp : TNewDocRequestOutput) => {
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
    }

    clearForm() {
        if (!!this.referenceReq) {
            this.name = this.referenceReq.Name
            this.description = this.referenceReq.Description
        } else {
            this.name = ""
            this.description = ""
        }
    }
}

</script>
