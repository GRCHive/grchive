<template>
    <div v-if="hasSchema">
        <v-row v-if="isLoading" align="center" justify="center">
            <v-progress-circular indeterminate size="64"></v-progress-circular>
        </v-row>

        <div v-else id="schemaTableGrid">
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
                            :height="400"
                        >
                        </database-table-viewer>
                    </div>
                </v-col>
            </v-row>

            <v-row v-if="pullInProgress" align="center" justify="center">
                <v-progress-circular indeterminate size="64"></v-progress-circular>
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
} from '../../ts/sql'
import { Watch } from 'vue-property-decorator'
import { contactUsUrl } from '../../ts/url'
import { PageParamsStore } from '../../ts/pageParams'
import DatabaseTableViewer from './DatabaseTableViewer.vue'

const Props = Vue.extend({
    props: {
        schema: {
            type: Object,
            default: () => null as DbSchema | null
        }
    }
})

const gridWidth : number = 4
const initialPull : number = 56

@Component({
    components: {
        DatabaseTableViewer
    }
})
export default class DatabaseSchemaViewer extends Props {
    allTables : DbTable[] | null = null
    filterString: string = ""
    pullInProgress :boolean = false
    hasMoreData : boolean = true

    get isLoading() : boolean {
        return this.allTables == null
    }

    get hasSchema() : boolean {
        return !!this.schema
    }

    get processedFilterString() : string {
        return this.filterString.trim()
    }

    get filteredTables() : DbTable[] {
        if (!this.allTables) {
            return []
        }

        if (this.processedFilterString == "") {
            return this.allTables
        }

        return this.allTables.filter((ele : DbTable) => {
            return ele.TableName.toLowerCase().includes(this.processedFilterString.toLowerCase())
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
    @Watch('processedFilterString')
    refreshData() {
        this.hasMoreData = true
        if (!this.schema) {
            this.allTables = null
            return
        }

        getSqlSchema({
            schemaId: this.schema!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            fnMode: false,
            start: 0,
            limit: initialPull,
            filter: this.processedFilterString,
        }).then((resp : TGetSqlSchemaOutput) => {
            this.allTables = resp.data.Schema!.Tables
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

    pullMoreTables() {
        if (this.pullInProgress || !this.hasMoreData || !this.allTables) {
            return
        }

        this.pullInProgress = true
        getSqlSchema({
            schemaId: this.schema!.Id,
            orgId: PageParamsStore.state.organization!.Id,
            fnMode: false,
            start: this.allTables!.length,
            limit: initialPull,
            filter: this.processedFilterString,
        }).then((resp : TGetSqlSchemaOutput) => {
            let toAdd = resp.data.Schema!.Tables
            this.hasMoreData = (toAdd.length > 0)
            this.allTables!.push(...toAdd)
        }).catch((err : any) => {
            // @ts-ignore
            this.$root.$refs.snackbar.showSnackBar(
                "Oops! Something went wrong. Try again.",
                true,
                "Contact Us",
                contactUsUrl,
                true);
        }).finally(() => {
            this.pullInProgress = false
        })
    }

    mounted() { 
        this.refreshData()

        window.addEventListener('scroll', this.handleWheel)
    }

    handleWheel(e : Event) {
        //@ts-ignore
        let ele : HTMLElement = document.querySelector(`#schemaTableGrid`)!

        let currentScroll = window.pageYOffset
        let maxScroll = document.documentElement.scrollHeight - document.documentElement.clientHeight

        if (currentScroll >= maxScroll) {
            this.pullMoreTables()
        }
    }
}

</script>
