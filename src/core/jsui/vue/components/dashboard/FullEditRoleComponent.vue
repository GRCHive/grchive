<template>
    <div class="ma-4">
       <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-dialog v-model="showHideDeleteRole" persistent max-width="40%">
                <generic-delete-confirmation-form
                    item-name="roles"
                    :items-to-delete="currentRolesToDelete"
                    v-on:do-cancel="showHideDeleteRole = false"
                    v-on:do-delete="doDelete"
                    :use-global-deletion="false"
                    :force-global-deletion="false">
                </generic-delete-confirmation-form>
            </v-dialog>
 
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Role: {{ currentRole.RoleMetadata.Name }}
                        <span v-if="currentRole.RoleMetadata.IsAdmin">&nbsp;(Admin)</span>
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ currentRole.RoleMetadata.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-spacer></v-spacer>

                <v-list-item-action>
                    <span>
                        <v-tooltip bottom
                                   v-if="!canDelete"
                        >
                            <template v-slot:activator="{ on }">
                                <v-icon small v-on="on">mdi-information</v-icon>
                            </template>
                            <span v-if="currentRole.RoleMetadata.IsDefault">You can not delete the default role. </span>
                            <span v-if="currentRole.RoleMetadata.IsAdmin">You can not delete the admin role. </span>
                            <span v-if="users.length > 0">You can not delete roles with users assigned to it. </span>
                        </v-tooltip>
                        <v-btn color="error"
                              :disabled="!canDelete"
                              @click="requestDelete">
                            Delete
                        </v-btn>
                    </span>
                </v-list-item-action>

            </v-list-item>
            <v-divider></v-divider>

            <v-container fluid>
                <v-row>
                    <v-col cols="9">
                        <create-new-role-form
                            ref="editForm"
                            :staged-edits="true"
                            :edit-mode="true"
                            :reference-role="currentRole"
                            @do-save="editRole">
                        </create-new-role-form>
                    </v-col>

                    <v-col cols="3">
                        <v-card>
                            <v-card-title>
                                Users
                                <v-spacer></v-spacer>

                                <v-dialog v-model="showHideAddUser" persistent max-width="40%">
                                    <template v-slot:activator="{ on }">
                                        <v-btn small color="primary" v-on="on">
                                            Add User
                                        </v-btn>
                                    </template>

                                    <v-card>
                                        <v-card-title>
                                            Select Users to Add to the
                                            <span class="font-weight-bold">&nbsp;{{ currentRole.RoleMetadata.Name }}&nbsp;</span>
                                            Role
                                        </v-card-title>
                                        <v-divider></v-divider>
                                        <user-table
                                            :users="selectableOrgUsers"
                                            v-model="usersToAdd"
                                            selectable
                                            multi>
                                        </user-table>

                                        <v-card-actions>
                                            <v-btn color="error" @click="showHideAddUser = false">
                                                Cancel
                                            </v-btn>
                                            <v-spacer></v-spacer>
                                            <v-btn color="success"
                                                   :disabled="usersToAdd.length == 0"
                                                   @click="addUsersToRole">
                                                Add
                                            </v-btn>
                                        </v-card-actions>
                                    </v-card>

                                </v-dialog>
                            </v-card-title>
                            <v-divider></v-divider>

                            <user-table :users="users">
                            </user-table>
                        </v-card>
                    </v-col>
                </v-row>
            </v-container>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import { PageParamsStore } from '../../../ts/pageParams'
import { FullRole } from '../../../ts/roles'
import { createOrgAllRolesUrl, contactUsUrl } from '../../../ts/url'
import { TGetSingleRoleInput, TGetSingleRoleOutput, getSingleRole} from '../../../ts/api/apiRoles'
import { TDeleteRoleInput, TDeleteRoleOutput, deleteRole} from '../../../ts/api/apiRoles'
import { TAddUsersToRoleInput, TAddUsersToRoleOutput, addUsersToRole} from '../../../ts/api/apiRoles'
import { lazyGetUserFromId } from '../../../ts/metadataUtils'
import CreateNewRoleForm from './CreateNewRoleForm.vue'
import GenericDeleteConfirmationForm from './GenericDeleteConfirmationForm.vue'
import UserTable from '../../generic/UserTable.vue'
import MetadataStore from '../../../ts/metadata'
import { listSubtract } from '../../../ts/sets'

export default Vue.extend({
    data: () => ({
        ready: false,
        expandDescription: false,
        currentRole: Object() as FullRole,
        users: [] as User[],
        showHideDeleteRole: false,
        showHideAddUser: false,
        usersToAdd: [] as User[]
    }),
    methods: {
        refreshRoleData() {
            let data = window.location.pathname.split('/')
            let roleId = Number(data[data.length - 1])

            getSingleRole(<TGetSingleRoleInput>{
                orgId: PageParamsStore.state.organization!.Id,
                roleId: roleId,
            }).then((resp : TGetSingleRoleOutput) => {
                this.currentRole = resp.data.role

                let promises = []
                for (let userId of resp.data.userIds) {
                    promises.push(lazyGetUserFromId(userId))
                }

                Promise.all(promises).then((vals) => {
                    this.users = vals
                })

                this.ready = true
            }).catch((err : any) => {
                window.location.replace('/404')
            })
        },
        editRole(role : FullRole) {
            Vue.set(this.currentRole, 'RoleMetadata', role.RoleMetadata)
            Vue.set(this.currentRole, 'Permissions', role.Permissions)
        },
        requestDelete() {
            this.showHideDeleteRole = true
        },
        doDelete() {
            deleteRole(<TDeleteRoleInput>{
                roleId: this.currentRole.RoleMetadata.Id,
                orgId: PageParamsStore.state.organization!.Id,
            }).then((resp : TDeleteRoleOutput) => {
                this.showHideDeleteRole = false
                window.location.replace(createOrgAllRolesUrl(PageParamsStore.state.organization!.OktaGroupName))
            }).catch((err : any) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        addUsersToRole() {
            addUsersToRole(<TAddUsersToRoleInput>{
                roleId: this.currentRole.RoleMetadata.Id,
                orgId: PageParamsStore.state.organization!.Id,
                userIds: this.usersToAdd.map((ele : User) => ele.Id)
            }).then((resp : TAddUsersToRoleOutput) => {
                this.users.push(...this.usersToAdd)
                this.usersToAdd = []
                this.showHideAddUser = false
            }).catch((err : any) => {
                //@ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Ensure that there is at least one user assigned to the admin role after this change and try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        }
    },
    computed: {
        canDelete(): boolean {
            return !this.currentRole.RoleMetadata.IsDefault && !this.currentRole.RoleMetadata.IsAdmin && this.users.length == 0
        },
        currentRolesToDelete() : string[] {
            return [this.currentRole.RoleMetadata.Name]
        },
        selectableOrgUsers() : User[] {
            return listSubtract<User>(MetadataStore.state.availableUsers, this.users)
        }
    },
    components: {
        CreateNewRoleForm,
        GenericDeleteConfirmationForm,
        UserTable
    },
    mounted() {
        this.refreshRoleData()
    }
})

</script>
