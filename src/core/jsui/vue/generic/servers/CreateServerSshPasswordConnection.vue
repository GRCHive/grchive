<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} SSH Connection (Password)
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="username"
                      label="Username"
                      filled
                      :rules="[rules.required]"
        >
        </v-text-field>

        <v-text-field v-model="password"
                    label="Password"
                    filled
                    type="password"
        >
        </v-text-field>
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
        >
            Save
        </v-btn>
    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    newServerSSHPasswordConnection,
    TNewServerSSHConnectionOutput
} from '../../../ts/api/apiServerConnection'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import * as rules from '../../../ts/formRules'

const Props = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        serverId: Number
    }
})

@Component
export default class CreateServerSshPasswordConnection extends Props {
    rules: any = rules
    formValid: boolean = false
    username: string = ""
    password: string = ""
 
    onSuccess(r : TNewServerSSHConnectionOutput) {
        this.$emit('do-save', r.data)
        if (!this.editMode) {
            this.clearForm()
        }
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
        }
    }

    clearForm() {
        this.username = ""
        this.password = ""
    }

    mounted() {
        this.clearForm()
    }
}

</script>
