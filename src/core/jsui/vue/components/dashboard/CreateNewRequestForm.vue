<template>

<v-card>
    <v-card-title>
        New Documentation Request 
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required]"
        ></v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
        ></v-textarea> 
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
import { newDocRequest, TNewDocRequestOutput } from '../../../ts/api/apiDocRequests'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
        catId: Number
    }
})

@Component
export default class CreateNewRequestForm extends Props {
    formValid: boolean = false
    rules = rules
    name: string = ""
    description: string = ""

    cancel() {
        this.$emit('do-cancel')
    }

    save() {
        newDocRequest({
            name: this.name,
            description: this.description,
            catId: this.catId,
            orgId: PageParamsStore.state.organization!.Id,
            requestedUserId: PageParamsStore.state.user!.Id,
        }).then((resp : TNewDocRequestOutput) => {
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
