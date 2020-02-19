<script lang="ts">

import Vue, {VNode} from 'vue'
import { VCheckbox } from 'vuetify/lib'
import Component from 'vue-class-component'
import BaseResourceTable from './BaseResourceTable.vue'
import ResourceTableProps from './ResourceTableProps'
import { PageParamsStore } from '../../ts/pageParams'
import MetadataStore from '../../ts/metadata'
import { createFrequencyDisplayString } from '../../ts/frequency'
import { createUserString } from '../../ts/users'
import { createControlUrl } from '../../ts/url'

@Component({
    components: {
        BaseResourceTable
    }
})
export default class ControlTable extends ResourceTableProps {
    get tableHeaders() : any[] {
        return [
            {
                text: 'Name',
                value: 'name',
            },
            {
                text: 'Type',
                value: 'type',
            },
            {
                text: 'Owner',
                value: 'owner',
            },
            {
                text: 'Frequency',
                value: 'frequency'
            },
            {
                text: 'Manual',
                value: 'manual'
            },
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            name: `${inp.Name} ${inp.Descrption}`,
            type: MetadataStore.getters.getControlTypeName(inp.ControlTypeId),
            owner: createUserString(MetadataStore.getters.getUser(inp.OwnerId)),
            frequency:createFrequencyDisplayString(inp.FrequencyType, inp.FrequencyInterval, inp.FrequencyOther),
            value: inp
        }
    }

    goToControl(item : any) {
        window.location.assign(createControlUrl(
            PageParamsStore.state.organization!.OktaGroupName,
            item.value.Id
        ))
    }

    renderIsManual(props : any): VNode {
        return this.$createElement(
            VCheckbox,
            {
                class: {
                    'ma-0': true,
                    'pa-0': true,
                },
                props: {
                    'input-value': props.item.value.Manual,
                    disabled: true,
                    'hide-details': true
                }
            }
        )
    }

    renderName(props: any) : VNode { 
        return this.$createElement(
            'div',
            [
                this.$createElement(
                    'p',
                    {
                        class: {
                            'ma-0': true,
                            'pa-0': true,
                            'body-1': true,
                            'font-weight-bold': true
                        },
                        domProps: {
                            innerHTML: props.item.value.Name
                        }
                    },
                ),
                this.$createElement(
                    'p',
                    {
                        class: {
                            'ma-0': true,
                            'pa-0': true,
                            'caption': true,
                            'font-weight-light': true
                        },
                        domProps: {
                            innerHTML: props.item.value.Description
                        }
                    }
                ),
            ]
        )
    }

    render() : VNode {
        return this.$createElement(
            BaseResourceTable,
            {
                props: {
                    ...this.$props,
                    tableHeaders: this.tableHeaders,
                    tableItems: this.tableItems,
                    resourceName: "controls",
                    useGlobalDeletion: true,
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any, g : boolean) => this.$emit('delete', item.value, g),
                    'click:row': this.goToControl
                },
                scopedSlots: {
                    'item.name': this.renderName,
                    'item.manual': this.renderIsManual,
                }
            }
        )
    }
}

</script>
