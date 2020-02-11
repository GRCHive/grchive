<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Role
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(256)]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :disabled="!canEdit"
                    hide-details
                    class="mb-2">
        </v-textarea> 

        <access-type-editor
            label="Organization Users"
            v-model="permissions.OrgUsersAccess"
            :disabled="!canEdit"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Organization Roles"
            v-model="permissions.OrgRolesAccess"
            :disabled="!canEdit"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Process Flows"
            v-model="permissions.ProcessFlowsAccess"
            :disabled="!canEdit"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Controls"
            v-model="permissions.ControlsAccess"
            :disabled="!canEdit"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Risks"
            v-model="permissions.RisksAccess"
            :disabled="!canEdit"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Control Documentation Metadata"
            v-model="permissions.ControlDocMetadataAccess"
            :disabled="!canEdit"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Control Documentation"
            v-model="permissions.ControlDocumentationAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="General Ledger"
            v-model="permissions.GLAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="Systems"
            v-model="permissions.SystemAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="Database"
            v-model="permissions.DbAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="Database Connection"
            v-model="permissions.DbConnectionAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="Document Requests"
            v-model="permissions.DocRequestAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="Deployments"
            v-model="permissions.DeploymentAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="Servers"
            v-model="permissions.ServerAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="Vendors"
            v-model="permissions.VendorAccess"
            :disabled="!canEdit"
        ></access-type-editor>

        <access-type-editor
            label="SQL Data"
            v-model="permissions.DbSqlAccess"
            :disabled="!canEdit"
        ></access-type-editor>
    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
            v-if="canEdit || dialogMode"
        >
            Cancel
        </v-btn>
        <v-spacer></v-spacer>
        <v-btn
            color="success"
            @click="save"
            :disabled="!canSubmit"
            v-if="canEdit"
        >
            Save
        </v-btn>

        <v-btn
            color="success"
            @click="edit"
            v-if="!canEdit"
        >
            Edit
        </v-btn>
    </v-card-actions>
</v-card>

</template>

<script lang="ts">

import Vue from 'vue'
import * as rules from "../../../ts/formRules"
import { PageParamsStore } from '../../../ts/pageParams'
import { Permissions, AccessType, FullRole } from '../../../ts/roles'
import { TNewRoleInput, TNewRoleOutput, newRole} from '../../../ts/api/apiRoles'
import { TEditRoleInput, TEditRoleOutput, editRole} from '../../../ts/api/apiRoles'
import { contactUsUrl } from '../../../ts/url'
import AccessTypeEditor from '../../generic/AccessTypeEditor.vue'

export default Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        stagedEdits: {
            type: Boolean,
            default: false
        },
        referenceRole: {
            type: Object as () => FullRole,
            default: null
        },
        dialogMode: {
            type: Boolean,
            default: false
        },
    },
    components: {
        AccessTypeEditor,
    },
    data : () => ({
        formValid: false,
        canEdit: false,
        rules,
        name: "",
        description: "",
        permissions: <Permissions>{
            OrgUsersAccess: AccessType.NoAccess,
            OrgRolesAccess: AccessType.NoAccess,
            ProcessFlowsAccess: AccessType.NoAccess,
            ControlsAccess: AccessType.NoAccess,
            ControlDocumentationAccess: AccessType.NoAccess,
            ControlDocMetadataAccess: AccessType.NoAccess,
            RisksAccess: AccessType.NoAccess,
            GLAccess: AccessType.NoAccess,
            SystemAccess: AccessType.NoAccess,
            DbAccess: AccessType.NoAccess,
            DbConnectionAccess: AccessType.NoAccess,
            DocRequestAccess: AccessType.NoAccess,
            DeploymentAccess: AccessType.NoAccess,
            ServerAccess: AccessType.NoAccess,
            VendorAccess: AccessType.NoAccess,
            DbSqlAccess: AccessType.NoAccess,
        },
    }),
    computed: {
        canSubmit() : boolean {
            return this.formValid
        }
    },
    methods: {
        cancel() {
            this.canEdit = false
            this.refreshFromReference()
            this.$emit('do-cancel')
        },
        save() {
            if (this.editMode) {
                this.doEdit()
            } else {
                this.doSave()
            }
        },
        doSave() {
            newRole(<TNewRoleInput>{
                orgId: PageParamsStore.state.organization!.Id,
                name: this.name,
                description: this.description,
                permissions: this.permissions
            }).then((resp : TNewRoleOutput) => {
                this.$emit('do-save', resp.data)
                this.refreshFromReference()
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        doEdit() {
            editRole(<TEditRoleInput>{
                orgId: PageParamsStore.state.organization!.Id,
                roleId: this.referenceRole.RoleMetadata.Id,
                name: this.name,
                description: this.description,
                permissions: this.permissions
            }).then((resp : TEditRoleOutput) => {
                this.canEdit = false
                this.$emit('do-save', resp.data)
            }).catch((err : any) => {
                // @ts-ignore
                this.$root.$refs.snackbar.showSnackBar(
                    "Oops! Something went wrong. Try again.",
                    true,
                    "Contact Us",
                    contactUsUrl,
                    true);
            })
        },
        edit() {
            this.canEdit = true
        },
        refreshFromReference() {
            if (!this.referenceRole) {
                this.name = ""
                this.description = ""
                this.permissions = <Permissions>{
                    OrgUsersAccess: AccessType.NoAccess,
                    OrgRolesAccess: AccessType.NoAccess,
                    ProcessFlowsAccess: AccessType.NoAccess,
                    ControlsAccess: AccessType.NoAccess,
                    ControlDocumentationAccess: AccessType.NoAccess,
                    ControlDocMetadataAccess: AccessType.NoAccess,
                    RisksAccess: AccessType.NoAccess,
                    GLAccess: AccessType.NoAccess,
                    SystemAccess: AccessType.NoAccess,
                    DbAccess: AccessType.NoAccess,
                    DbConnectionAccess: AccessType.NoAccess,
                    DocRequestAccess: AccessType.NoAccess,
                    DeploymentAccess: AccessType.NoAccess,
                    ServerAccess: AccessType.NoAccess,
                    VendorAccess: AccessType.NoAccess,
                    DbSqlAccess: AccessType.NoAccess,
                }
            } else {
                this.name = this.referenceRole.RoleMetadata.Name
                this.description = this.referenceRole.RoleMetadata.Description
                this.permissions = Object.assign({}, this.referenceRole.Permissions)
            }
        }
    },
    mounted() {
        this.canEdit = (!this.stagedEdits || !this.editMode)
        this.refreshFromReference()
    },
    watch: {
        referenceRole() {
            this.refreshFromReference()
        }
    }
})

</script>
