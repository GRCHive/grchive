<template>
    <div v-if="hasRefresh">
        <v-row v-if="isPending" align="center" justify="center">
            <p class="display-1">We are retrieving your schemas. Check back soon!</p>
        </v-row>

        <v-row v-else-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <v-row v-else-if="!refresh.RefreshSuccess" align="center" justify="center">
            <div>
                <p class="display-1">Oops! We were unable to retrieve your schemas.</p>
                <v-divider class="mb-3"></v-divider>
                <pre>{{ refresh.RefreshErrors }}</pre>
            </div>
        </v-row>

        <div v-else>
            <v-select
                v-model="selectedSchema"
                label="Schemas"
                filled
                :items="schemaItems"
                hide-details
                dense
            >
            </v-select>

            <div class="mt-2">
                <database-function-viewer
                    :schema="selectedSchema"
                    v-if="fnMode"
                >
                </database-function-viewer>

                <database-schema-viewer
                    :schema="selectedSchema"
                    v-else
                >
                </database-schema-viewer>
            </div>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import {
    DbRefresh,
    DbSchema
} from '../../ts/sql'
import { allSqlSchemas, TAllSqlSchemasOutput } from '../../ts/api/apiSqlSchemas'
import { Watch } from 'vue-property-decorator'
import { PageParamsStore } from '../../ts/pageParams'
import { contactUsUrl } from '../../ts/url'
import DatabaseSchemaViewer from './DatabaseSchemaViewer.vue'
import DatabaseFunctionViewer from './DatabaseFunctionViewer.vue'

const Props = Vue.extend({
    props: {
        refresh: {
            type: Object,
            default: Object() as () => DbRefresh | null
        },
        fnMode: {
            type: Boolean,
            default: false,
        }
    }
})

@Component({
    components: {
        DatabaseSchemaViewer,
        DatabaseFunctionViewer,
    }
})
export default class DatabaseRefreshViewer extends Props {
    availableSchemas : DbSchema[] | null = null
    selectedSchema : DbSchema | null = null

    get hasRefresh() : boolean {
        return !!this.refresh
    }

    get isPending() : boolean {
        if (!this.refresh) {
            return false
        }

        return this.refresh.RefreshFinishTime == null
    }

    get isLoading() : boolean {
        return (this.availableSchemas == null)
    }

    get schemaItems() : any[] {
        if (!this.availableSchemas) {
            return []
        }

        return this.availableSchemas.map((ele : DbSchema) => ({
            text: ele.SchemaName,
            value: ele
        }))
    }

    @Watch('refresh')
    refreshData() {
        if (!this.refresh) {
            this.availableSchemas = null
            return
        }

        allSqlSchemas({
            refreshId: this.refresh!.Id,
            orgId: PageParamsStore.state.organization!.Id,
        }).then((resp : TAllSqlSchemasOutput) => {
            this.availableSchemas = resp.data
            if (this.availableSchemas!.length > 0) {
                this.selectedSchema = this.availableSchemas[0]
            }
        }).catch((err: any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        })
    }

    mounted() {
        this.refreshData()
    }
}

</script>
