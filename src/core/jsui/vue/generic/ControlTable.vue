<script lang="ts">

import Vue, {VNode} from 'vue'
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
                text: 'Control',
                value: 'control',
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
        ]
    }

    get tableItems(): any[] {
        return this.resources.map(this.transformInputResourceToTableItem)
    }

    transformInputResourceToTableItem(inp : any) : any {
        return {
            id: inp.Id,
            control: `${inp.Name} ${inp.Descrption}`,
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

    renderControl(props: any) : VNode { 
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
                    resourceName: "controls"
                },
                on: {
                    input: (items : any[]) => this.$emit('input', items.map((ele : any) => ele.value)),
                    delete: (item : any) => this.$emit('delete', item.value),
                    'click:row': this.goToControl
                },
                scopedSlots: {
                    'item.control': this.renderControl,
                }
            }
        )
    }
}

</script>
