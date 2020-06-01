<template>

<v-card :flat="hideHeader">
    <template v-if="!hideHeader">
        <v-card-title>
            {{ editMode ? "Edit" : "New" }} SSH Connection (Key)
        </v-card-title>
        <v-divider></v-divider>
    </template>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="username"
                      label="Username"
                      filled
                      :rules="[rules.required]"
                      :readonly="!canEdit"
        >
        </v-text-field>

        <v-file-input
            label="Private Key"
            v-model="file"
            :rules="[rules.required]"
            v-if="canEdit"
        ></v-file-input>

        <v-btn
            block
            @click="downloadKey"
            color="primary"
            v-else
        >
            Download
        </v-btn>
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
            :disabled="!formValid"
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
import Component from 'vue-class-component'
import {
    newServerSSHKeyConnection, editServerSSHKeyConnection,
    TNewServerSSHConnectionOutput,
    getServerSSHKeyConnection, TGetServerKeyConnectionOutput
} from '../../../ts/api/apiServerConnection'
import { contactUsUrl } from '../../../ts/url'
import { PageParamsStore } from '../../../ts/pageParams'
import {
    ServerSSHConnectionGeneric,
} from '../../../ts/infrastructure'
import * as rules from '../../../ts/formRules'
import { saveAs } from 'file-saver'

const Props = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        serverId: Number,
        referenceConnection: {
            type: Object,
            default: () => null as ServerSSHConnectionGeneric | null
        },
        hideHeader: {
            type: Boolean,
            default: false,
        }
    }
})

@Component
export default class CreateServerSshKeyConnection extends Props {
    rules: any = rules

    canEdit: boolean = false
    formValid: boolean = false

    username: string = ""
    file : File | null = null

    onSuccess(r : TNewServerSSHConnectionOutput) {
        this.$emit('do-save', r.data)
        if (!this.editMode) {
            this.clearForm()
        } else {
            this.canEdit = false
        }
    }

    onError(err : any) {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops. Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true);
    }

    doSave() {
        newServerSSHKeyConnection({
            orgId: PageParamsStore.state.organization!.Id,
            serverId: this.serverId,
            username: this.username,
            file: this.file!,
        }).then(this.onSuccess).catch(this.onError)
    }

    doEdit() {
        editServerSSHKeyConnection({
            orgId: PageParamsStore.state.organization!.Id,
            serverId: this.serverId,
            connectionId: this.referenceConnection!.Id,
            username: this.username,
            file: this.file!,
        }).then(this.onSuccess).catch(this.onError)
    }

    save() {
        if (this.editMode) {
            this.doEdit()
        } else {
            this.doSave()
        }
    }

    cancel() {
        this.$emit('do-cancel')
        if (this.editMode) {
            this.clearForm()
            this.canEdit = false
        }
    }

    edit() {
        this.canEdit = true
    }

    clearForm() {
        if (!!this.referenceConnection) {
            this.username = this.referenceConnection!.Username
            this.file = null
        } else {
            this.username = ""
            this.file = null
        }
    }

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }

    downloadKey() {
        getServerSSHKeyConnection({
            orgId: PageParamsStore.state.organization!.Id,
            serverId: this.serverId,
            connectionId: this.referenceConnection!.Id,
        }).then((resp : TGetServerKeyConnectionOutput) => {
            let blob = new Blob([resp.data.PrivateKey], {type: "text/plain;charset=utf-8"})
            saveAs(blob, "key")
        }).catch(this.onError)
    }
}

</script>
