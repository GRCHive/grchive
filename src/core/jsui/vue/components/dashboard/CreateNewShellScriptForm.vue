<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} {{ shellTypeStr }} Script
    </v-card-title>
    <v-divider></v-divider>

    <v-row class="mx-0">
        <v-col :cols="editMode ? 12 : 6">
            <v-form class="ma-4" ref="form" v-model="formValid">
                <v-text-field v-model="name"
                              label="Name"
                              filled
                              :rules="[rules.required]"
                              :readonly="!canEdit">
                </v-text-field>

                <v-textarea v-model="description"
                            label="Description"
                            filled
                            :readonly="!canEdit">
                </v-textarea> 
            </v-form>
        </v-col>

        <v-col cols="6" v-if="!editMode">
            <generic-code-editor
                v-model="script"
                :lang="shellCMLanguage"
            >
            </generic-code-editor>
        </v-col>
    </v-row>

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
            @click="edit"
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
import GenericCodeEditor from '../../generic/code/GenericCodeEditor.vue'
import { Watch } from 'vue-property-decorator'
import * as rules from '../../../ts/formRules'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { ShellTypes, ShellTypeToCodeMirror, ShellScript } from '../../../ts/shell'
import {
    newShellScript, TNewShellScriptOutput,
    editShellScript, 
}  from '../../../ts/api/apiShell'

const VueComponent = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        referenceScript: {
            type: Object,
            default: () => null as ShellScript | null
        },
        shellType: Number
    }
})

@Component({
    components: {
        GenericCodeEditor,
    }
})
export default class CreateNewShellScriptForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    name : string = ""
    description: string = ""
    script : string = ""

    get shellTypeStr() : string {
        return ShellTypes[this.shellType]
    }

    get shellCMLanguage() : string {
        return ShellTypeToCodeMirror.get(<ShellTypes>this.shellType)!
    }

    @Watch('referenceScript')
    clearForm() {
        if (!!this.referenceScript) {
            this.name = this.referenceScript.Name
            this.description = this.referenceScript.Description
        } else {
            this.name = ""
            this.description = ""
            this.script = ""
        }
    }

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }


    onSuccess(resp : TNewShellScriptOutput) {
        this.$emit('do-save', resp.data)
        if (this.editMode) {
            this.canEdit = false
        }

        if (!this.editMode) {
            this.clearForm()
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
        newShellScript({
            orgId: PageParamsStore.state.organization!.Id,
            shellType: this.shellType,
            name: this.name,
            description: this.description,
            script: this.script,
        }).then((resp: TNewShellScriptOutput) => {
            this.onSuccess(resp)
        }).catch((err : any) => {
            this.onError()
        })
    }

    doEdit() {
        editShellScript({
            orgId: PageParamsStore.state.organization!.Id,
            shellId: this.referenceScript!.Id,
            name: this.name,
            description: this.description,
        }).then((resp: TNewShellScriptOutput) => {
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

    edit() {
        this.canEdit = true
    }

    cancel() {
        this.$emit('do-cancel')

        if (this.editMode) {
            this.canEdit = false
        }
    }
}

</script>

