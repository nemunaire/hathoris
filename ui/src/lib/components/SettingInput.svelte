<script>
 import { createEventDispatcher } from 'svelte';

 import BasicInput from '$lib/components/SettingInputBasic.svelte';

 const dispatch = createEventDispatcher();

 export let field = { };
 export let id = null;
 export let value = undefined;

 if (field.type && (value === undefined || value === null)) {
     if (field.type.startsWith('[]')) {
         value = [];
     } else if (field.type === 'int') {
         value = 0;
     } else {
         value = "";
     }
 }

 function addElement() {
     if (field.type.endsWith('int')) {
         value.push(0);
     } else {
         value.push("");
     }
     value = value;
 }
</script>

{#if field.type.startsWith("[]")}
    {#each value as val, k}
        {#if k !== 0}<br>{/if}
        <BasicInput
            {field}
            {id}
            bind:value={value[k]}
            on:change={(e) => dispatch('change', e)}
            on:input={(e) => dispatch('input', e)}
        />
    {/each}
    <button
        class="btn btn-sm btn-info"
        on:click={addElement}
        on:keypress={addElement}
    >
        <i class="bi bi-plus" />
    </button>
{:else}
    <BasicInput
        {field}
        {id}
        bind:value={value}
        on:change={(e) => dispatch('change', e)}
        on:input={(e) => dispatch('input', e)}
    />
{/if}
