<template>

<v-card>
    <v-card-title>
        Edit Risk
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
import { newRisk } from "../../../ts/api/apiRisks"

export default Vue.extend({
    props : {
        nodeId: Number
    },
    data: () => ({
        name: "",
        description: "",
        rules,
        formValid: false
    }),
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.name.length > 0;
        }
    },
    methods: {
        clearForm() {
            this.name = ""
            this.description = ""
        },
        cancel() {
            this.$emit('do-cancel')
            this.clearForm()
        },
        save() {
            //@ts-ignore
            if (!this.canSubmit) {
                return;
            }

            newRisk(<TNewRiskInput>{
                //@ts-ignore
                csrf : this.$root.csrf,
                name : this.name,
                description: this.description,
                nodeId: this.nodeId
            }).then((resp : TNewRiskOutput) => {
                this.clearForm()
                this.$emit('do-save', resp.data)
            }).catch((err) => {
                if (!!err.response && err.response.data.IsDuplicate) {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        "A risk with this name exists already. Pick another name.",
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

