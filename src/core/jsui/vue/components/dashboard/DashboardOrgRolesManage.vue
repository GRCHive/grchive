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

        <v-list-item>
            <v-list-item-content class="font-weight-bold">
                <v-list-item-title>
                    Role
                </v-list-item-title>
            </v-list-item-content>

            <v-list-item-content class="font-weight-bold">
                <v-list-item-title>
                    Default
                </v-list-item-title>
            </v-list-item-content>

            <v-list-item-content class="font-weight-bold">
                <v-list-item-title>
                    Admin
                </v-list-item-title>
            </v-list-item-content>
        </v-list-item>

        <v-list two-line>
            <v-list-item
                v-for="(item, index) in filteredRoles"
                :key="index"
                :href="createOrgRoleUrl(orgGroupName, item.Id)"
            >
                <v-list-item-content>
                    <v-list-item-title v-html="`${item.Name}`">
                    </v-list-item-title>

                    <v-list-item-subtitle v-html="`${item.Description}`">
                    </v-list-item-subtitle>
                </v-list-item-content>

                <v-list-item-content>
                    <v-checkbox class="ma-0" v-model="item.IsDefault" :hide-details="true" disabled></v-checkbox>
                </v-list-item-content>

                <v-list-item-content>
                    <v-checkbox class="ma-0" v-model="item.IsAdmin" :hide-details="true" disabled></v-checkbox>
                </v-list-item-content>
            </v-list-item>
        </v-list>

    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import { RoleMetadata, FullRole } from '../../../ts/roles'
import { TGetAllOrgRolesInput, TGetAllOrgRolesOutput, getAllOrgRoles } from '../../../ts/api/apiRoles'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl, createOrgRoleUrl } from '../../../ts/url'
import CreateNewRoleForm from './CreateNewRoleForm.vue'

export default Vue.extend({
    data : () => ({
        filterText: "",
        showHideCreate: false,
        roles: [] as RoleMetadata[]
    }),
    components: {
        CreateNewRoleForm
    },
    computed: {
        filteredRoles() : RoleMetadata[] {
            return this.roles
        },
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
            console.log(role)
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
