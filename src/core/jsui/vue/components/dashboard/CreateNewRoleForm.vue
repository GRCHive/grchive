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
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Organization Roles"
            v-model="permissions.OrgRolesAccess"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Process Flows"
            v-model="permissions.ProcessFlowsAccess"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Controls"
            v-model="permissions.ControlsAccess"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Risks"
            v-model="permissions.RisksAccess"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Control Documentation Metadata"
            v-model="permissions.ControlDocMetadataAccess"
        ></access-type-editor>
        <v-divider></v-divider>

        <access-type-editor
            label="Control Documentation"
            v-model="permissions.ControlDocumentationAccess"
        ></access-type-editor>
    </v-form>

    <v-card-actions>
        <v-btn
            color="error"
            @click="cancel"
            v-if="canEdit"
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
import { Permissions, AccessType } from '../../../ts/roles'
import { TNewRoleInput, TNewRoleOutput, newRole} from '../../../ts/api/apiRoles'
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
        }
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
        },
    }),
    computed: {
        canSubmit() : boolean {
            return this.formValid
        }
    },
    methods: {
        cancel() {
            this.$emit('do-cancel')
        },
        save() {
            newRole(<TNewRoleInput>{
                orgId: PageParamsStore.state.organization!.Id,
                name: this.name,
                description: this.description,
                permissions: this.permissions
            }).then((resp : TNewRoleOutput) => {
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
        }
    },
    mounted() {
        this.canEdit = (!this.stagedEdits || !this.editMode)
    }
})

</script>
