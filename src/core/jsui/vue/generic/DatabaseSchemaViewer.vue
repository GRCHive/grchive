<template>
    <div v-if="hasSchema">
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else>
            <v-text-field outlined
                          v-model="filterString"
                          prepend-inner-icon="mdi-magnify"
                          label="Table Name Filter"
                          hide-details
                          clearable
            ></v-text-field>

            <v-row
                v-for="(tableRow, y) in tableGrid"
                :key="`row${y}`"
            >
                <v-col
                    v-for="(table, x) in tableRow"
                    :key="`col${y}-${x}`"
                    cols="3"
                >
                    <div class="tableGrid">
                        <database-table-viewer
                            :table="table"
                            :columns="allColumns[table.Id]"
                            :height="400"
                        >
                        </database-table-viewer>
                    </div>
                </v-col>
            </v-row>
        </div>
    </div>
</template>

<script lang="ts">

import Vue from 'vue'
import Component from 'vue-class-component'
import { getSqlSchema, TGetSqlSchemaOutput } from '../../ts/api/apiSqlSchemas'
import {
    DbSchema,
    DbTable,
    DbColumn
} from '../../ts/sql'
import { Watch } from 'vue-property-decorator'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import DatabaseTableViewer from './DatabaseTableViewer.vue'

const Props = Vue.extend({
    props: {
        schema: {
            type: Object,
            default: Object() as () => DbSchema | null
        }
    }
})

const gridWidth : number = 4

@Component({
    components: {
        DatabaseTableViewer
    }
})
export default class DatabaseSchemaViewer extends Props {
    allTables : DbTable[] | null = null
    allColumns : Record<number, DbColumn[]> | null = null
    filterString: string | null = ""

    get isLoading() : boolean {
        return this.allTables == null || this.allColumns == null
    }

    get hasSchema() : boolean {
        return !!this.schema
    }

    get filteredTables() : DbTable[] {
        if (!this.allTables) {
            return []
        }

        if (!this.filterString) {
            return this.allTables
        }

        let trimmedFilter = this.filterString.trim()
        if (trimmedFilter == "") {
            return this.allTables
        }

        return this.allTables.filter((ele : DbTable) => {
            return ele.TableName.toLowerCase().includes(trimmedFilter.toLowerCase())
        })
    }

    get tableGrid() : DbTable[][] {
        let fullGrid : DbTable[][] = []
        let fullIdx = 0
        let currentRow : DbTable[] = []
        while (fullIdx < this.filteredTables.length) {
            for (let x = 0; x < gridWidth && fullIdx < this.filteredTables.length; ++x) {
                currentRow.push(this.filteredTables[fullIdx++])
            }

            fullGrid.push(currentRow)
            currentRow = []
        }
        return fullGrid
    }

    @Watch('schema')
    refreshData() {
        if (!this.schema) {
            this.allTables = null
            this.allColumns = null
            return
        }

        getSqlSchema({
            schemaId: this.schema!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            fnMode: false,
        }).then((resp : TGetSqlSchemaOutput) => {
            this.allTables = resp.data.Schema!.Tables
            this.allColumns = resp.data.Schema!.Columns
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

    mounted() { 
        this.refreshData()
    }
}

</script>
