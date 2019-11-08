<template>
    <v-form v-model="formValid" ref="form" @submit="submit" onSubmit="return false;">
        <v-text-field
            v-model="name"
            label="Name"
            filled
            :rules="[rules.required, rules.createMaxLength(320)]"
            required
        ></v-text-field>

        <v-text-field
            v-model="email"
            label="Email"
            type="email"
            filled
            :rules="[rules.required, rules.createMaxLength(320), rules.email]"
            required
        ></v-text-field>

        <v-switch
            v-model="agree"
            :label="`I agree to ${companyName} collecting and storing my personal information to send me updates about future services and products.`"
            required
            :rules="[rules.required]"
        >
        </v-switch>

        <v-btn
            color="success"
            class="my-2"
            :disabled="!canSubmit"
            @click="submit"
        >
            Submit
        </v-btn>
    </v-form>
</template>

<script lang="ts">

import * as rules from "../../ts/formRules"
import { contactUsUrl } from "../../ts/url"
import { postFormUrlEncoded } from "../../ts/http"
import Vue from 'vue';

export default Vue.extend({
    data: () => ({
        contactUsUrl,
        name: undefined,
        email: undefined,
        agree: false,
        rules: rules,
        formValid: false
    }),
    props: {
        companyName: String
    },
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.name && this.$data.email && this.$data.agree;
        }
    },
    methods: {
        submit() {
            if (!this.canSubmit) {
                return;
            }
            // ts-ignore is needed since as of Vuetify 2.0.19 we can't type assert the form
            // to a Vuetify VForm to have TypeScript properly fin the validate function.
            // @ts-ignore
            if (!this.$refs.form.validate()) {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Please refresh the page and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
                return;
            }

            postFormUrlEncoded('#', {
                name: this.$data.name,
                email: this.$data.email
            }, {}).then(() => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Success! We will reach out soon.",
                    false,
                    "Contact Us",
                    contactUsUrl,
                    false);
            }).catch((err) => {
                if (!!err.response && err.response.data.IsDuplicate) {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        "Oops! It looks like you already gave us your email address. Contact us if you require assistance.",
                        true,
                        "Contact Us",
                        contactUsUrl,
                        true);

                } else {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        "Oops! It looks like something went wrong on our end. Try again later or get in touch directly.",
                        true,
                        "Contact Us",
                        contactUsUrl,
                        true);
                }
            })
        }
    }
});

</script>
