<template>

<v-card>
    <v-card-title>
        New Process Flow
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name" label="Name" filled :rules="[rules.required, rules.createMaxLength(256)]">
        </v-text-field>

        <v-textarea v-model="description" label="Description" filled>
        </v-textarea> 

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
            :disabled="!canSubmit"
        >
            Save
        </v-btn>
    </v-card-actions>
</v-card>
    
</template>

<script lang="ts">

import Vue from 'vue'
import * as rules from "../../../ts/formRules"
import { postFormUrlEncoded } from "../../../ts/http"
import { contactUsUrl, newProcessFlowAPIUrl } from "../../../ts/url"

interface ResponseData {
    data: {
        Name: string,
        Id: number
    }
}

export default Vue.extend({
    data: () => ({
        name: undefined,
        description: undefined,
        rules,
        formValid: false
    }),
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.name;
        }
    },
    methods: {
        cancel() {
            this.$emit('do-cancel')
        },
        save() {
            //@ts-ignore
            if (!this.canSubmit) {
                return;
            }

            // Post request: name, description, csrf, organization.
            //@ts-ignore
            postFormUrlEncoded<ResponseData>(newProcessFlowAPIUrl, {
                name: this.name,
                description: this.description || "",
                //@ts-ignore
                organization: this.$root.orgGroupId,
                //@ts-ignore
                csrf: this.$root.csrf
            }).then((resp : ResponseData) => {
                this.name = undefined;
                this.description = undefined;
                this.$emit('do-save', resp.data.Name, resp.data.Id)
            }).catch((err) => {
                if (!!err.response && err.response.data.IsDuplicate) {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        "A process flow with this name exists already. Pick another name.",
                        false,
                        "",
                        contactUsUrl,
                        true);
                } else {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        "Oops! Something went wrong. Try again.",
                        true,
                        "Contact Us",
                        contactUsUrl,
                        true);
                }
            })

        }
    }
})

</script>
