<template>
    <v-form class="ma-4" v-model="formValid" ref="form" @submit="submit" onSubmit="return false;">
        <v-text-field
            v-model="email"
            label="Email"
            type="email"
            filled
            :rules="[rules.required, rules.createMaxLength(320), rules.email]"
            required
        ></v-text-field>

        <v-btn
            color="success"
            class="my-2"
            :disabled="!canSubmit"
            @click="submit"
        >
            Next
        </v-btn>

        <p class="body-1 my-2">
            Don't have an account?
            <a :href="getStartedUrl">Get started.</a>
        </p>
    </v-form>
</template>

<script lang="ts">

import { getStartedUrl, contactUsUrl } from '../../ts/url'
import * as rules from "../../ts/formRules"
import { postFormUrlEncoded } from "../../ts/http"
import { getCurrentCSRF } from "../../ts/csrf"
import Vue from 'vue';

interface ResponseData {
    data: {
        LoginUrl : string
    }
}

export default Vue.extend({
    data: () => ({
        getStartedUrl,
        rules,
        email: undefined,
        formValid: false,
    }),
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.email;
        }
    },
    methods: {
        submit() {
            // Need this since hitting enter on the form can trigger this...
            if (!this.canSubmit) {
                return;
            }

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

            postFormUrlEncoded<ResponseData>('#', {
                email: this.$data.email,
                csrf: getCurrentCSRF(),
            }).then((resp : ResponseData) => {
                window.location.assign(resp.data.LoginUrl);
            }).catch((err) => {
                if (!!err.response && err.response.data.CanNotFindIdP) {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        "Oops! Have your organization contact us to get started.",
                        true,
                        "Get Started",
                        getStartedUrl,
                        true);
                } else {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        "Oops! It looks like something went wrong on our end. Try again later or contact support.",
                        true,
                        "Contact Us",
                        contactUsUrl,
                        true);
                }
            });
        }
    }
})

</script>
