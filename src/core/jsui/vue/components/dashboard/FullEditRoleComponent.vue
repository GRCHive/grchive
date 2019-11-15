<template>
    <div class="ma-4">
        <v-overlay :value="!ready">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-overlay>

        <div v-if="ready">
            <v-list-item two-line class="pa-0">
                <v-list-item-content>
                    <v-list-item-title class="title">
                        Role: {{ currentRole.RoleMetadata.Name }}
                        <v-btn icon @click="expandDescription = !expandDescription">
                            <v-icon small v-if="!expandDescription" >mdi-chevron-down</v-icon>
                            <v-icon small v-else>mdi-chevron-up</v-icon>
                        </v-btn>
                    </v-list-item-title>

                    <v-list-item-subtitle :class="expandDescription ? `long-text` : `hide-long-text`">
                        {{ currentRole.RoleMetadata.Description }}
                    </v-list-item-subtitle>
                </v-list-item-content>
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
import { TGetSingleRoleInput, TGetSingleRoleOutput, getSingleRole} from '../../../ts/api/apiRoles'
import { lazyGetUserFromId } from '../../../ts/metadataUtils'
import CreateNewRoleForm from './CreateNewRoleForm.vue'

export default Vue.extend({
    data: () => ({
        ready: false,
        expandDescription: false,
        currentRole: Object() as FullRole,
        users: [] as User[]
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
        }
    },
    components: {
        CreateNewRoleForm,
    },
    mounted() {
        this.refreshRoleData()
    }
})

</script>
