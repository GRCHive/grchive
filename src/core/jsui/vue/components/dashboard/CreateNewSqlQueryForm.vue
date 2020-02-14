<template>

<v-card>
    <v-card-title class="pl-3">
        {{ editMode ? "Edit" : "New" }} SQL Query
    </v-card-title>

    <v-divider></v-divider>

    <v-form class="ma-4" ref="form" v-model="formValid">
        <v-text-field v-model="name"
                      label="Name"
                      filled
                      :rules="[rules.required, rules.createMaxLength(255)]"
                      :readonly="!canEdit"
        ></v-text-field>

        <v-textarea v-model="description"
                    label="Description"
                    filled
                    :readonly="!canEdit"
                    hide-details
        ></v-textarea> 

        <div v-if="allowQueryEdit">
            <p class="subtitle-2 mt-4">SQL Query</p>
            <sql-text-area
                v-model="sqlQuery"
                class="my-4"
                :key="queryKey"
            ></sql-text-area>
        </div>
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
            @click="canEdit = true"
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
import { contactUsUrl } from '../../../ts/url'
import SqlTextArea from '../../generic/SqlTextArea.vue'
import {
    newSqlQuery, TNewSqlQueryOutput,
    updateSqlQuery, TUpdateSqlQueryOutput, TUpdateSqlQueryInput
} from '../../../ts/api/apiSqlQueries'
import { 
    DbSqlQueryMetadata,
    DbSqlQuery,
} from '../../../ts/sql'
import { PageParamsStore } from '../../../ts/pageParams'

const Props = Vue.extend({
    props: {
        dbId: Number,
        editMode: {
            type: Boolean,
            default: false,
        },
        enableQuery: {
            type: Boolean,
            default: true
        },
        referenceMetadata: {
            type: Object as () => DbSqlQueryMetadata | null,
            default: null
        },
        referenceQuery: {
            type: Object as () => DbSqlQuery | null,
            default: null
        }
    }
})

@Component({
    components: {
        SqlTextArea
    }
})
export default class CreateNewSqlQueryForm extends Props {
    canEdit: boolean = false
    formValid: boolean = false
    rules : any = rules

    name : string = ""
    description : string = ""
    sqlQuery : string = ""
    queryKey : number = 0

    get allowQueryEdit() : boolean {
        return !this.editMode || !!this.referenceQuery
    }

    onSuccess(metadata : DbSqlQueryMetadata | null, query : DbSqlQuery | null) {
        this.$emit('do-save', metadata, query)
    }

    onError(err : any) {
        // @ts-ignore
        this.$root.$refs.snackbar.showSnackBar(
            "Oops! Something went wrong. Try again.",
            true,
            "Contact Us",
            contactUsUrl,
            true);
    }

    doSave() {
        newSqlQuery({
            dbId: this.dbId,
            orgId: PageParamsStore.state.organization!.Id,
            name: this.name,
            description: this.description,
            uploadUserId: PageParamsStore.state.user!.Id,
            query: this.sqlQuery,
        }).then((resp : TNewSqlQueryOutput) => {
            this.onSuccess(resp.data.Metadata, resp.data.Query)
            this.clearForm()
        }).catch((err : any) => {
            this.onError(err)
        })
    }

    doEdit() {
        let editData = <TUpdateSqlQueryInput>{
            orgId: PageParamsStore.state.organization!.Id,
            metadataId: this.referenceMetadata!.Id,
        }

        editData.metadata = {
            name: this.name,
            description: this.description,
        }

        if (this.allowQueryEdit) {
            editData.query = {
                query: this.sqlQuery,
                uploadUserId: PageParamsStore.state.user!.Id,
            }
        }

        updateSqlQuery(editData).then((resp : TUpdateSqlQueryOutput) => {
            this.onSuccess(resp.data.Metadata, resp.data.Query)
            this.canEdit = false
        }).catch((err : any) => {
            this.onError(err)
        })
    }

    save() {
        if (!this.formValid) {
            return
        }

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

    mounted() {
        this.canEdit = !this.editMode
        this.clearForm()
    }

    clearForm() {
        if (!!this.referenceMetadata)  {
            this.name = this.referenceMetadata.Name
            this.description = this.referenceMetadata.Description
        } else {
            this.name = ""
            this.description = ""
        }

        if (!!this.referenceQuery) {
            this.sqlQuery = this.referenceQuery.Query
        } else {
            this.sqlQuery = ""
        }

        // Do this so that we can recreate the QuillJS environment so
        // that we can properly clear the text while maintaining the code formatting...:\
        this.queryKey += 1
    }
}

</script>
