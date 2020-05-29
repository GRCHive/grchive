<template>

<v-card :flat="hideHeader">
    <template v-if="!hideHeader">
        <v-card-title>
            {{ editMode ? "Edit" : "New" }} SSH Connection (Password)
        </v-card-title>
        <v-divider></v-divider>
    </template>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="username"
                      label="Username"
                      filled
                      :rules="[rules.required]"
                      :readonly="!canEdit"
        >
        </v-text-field>

        <v-text-field v-model="password"
                    label="Password"
                    filled
                    :type="(!canEdit && showPassword) ? 'text' : 'password'"
                    :readonly="!canEdit"
        >
            <template v-slot:append v-if="!canEdit">
                <v-btn icon v-if="!showPassword" @click="onShowPassword" :loading="loadingPassword">
                    <v-icon>mdi-eye</v-icon>
                </v-btn>

                <v-btn icon v-else @click="onHidePassword">
                    <v-icon>mdi-eye-off</v-icon>
                </v-btn>
            </template>
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
            @click="edit"
            v-if="!canEdit"
            :loading="loadingPassword"
        >
            Edit
        </v-btn>
    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    newServerSSHPasswordConnection, editServerSSHPasswordConnection,
    TNewServerSSHConnectionOutput,
    getServerSSHPasswordConnection, TGetServerConnectionOutput,
} from '../../../ts/api/apiServerConnection'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import {
    ServerSSHConnectionGeneric,
} from '../../../ts/infrastructure'
import * as rules from '../../../ts/formRules'

const Props = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        serverId: Number,
        referenceConnection: {
            type: Object,
            default: () => null as ServerSSHConnectionGeneric | null
        },
        hideHeader: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class CreateServerSshPasswordConnection extends Props {
    rules: any = rules

    canEdit: boolean = false
    formValid: boolean = false
    username: string = ""
    password: string = ""

    realDecryptedPassword : string | null = null
    showPassword : boolean = false
    loadingPassword: boolean = false

    onShowPassword() {
        if (!!this.realDecryptedPassword) {
            this.password = this.realDecryptedPassword!
            this.showPassword = true
        } else {
            this.loadingPassword = true
            getServerSSHPasswordConnection({
                orgId: PageParamsStore.state.organization!.Id,
                serverId: this.serverId,
                connectionId: this.referenceConnection!.Id,
            }).then((resp : TGetServerConnectionOutput) => {
                this.realDecryptedPassword = resp.data.Password
                this.password = this.realDecryptedPassword!
                this.showPassword = true
            }).catch(this.onError).finally(() => {
                this.loadingPassword = false
            })
        }
    }

    onHidePassword() {
        this.showPassword = false
    }
 
    onSuccess(r : TNewServerSSHConnectionOutput) {
        this.$emit('do-save', r.data)
        if (!this.editMode) {
            this.clearForm()
        } else {
            this.canEdit = false
        }

        // This needs to get reset just in case the password changed so force
        // the re-grab here just to be safe.
        this.showPassword = false
        this.realDecryptedPassword = null
    }

    onError(err : any) {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops. Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true);
    }

    doSave() {
        newServerSSHPasswordConnection({
            orgId: PageParamsStore.state.organization!.Id,
            serverId: this.serverId,
            username: this.username,
            password: this.password,
        }).then(this.onSuccess).catch(this.onError)
    }

    doEdit() {
        editServerSSHPasswordConnection({
            orgId: PageParamsStore.state.organization!.Id,
            serverId: this.serverId,
            connectionId: this.referenceConnection!.Id,
            username: this.username,
            password: this.password,
        }).then(this.onSuccess).catch(this.onError)
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
            this.clearForm()
            this.canEdit = false
        }
    }

    edit() {
        this.canEdit = true
        // Force user to create a new password.
        this.password = ""
    }

    clearForm() {
        if (!!this.referenceConnection) {
            this.username = this.referenceConnection!.Username

            // The password isn't sent to us unless the user directly requests the password to be sent and when it is
            // sent to us, it's stored in realDecryptedPassword.
            if (!!this.realDecryptedPassword) {
                this.password = this.realDecryptedPassword
            } else {
                // Some garbage to show something.
                this.password = Math.random().toString(16)
            }
        } else {
            this.username = ""
            this.password = ""
        }
    }

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }
}

</script>
