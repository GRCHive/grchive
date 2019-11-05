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
import { contactUsUrl } from "../../../ts/url"
import { getCurrentCSRF } from '../../../ts/csrf'
import { TNewProcessFlowInput, TNewProcessFlowOutput, newProcessFlow }  from '../../../ts/api/apiProcessFlow'

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
            if (!this.canSubmit) {
                return;
            }

            // Post request: name, description, csrf, organization.
            newProcessFlow(<TNewProcessFlowInput>{
                name: this.name || "",
                description: this.description || "",
                //@ts-ignore
                organization: this.$root.orgGroupId,
                csrf: getCurrentCSRF()
            }).then((resp : TNewProcessFlowOutput ) => {
                this.name = undefined;
                this.description = undefined;
                this.$emit('do-save', resp.data.Name, resp.data.Id)
            }).catch((err : any) => {
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
