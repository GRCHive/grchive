<template>

<v-card>
    <v-card-title>
        New Database Connection
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="host"
                      label="Host"
                      filled
                      :rules="[rules.required]">
        </v-text-field>

        <v-text-field :value="port"
                      @input="port = Number(arguments[0])"
                      label="Port"
                      filled
                      type="number"
                      :rules="[rules.required]">
        </v-text-field>

        <v-text-field v-model="dbName"
                      label="Database Name"
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
                      type="password">
        </v-text-field>

        <span class="body-1">Extra Parameters</span>
        <v-divider></v-divider>
        <string-dict-form-component
            v-model="parameters"
        >
        </string-dict-form-component>
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
import StringDictFormComponent from '../../generic/StringDictFormComponent.vue'

const VueComponent = Vue.extend({
    props: {
        dbId: Number,
    }
})

@Component({
    components: {
        StringDictFormComponent
    }
})
export default class CreateNewDbConnectionForm extends VueComponent {
    formValid: boolean = false
    rules: any = rules

    host : string = ""
    port : number = 0
    dbName : string = ""
    parameters : Record<string, string> = Object()

    username : string = ""
    password: string = ""

    cancel() {
        this.$emit('do-cancel')
    }

    clearForm() {
        this.host = ""
        this.port = 0
        this.dbName = ""
        this.username = ""
        this.password = ""
    }

    save() {
        newDatabaseConnection({
            dbId: this.dbId,
            orgId: PageParamsStore.state.organization!.Id,
            host: this.host,
            port: this.port,
            dbName: this.dbName,
            parameters: this.parameters,
            username: this.username,
            password: this.password,
        }).then((resp : TNewDbConnOutputs) => {
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
}

</script>
