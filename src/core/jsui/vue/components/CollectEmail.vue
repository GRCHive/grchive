<template>
    <v-form v-model="formValid" ref="form">
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
            :disabled="!formValid || !name || !email || !agree"
            @click="submit"
        >
            Submit
        </v-btn>

        <v-snackbar 
            v-model="showSnack"
            :color="snackIsError ? 'error' : 'success'"
            bottom
            timeout=10000
        >
            {{ snackText }} 

            <v-btn
                color="primary"
                v-if="snackShowContact"
                :href="contactUsUrl"
            >
                Contact Us
            </v-btn>

            <v-btn
                color="secondary"
                @click="showSnack = false"
            >
                Close
            </v-btn>
        </v-snackbar>
    </v-form>
</template>

<script lang="ts">

import * as rules from "../../ts/formRules"
import { contactUsUrl } from "../../ts/url"
import { postFormUrlEncoded} from "../../ts/http"
import Vue from 'vue';

export default Vue.extend({
    props: {
        'companyName' : String,
    },
    data: () => ({
        contactUsUrl,
        name: undefined,
        email: undefined,
        agree: false,
        rules: rules,
        formValid: false,
    }),
    methods: {
        submit() {
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
            }).then(() => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Success! We will reach out soon.",
                    false,
                    "Contact Us",
                    contactUsUrl,
                    false);

            }).catch((err) => {
                if (err.response.data.IsDuplicate) {
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
