<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Form
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required]"
                      :readonly="!canEdit">
        </v-text-field>
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
import { PageParamsStore } from '../../../ts/pageParams'
import { newFolder, TNewFolderOutput } from '../../../ts/api/apiFolders'
import { updateFolder, TUpdateFolderOutput } from '../../../ts/api/apiFolders'
import { FileFolder } from '../../../ts/folders'

const Props = Vue.extend({
    props: {
        controlId: Number,
        editMode: {
            type: Boolean,
            default: false
        },
        dialogMode: {
            type: Boolean,
            default: false
        },
        referenceFolder: {
            type : Object,
            default: () => null as FileFolder | null
        }
    }
})

@Component
export default class CreateNewFolderForm extends Props {
    name : string = ""
    canEdit: boolean = false
    formValid : boolean = false
    rules: any = rules

    cancel() {
        if (this.editMode && !this.dialogMode) {
            this.canEdit = false
        }
        this.clearForm()
        this.$emit('do-cancel')
    }

    doSave() {
        newFolder({
            name: this.name,
            orgId: PageParamsStore.state.organization!.Id,
            controlId: this.controlId
        }).then((resp : TNewFolderOutput) => {
            this.$emit('do-save', resp.data)
            this.clearForm()
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    doEdit() {
        updateFolder({
            name: this.name,
            orgId: PageParamsStore.state.organization!.Id,
            controlId: this.controlId,
            folderId: this.referenceFolder!.Id
        }).then((resp : TUpdateFolderOutput) => {
            this.$emit('do-save', resp.data)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
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
        this.canEdit = !this.editMode || this.dialogMode
        this.clearForm()
    }

    @Watch('referenceFolder')
    clearForm() {
        if (!!this.referenceFolder) {
            this.name = this.referenceFolder.Name
        } else {
            this.name = ""
        }
    }
}

</script>
