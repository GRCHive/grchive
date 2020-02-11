<template>

<v-card>
    <v-card-title>
        {{ editMode ? "Edit" : "New" }} Database
    </v-card-title>
    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-autocomplete
            v-model="typeId"
            filled
            :items="databaseTypeSelect"
            label="Type"
            hide-no-data
            :rules="[rules.required]"
            :disabled="!canEdit"
        >
        </v-autocomplete>

        <v-text-field v-if="isOtherType"
                      v-model="otherType"
                      label="Please Specify"
                      filled
                      :rules="[rules.required]"
                      :disabled="!canEdit">
        </v-text-field>

        <v-text-field v-model="version"
                      label="Version"
                      filled
                      :rules="[rules.required]"
                      :disabled="!canEdit">
        </v-text-field>
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
import * as rules from '../../../ts/formRules'
import { TNewDatabaseOutputs, newDatabase} from '../../../ts/api/apiDatabases'
import { TEditDatabaseOutputs, editDatabase} from '../../../ts/api/apiDatabases'
import { PageParamsStore } from '../../../ts/pageParams'
import { contactUsUrl } from '../../../ts/url'
import { otherTypeId, DatabaseType, Database } from '../../../ts/databases'
import MetadataStore from '../../../ts/metadata'

const VueComponent = Vue.extend({
    props: {
        editMode: {
            type: Boolean,
            default: false
        },
        dialogMode: {
            type: Boolean,
            default: false
        },
        referenceDb: {
            type: Object as () => Database | null,
            default: null
        }
    }
})

@Component
export default class CreateNewDatabaseForm extends VueComponent {
    canEdit : boolean = false
    formValid: boolean = false
    rules: any = rules

    name : string = ""
    typeId: number | null = null
    otherType: string = ""
    version: string = ""

    get isOtherType() : boolean {
        return this.typeId == otherTypeId
    }

    get databaseTypeSelect() : any[] {
        return MetadataStore.state.availableDbTypes.map((ele : DatabaseType) => ({
            text: ele.Name,
            value: ele.Id
        }))
    }

    get isDatabaseTypeLoading() : boolean {
        return !MetadataStore.state.dbTypesInitialized
    }

    doSave() {
        newDatabase({
            name: this.name,
            orgId: PageParamsStore.state.organization!.Id,
            typeId: this.typeId!,
            otherType: this.otherType,
            version: this.version,
        }).then((resp : TNewDatabaseOutputs) => {
            this.$emit('do-save', resp.data)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
    }

    doEdit() {
        editDatabase({
            dbId: this.referenceDb!.Id,
            name: this.name,
            orgId: PageParamsStore.state.organization!.Id,
            typeId: this.typeId!,
            otherType: this.otherType,
            version: this.version,
        }).then((resp : TNewDatabaseOutputs) => {
            this.$emit('do-save', resp.data)
            this.canEdit = false
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops. Something went wrong. Try again.",
                false,
                "",
                contactUsUrl,
                true);
        })
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
        this.clearForm()

        if (this.editMode) {
            this.canEdit = false
        }

    }

    edit() {
        this.canEdit = true
    }

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }

    clearForm() {
        if (!!this.referenceDb) {
            this.name = this.referenceDb.Name
            this.typeId = this.referenceDb.TypeId
            this.otherType = this.referenceDb.OtherType
            this.version = this.referenceDb.Version
        } else {
            this.name = ""
            this.typeId = null
            this.otherType = ""
            this.version = ""
        }
    }
}

</script>
