<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Data Object
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                      :readonly="!canEdit">
        </v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :readonly="!canEdit"
                    hide-details
        >
        </v-textarea> 

        <data-source-form-component
            v-model="source"
            :rules="[rules.required]"
        >
        </data-source-form-component>
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
import * as rules from '../../../ts/formRules'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { DataSourceLink } from '../../../ts/clientData'
import { TNewClientDataOutput, newClientData } from '../../../ts/api/apiClientData'
import DataSourceFormComponent from '../../generic/DataSourceFormComponent.vue'

const Props = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
    }
})

@Component({
    components: {
        DataSourceFormComponent
    }
})
export default class CreateNewClientDataForm extends Props {
    rules:  any = rules
    formValid : boolean = false 
    canEdit : boolean = false

    name : string = ""
    description : string = ""
    source : DataSourceLink | null = null

    cancel() {
        this.clearForm()
        this.$emit('do-cancel')
    }

    doSave() {
        newClientData({
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            description: this.description,
            sourceId: this.source!.SourceId,
            sourceTarget: this.source!.SourceTarget,
        }).then((resp : TNewClientDataOutput) => {
            this.$emit('do-save', resp.data)
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

    doEdit() {
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

    clearForm() {
        this.name = ""
        this.description = ""
    }
}

</script>
