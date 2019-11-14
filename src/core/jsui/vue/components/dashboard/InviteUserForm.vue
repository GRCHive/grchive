<template>

<v-card>
    <v-card-title class="pl-3">
        Invite Users
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-combobox chips
                    clearable
                    deletable-chips
                    filled
                    label="Emails"
                    :delimiters="[,]"
                    multiple
                    :hide-no-data="true"
                    :rules="[rules.required,
                             rules.createPerElement(rules.email),
                             rules.createPerElement(rules.createMaxLength(320))]"
                    v-model="emails">
        </v-combobox>
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
            Send
        </v-btn>
    </v-card-actions>
</v-card>
    

</template>

<script lang="ts">

import Vue from 'vue'
import * as rules from '../../../ts/formRules'
import { contactUsUrl } from '../../../ts/url'
import { TInviteUsersInput, TInviteUsersOutput, inviteUsers } from '../../../ts/api/apiUsers'
import { PageParamsStore } from '../../../ts/pageParams'

export default Vue.extend({
    data: () => ({
        emails: [],
        formValid: false,
        rules
    }),
    computed: {
        canSubmit() : boolean {
            return this.formValid && this.emails.length > 0
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

            inviteUsers(<TInviteUsersInput>{
                fromUserId: PageParamsStore.state.user!.Id,
                fromOrgId: PageParamsStore.state.organization!.Id,
                toEmails: this.emails
            }).then((resp : TInviteUsersOutput) => {
                this.$emit('do-save')
            }).catch((err : any) => {
                if (!!err.response && err.response.data.FailedEmail) {
                    // @ts-ignore
                    this.$root.$refs.snackbar.showSnackBar(
                        `Oops! Some invites failed to send. Please check the email ${err.response.data.FailedEmail} and try again. There may already be a pending invitation to this user from your organization or they may already have access to your organization.`,
                        true,
                        "Contact Us",
                        contactUsUrl,
                        true);

                    let failIdx = this.emails.findIndex((ele : string) => (ele == err.response.data.FailedEmail))
                    if (failIdx != -1) {
                        this.emails.splice(0, failIdx)
                    }
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
        },
    }
})

</script>
