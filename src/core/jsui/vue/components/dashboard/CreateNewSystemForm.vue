<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} System
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required]"
                      :readonly="!canEdit">
        </v-text-field>

        <v-text-field v-model="purpose"
                      label="Purpose"
                      filled
                      :readonly="!canEdit">
        </v-text-field> 

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :readonly="!canEdit">
        </v-textarea> 

    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
            v-if="canEdit || dialogMode"
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
import * as rules from '../../../ts/formRules'
import { TNewSystemOutputs, newSystem} from '../../../ts/api/apiSystems'
import { TEditSystemOutputs, editSystem} from '../../../ts/api/apiSystems'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { System } from '../../../ts/systems'

const VueComponent = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        dialogMode: {
            type: Boolean,
            default: false
        },
        referenceSystem: {
            type: Object as () => System | null,
            default: null
        }
    }
})

@Component
export default class CreateNewSystemForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    name : string = ""
    purpose : string = ""
    description: string = ""

    doSave() {
        newSystem({
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            purpose: this.purpose,
            description: this.description,
        }).then((resp : TNewSystemOutputs) => {
            this.$emit('do-save', resp.data)
            this.clearForm()
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    doEdit() {
        editSystem({
            sysId: this.referenceSystem!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            purpose: this.purpose,
            description: this.description,
        }).then((resp : TNewSystemOutputs) => {
            this.$emit('do-save', resp.data)
            this.canEdit = false
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
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

    cancel() {
        this.$emit('do-cancel')
        if (this.editMode) {
            this.canEdit = false
        }
    }

    edit() {
        this.canEdit = true
    }

    mounted() {
        this.canEdit = !this.editMode
    }

    clearForm() {
        if (!!this.referenceSystem) {
            this.name = this.referenceSystem.Name
            this.purpose = this.referenceSystem.Purpose
            this.description = this.referenceSystem.Description
        } else {
            this.name = ""
            this.purpose = ""
            this.description = ""
        }
    }
}

</script>
