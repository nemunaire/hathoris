<script>
 import { createEventDispatcher } from 'svelte';

 const dispatch = createEventDispatcher();

 export let field = { };
 export let id = null;
 export let value = "";
</script>

{#if field.type.endsWith('int')}
    <input
        type="number"
        id={id}
        class="form-control form-control-sm"
        placeholder={field.placeholder}
        bind:value={value}
        on:change={(e) => dispatch('change', e)}
        on:input={(e) => dispatch('input', e)}
    >
{:else}
    <input
        type="text"
        id={id}
        class="form-control form-control-sm"
        placeholder={field.placeholder}
        bind:value={value}
        on:change={(e) => dispatch('change', e)}
        on:input={(e) => dispatch('input', e)}
    >
{/if}
