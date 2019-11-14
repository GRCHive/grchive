<template>
    <div>
        <v-form class="ma-4" v-model="formValid" ref="form" @submit="register" onSubmit="return false;">
            <v-text-field
                v-model="email"
                label="Email"
                type="email"
                filled
                :rules="[rules.required, rules.createMaxLength(320), rules.email]"
                required
                :disabled="disableEmail"
            ></v-text-field>

            <v-text-field
                v-model="firstName"
                label="First Name"
                type="text"
                filled
                :rules="[rules.required, rules.createMaxLength(256)]"
                required
            ></v-text-field>

            <v-text-field
                v-model="lastName"
                label="Last Name"
                type="text"
                filled
                :rules="[rules.required, rules.createMaxLength(256)]"
                required
            ></v-text-field>

            <v-text-field
                v-model="password"
                label="Password"
                type="password"
                filled
                :rules="[rules.required, rules.password, passwordConfirmationMustMatch]"
                required
            ></v-text-field>

            <v-text-field
                v-model="passwordConfirm"
                label="Password Confirmation"
                type="password"
                filled
                :rules="[rules.required, passwordConfirmationMustMatch]"
                required
            ></v-text-field>

            <v-text-field
                v-model="inviteCode"
                label="Invite Code"
                type="text"
                filled
                :rules="[rules.required]"
                required
                :disabled="disableInviteCode"
            ></v-text-field>

            <v-btn
                color="success"
                class="my-2"
                :disabled="!canSubmit"
                @click="register"
            >
                Register
            </v-btn>
        </v-form>

        <div class="mx-4">
            <p class="body-1 my-2">
                Already have an account?
                <a :href="loginPageUrl">Login.</a>
            </p>

            <p class="body-1 my-2">
                Don't have an invitation code?
                <a :href="getStartedUrl">Get started.</a>
            </p>
        </div>
    </div>
</template>

<script lang="ts">

import { getStartedUrl, contactUsUrl, loginPageUrl } from '../../ts/url'
import * as rules from "../../ts/formRules"
import { postFormUrlEncoded } from "../../ts/http"
import { getCurrentCSRF } from "../../ts/csrf"
import Vue from 'vue';
import * as qs from 'query-string'
import 'url-search-params-polyfill';

interface ResponseData {
    data: {
        RedirectUrl : string
    }
}

export default Vue.extend({
    data: () => ({
        rules,
        getStartedUrl,
        contactUsUrl,
        loginPageUrl,
        formValid: false,
        firstName: "",
        lastName: "",
        email: "",
        password: "",
        passwordConfirm: "",
        inviteCode: "",
        disableEmail: false,
        disableInviteCode: false,
    }),
    computed: {
        canSubmit() : boolean {
            return this.formValid;
        }
    },
    methods: {
        register() {
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
                email: this.email,
                firstName: this.firstName,
                lastName: this.lastName,
                password: this.password,
                inviteCode: this.inviteCode,
                csrf: getCurrentCSRF(),
            }, {}).then((resp : ResponseData) => {
                window.location.assign(resp.data.RedirectUrl + '?' + qs.stringify({
                    email: this.email,
                    selfLogin: true,
                    fromRegistration: true
                }))
            }).catch((err) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Verify your invite code and try again or contact support.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            });
        },
        passwordConfirmationMustMatch() : boolean | string {
            return (this.password == this.passwordConfirm) || "Passwords must match."
        }
    },
    watch: {
        passwordConfirm() {
            // We need this to happen so that we reset the validation on the password
            // input based on what the password confirmation field is.
            //@ts-ignore
            this.$refs.form.validate()
        }
    },
    mounted() {
        let params = new URLSearchParams(window.location.search)

        if (params.has("inviteCode")) {
            this.inviteCode = params.get("inviteCode")!
            this.disableInviteCode = true
        }

        if (params.has("email")) {
            this.email = params.get("email")!
            this.disableEmail = true
        }
    }
})

</script>

