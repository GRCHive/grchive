<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Users
                </v-list-item-title>
            </v-list-item-content>
            <v-list-item-action>
                <v-text-field outlined
                              v-model="filterText"
                              prepend-inner-icon="mdi-magnify"
                              hide-details
                ></v-text-field>
            </v-list-item-action>
            <v-spacer></v-spacer>

            <v-list-item-action>
                <v-dialog v-model="showHideInvite" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            Invite
                        </v-btn>
                    </template>
                    <invite-user-form
                      @do-cancel="cancelInvite"
                      @do-save="onInvite"
                    ></invite-user-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <user-table
            :resources="users"
            :search="filterText"
        ></user-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import MetadataStore from '../../../ts/metadata'
import InviteUserForm from './InviteUserForm.vue'
import { replaceWithMark, sanitizeTextForHTML } from '../../../ts/text'
import UserTable from '../../generic/UserTable.vue'

export default Vue.extend({
    data: () => ({
        filterText: "",
        showHideInvite: false,
    }),
    components: {
        InviteUserForm,
        UserTable
    },
    computed: {
        filter() : (a : User) => boolean {
            const filterText = this.filterText.trim()
            return (ele : User) : boolean => {
                return ele.FirstName.toLocaleLowerCase().includes(filterText.toLocaleLowerCase()) ||
                    ele.LastName.toLocaleLowerCase().includes(filterText.toLocaleLowerCase()) || 
                    ele.Email.toLocaleLowerCase().includes(filterText.toLocaleLowerCase())
            }
        },
        users() : User[] {
            return MetadataStore.state.availableUsers
        },
    },
    methods: {
        highlightText(input : string) : string {
            const safeInput = sanitizeTextForHTML(input)
            const useFilter = this.filterText.trim()
            if (useFilter.length == 0) {
                return safeInput
            }
            return replaceWithMark(
                safeInput,
                sanitizeTextForHTML(useFilter))
        },
        onInvite() {
            this.showHideInvite = false
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Successfully sent invite(s)!",
                false,
                "",
                "",
                false);
        },
        cancelInvite() {
            this.showHideInvite = false
        }
    }
})

</script>
