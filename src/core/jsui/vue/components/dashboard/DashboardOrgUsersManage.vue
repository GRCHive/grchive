<template>
    <div class="ma-4" v-if="ready">
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
                      :preload-roles="roleList"
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
            show-role
            :available-roles="roles"
        ></user-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import MetadataStore from '../../../ts/metadata'
import InviteUserForm from './InviteUserForm.vue'
import { replaceWithMark, sanitizeTextForHTML } from '../../../ts/text'
import UserTable from '../../generic/UserTable.vue'
import { RoleMetadata } from '../../../ts/roles'
import { TGetAllOrgRolesOutput, getAllOrgRoles } from '../../../ts/api/apiRoles'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'

export default Vue.extend({
    data: () => ({
        filterText: "",
        showHideInvite: false,
        roleList: [] as RoleMetadata[],
        roles: null as Record<number, RoleMetadata> | null,
    }),
    components: {
        InviteUserForm,
        UserTable
    },
    computed: {
        ready() {
            return MetadataStore.state.usersInitialized && !!this.roles
        },

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
    },
    mounted() {
        getAllOrgRoles({
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TGetAllOrgRolesOutput) => {
            this.roles = new Object() as Record<number, RoleMetadata>
            for (let role of resp.data) {
                this.roles[role.Id] = role
            }
            this.roleList = resp.data
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }
})

</script>
