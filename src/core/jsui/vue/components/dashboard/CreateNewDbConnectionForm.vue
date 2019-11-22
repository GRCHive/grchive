<template>

<v-card>
    <v-card-title>
        New Database Connection
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="connString"
                      label="Connection String"
                      filled
                      :rules="[rules.required]">
        </v-text-field>

        <v-text-field v-model="username"
                      label="Username"
                      filled
                      :rules="[rules.required]">
        </v-text-field>

        <v-text-field v-model="password"
                      label="Password"
                      filled
                      :rules="[rules.required]"
                      type="password">
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
import * as rules from '../../../ts/formRules'
import { TNewDbConnOutputs, newDatabaseConnection } from '../../../ts/api/apiDatabases'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const VueComponent = Vue.extend({
    props: {
        dbId: Number
    }
})

@Component
export default class CreateNewDbConnectionForm extends VueComponent {
    formValid: boolean = false
    rules: any = rules

    connString : string = ""
    username : string = ""
    password: string = ""

    cancel() {
        this.$emit('do-cancel')
    }

    save() {
        newDatabaseConnection({
            dbId: this.dbId,
            orgId: PageParamsStore.state.organization!.Id,
            connectionString: this.connString,
            username: this.username,
            password: this.password,
        }).then((resp : TNewDbConnOutputs) => {
            this.$emit('do-save', resp.data)
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
}

</script>
