<template>
    <section id="content" >
        <p class="display-1">My Profile</p>
        <v-form onSubmit="return false;" v-model="formValid">
            <v-text-field v-model="formData.firstName"
                          label="First Name"
                          :disabled="!canEdit"
                          filled
                          :rules="[rules.required, rules.createMaxLength(320)]"
            >
            </v-text-field>

            <v-text-field v-model="formData.lastName"
                          label="Last Name"
                          :disabled="!canEdit"
                          filled
                          :rules="[rules.required, rules.createMaxLength(320)]"
            >
            </v-text-field>

            <v-text-field v-model="formData.email"
                          label="Email"
                          disabled
                          filled
                          :rules="[rules.required, rules.createMaxLength(320), rules.email]"
            >
            </v-text-field>

            <v-btn v-if="canEdit" color="error" @click="cancelEdit">
                Cancel
            </v-btn>

            <v-btn v-if="canEdit" color="success" @click="save">
                Save
            </v-btn>

            <v-btn v-else color="warning" @click="startEdit">
                Edit
            </v-btn>
        </v-form>

        <v-divider class="my-2"></v-divider>

        <span>
            Email Verification:

            <template v-if="isEmailVerified">
                <v-icon color="success">mdi-check-circle</v-icon>
            </template>

            <template v-else>
                <v-icon color="error">mdi-close-circle</v-icon>
                <v-btn color="warning" :disabled="!canSendVerification" @click="sendVerification">
                    Resend Verification Email
                </v-btn>
            </template>
        </span>

    </section>
</template>

<script lang="ts">

import * as rules from "../../../ts/formRules"
import { contactUsUrl } from "../../../ts/url"
import { TEditUserProfileInput, TEditUserProfileOutput, editUserProfile } from '../../../ts/api/apiUsers'
import Vue from 'vue'
import { PageParamsStore } from '../../../ts/pageParams'
import { TRequestVerificationEmailInput, TRequestVerificationEmailOutput, requestResendVerificationEmail } from '../../../ts/api/apiUsers'

export default Vue.extend({
    data: function() {
        return {
            rules,
            formValid: true,
            formData: {
                firstName: PageParamsStore.state.user!.FirstName,
                lastName: PageParamsStore.state.user!.LastName,
                email: PageParamsStore.state.user!.Email
            },
            canEdit: false,
            savedState: { firstName: "", lastName: "", email: ""},
            canSendVerification: true
        }
    },
    computed: {
        canSubmit() : boolean {
            return this.canEdit && this.formValid && this.formData.firstName.length > 0
        },
        isEmailVerified() : boolean {
            return PageParamsStore.state.user!.Verified
        },
    },
    methods: {
        startEdit: function() {
            this.savedState = JSON.parse(JSON.stringify(this.formData))
            this.canEdit = true
        },
        cancelEdit: function() {
            this.formData = this.savedState
            this.canEdit = false
        },
        save() {
            if (!this.canSubmit) {
                return;
            }

            this.canEdit = false
            //@ts-ignore
            editUserProfile(PageParamsStore.state.user!.Id, <TEditUserProfileInput>{
                firstName: this.formData.firstName,
                lastName: this.formData.lastName
            }).then((resp : TEditUserProfileOutput) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Success!.",
                    false,
                    "",
                    "",
                    false);
                window.location.reload(false);
            }).catch((err : any) => {
                this.formData = this.savedState
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        sendVerification() {
            this.canSendVerification = false
            requestResendVerificationEmail(<TRequestVerificationEmailInput>{
                userId: PageParamsStore.state.user!.Id
            }).then((resp : TRequestVerificationEmailOutput) => {
            }).catch ((err : any) => {
                this.canSendVerification = true
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! It looks like something went wrong on our end. Try again later or get in touch directly.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }
    }
})

</script>

<style scoped>

#content {
    width: 50%
}

</style>
