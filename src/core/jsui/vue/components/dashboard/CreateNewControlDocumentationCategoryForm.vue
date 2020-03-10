<template>

<v-card>
    <v-card-title class="pl-3">
        {{ editMode ? "Edit" : "New" }} Documentation Category
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                      :readonly="!canEdit"
        ></v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :readonly="!canEdit"
        ></v-textarea> 

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
            v-if="!canEdit"
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
import * as rules from "../../../ts/formRules"
import { contactUsUrl } from "../../../ts/url"
import { newControlDocCat, TNewControlDocCatInput, TNewControlDocCatOutput } from '../../../ts/api/apiControlDocumentation'
import { editControlDocCat, TEditControlDocCatInput, TEditControlDocCatOutput } from '../../../ts/api/apiControlDocumentation'
import { ControlDocumentationCategory } from '../../../ts/controls'
import { PageParamsStore } from '../../../ts/pageParams'

const FormProps = Vue.extend({
    props : {
        editMode: {
            type: Boolean,
            default: false
        },
        defaultName: {
            type: String,
            default: ""
        },
        defaultDescription: {
            type: String,
            default: ""
        },
        catId: {
            type: Number,
            default: -1
        },
    },
})

@Component
export default class CreateNewControlDocumentationCategoryForm extends FormProps {
    name: string = ""
    description: string = ""
    rules = rules
    formValid: boolean = false
    canEdit: boolean = false

    get canSubmit() : boolean {
        return this.$data.formValid && this.$data.name.length > 0;
    }

    @Watch('defaultName')
    @Watch('defaultDescription')
    clearForm() {
        this.name = this.defaultName
        this.description = this.defaultDescription
    }

    cancel() {
        this.$emit('do-cancel')
        this.clearForm()

        if (this.editMode) {
            this.canEdit = false
        }
    }

    save() {
        //@ts-ignore
        if (!this.canSubmit) {
            return;
        }

        if (this.editMode) {
            this.doEdit()
            this.canEdit = false
        } else {
            this.doSave()
        }
    }

    onSuccess(cat : ControlDocumentationCategory) {
        this.$emit('do-save', cat)
        if (!this.editMode) {
            this.clearForm()
        }
    }

    onError( err: any) {
        if (!!err.response && err.response.data.IsDuplicate) {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "A documentation category with this name exists already. Pick another name.",
                false,
                "",
                contactUsUrl,
                true);
        } else {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }
    }

    doSave() {
        newControlDocCat(<TNewControlDocCatInput>{
            name: this.name,
            description: this.description,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TNewControlDocCatOutput) => {
            this.onSuccess(resp.data)
        }).catch((err : any) => {
            this.onError(err)
        })
    }

    doEdit() {
        editControlDocCat(<TEditControlDocCatInput>{
            catId: this.catId,
            name: this.name,
            description: this.description,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TEditControlDocCatOutput) => {
            this.onSuccess(resp.data)
        }).catch((err : any) => {
            this.onError(err)
        })
    }

    mounted() {
        this.canEdit = !this.editMode
    }
}

</script>
