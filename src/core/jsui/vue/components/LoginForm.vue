<template>
    <div>
        <v-form class="ma-4" v-model="formValid" ref="form" @submit="checkIdp" onSubmit="return false;" v-if="!selfSignInMode">
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
                @click="checkIdp"
            >
                Next
            </v-btn>

            <p class="body-1 my-2">
                Don't have an account?
                <a :href="getStartedUrl">Get started.</a>
            </p>
        </v-form>

        <div id="oktaLogin" v-else></div>
    </div>
</template>

<script lang="ts">

import { getStartedUrl, contactUsUrl } from '../../ts/url'
import * as rules from "../../ts/formRules"
import { postFormUrlEncoded } from "../../ts/http"
import { getCurrentCSRF } from "../../ts/csrf"
import { PageParamsStore } from "../../ts/pageParams"
import Vue from 'vue';

import OktaSignIn from '@okta/okta-signin-widget';
import '@okta/okta-signin-widget/dist/css/okta-sign-in.min.css';

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
        selfSignInMode: false
    }),
    computed: {
        canSubmit() : boolean {
            return this.$data.formValid && this.$data.email;
        }
    },
    methods: {
        switchToOktaLogin() {
            this.selfSignInMode = true

            Vue.nextTick(() => {
                let signInForm = new OktaSignIn({
                    baseUrl: PageParamsStore.state.auth!.OktaServer,
                    features: {
                        registration: true,
                    },
                    username: this.email,
                    clientId: PageParamsStore.state.auth!.OktaClientId,
                    redirectUri: PageParamsStore.state.auth!.OktaRedirectUri,
                    authParams: {
                        issuer: `${PageParamsStore.state.auth!.OktaServer}/oauth2/default`,
                        display: 'page',
                        responseMode: 'query',
                        responseType: 'code',
                        scopes: PageParamsStore.state.auth!.OktaScope.split(' '),
                        getAccessToken: true,
                        getIdToken: true,
                        state: getCurrentCSRF(),
                    }
                })

                signInForm.renderEl({
                    el: '#oktaLogin',
                })
            })
        },
        checkIdp() {
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
            }, {}).then((resp : ResponseData) => {
                window.location.assign(resp.data.LoginUrl);
            }).catch((err) => {
                if (!!err.response && err.response.data.CanNotFindIdP) {
                    this.switchToOktaLogin()
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
