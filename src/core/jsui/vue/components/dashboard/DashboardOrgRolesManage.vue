<template>
    <div class="ma-4">
        <v-list-item class="pa-0">
            <v-list-item-content class="disable-flex mr-4">
                <v-list-item-title class="title">
                    Roles
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
                <v-dialog v-model="showHideCreate" persistent max-width="40%">
                    <template v-slot:activator="{ on }">
                        <v-btn color="primary" v-on="on">
                            New
                        </v-btn>
                    </template>
                    <create-new-role-form
                        @do-cancel="cancelNew"
                        @do-save="saveNew"
                    ></create-new-role-form>
                </v-dialog>
            </v-list-item-action>
        </v-list-item>
        <v-divider></v-divider>

        <role-table
            :resources="roles"
            :search="filterText">
        </role-table>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import { RoleMetadata, FullRole } from '../../../ts/roles'
import { TGetAllOrgRolesInput, TGetAllOrgRolesOutput, getAllOrgRoles } from '../../../ts/api/apiRoles'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgRoleUrl } from '../../../ts/url'
import CreateNewRoleForm from './CreateNewRoleForm.vue'
import RoleTable from '../../generic/RoleTable.vue'

export default Vue.extend({
    data : () => ({
        filterText: "",
        showHideCreate: false,
        roles: [] as RoleMetadata[]
    }),
    components: {
        CreateNewRoleForm,
        RoleTable
    },
    computed: {
        orgGroupName() : string {
            return PageParamsStore.state.organization!.OktaGroupName
        }
    },
    methods: {
        cancelNew() {
            this.showHideCreate = false
        },
        saveNew(role : FullRole) {
            this.showHideCreate = false
            this.roles.push(role.RoleMetadata)
        },
        createOrgRoleUrl
    },
    mounted() {
        getAllOrgRoles(<TGetAllOrgRolesInput>{
            orgId: PageParamsStore.state.organization!.Id
        }).then((resp : TGetAllOrgRolesOutput) => {
            this.roles = resp.data
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
