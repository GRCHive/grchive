<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Server
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-text-field v-model="os"
                      label="Operating System"
                      filled
                      :disabled="!canEdit">
        </v-text-field>

        <v-text-field v-model="location"
                      label="Location"
                      filled
                      :disabled="!canEdit">
        </v-text-field>

        <v-text-field v-model="ip"
                      label="IP Address"
                      filled
                      :disabled="!canEdit">
        </v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :disabled="!canEdit">
        </v-textarea> 
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
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { Server } from '../../../ts/infrastructure'
import { updateServer, newServer, TNewServerOutput } from '../../../ts/api/apiServers'

const VueComponent = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        referenceServer: {
            type: Object as () => Server | null,
            default: null
        }
    }
})

@Component
export default class CreateNewServerForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    name : string = ""
    os : string = ""
    location: string = ""
    ip : string = ""
    description: string = ""

    onSuccess(resp : TNewServerOutput) {
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
        newServer({
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            description: this.description,
            os: this.os,
            ip : this.ip,
            location: this.location,
        }).then((resp : TNewServerOutput) => {
            this.onSuccess(resp)
        }).catch((err : any) => {
            this.onError()
        })
    }

    doEdit() {
        updateServer({
            serverId: this.referenceServer!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            description: this.description,
            os: this.os,
            ip : this.ip,
            location: this.location,
        }).then((resp : TNewServerOutput) => {
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
        if (!!this.referenceServer) {
            this.name = this.referenceServer!.Name
            this.os = this.referenceServer!.OperatingSystem
            this.location = this.referenceServer!.Location
            this.ip = this.referenceServer!.IpAddress
            this.description = this.referenceServer!.Description
        } else {
            this.name = ""
            this.os = ""
            this.location = ""
            this.ip = ""
            this.description = ""
        }
    }
}

</script>
