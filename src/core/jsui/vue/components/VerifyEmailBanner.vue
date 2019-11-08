<template>
    <v-alert type="warning"
             v-model="showBanner"
             dismissible
             class="pa-2 ma-0"
             width="100%"
             tile
    >
        <template v-slot:prepend>
            <v-icon>mdi-exclamation</v-icon>
        </template>

        <v-row align="center">
            <v-col class="pl-4 py-0 grow">
                Please verify your email address to access all available functionality.
            </v-col>

            <v-col class="py-0 pr-2 shrink">
                <v-btn small outlined @click="sendVerification" v-if="!sentEmail">
                    Resend Verification Email
                </v-btn>

                <v-btn small outlined disabled v-else>
                    Email Sent!
                </v-btn>
            </v-col>
        </v-row>
    </v-alert>
</template>

<script lang="ts">

import Vue from 'vue'
import { PageParamsStore } from '../../ts/pageParams'
import { TRequestVerificationEmailInput, TRequestVerificationEmailOutput, requestResendVerificationEmail } from '../../ts/api/apiUsers'
import { contactUsUrl } from '../../ts/url'

export default Vue.extend({
    data: () => ({
        showBanner: Boolean(!PageParamsStore.state.user!.Verified),
        sentEmail: false
    }),
    methods: {
        sendVerification() {
            this.sentEmail = true

            requestResendVerificationEmail(<TRequestVerificationEmailInput>{
                userId: PageParamsStore.state.user!.Id
            }).then((resp : TRequestVerificationEmailOutput) => {
                setTimeout(() => {
                    this.showBanner = false
                }, 5000)
            }).catch ((err : any) => {
                this.sentEmail = false
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! It looks like something went wrong on our end. Try again later or get in touch directly.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }
    },
    watch: {
        showBanner() {
            this.$emit('toggle-banner', this.showBanner)
        }
    },
    mounted() {
        this.$emit('toggle-banner', this.showBanner)
    }
})

</script>
