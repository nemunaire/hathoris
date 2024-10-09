<script>
 import { invalidate } from '$app/navigation';
 import { tick } from 'svelte';

 import Input from '$lib/components/SettingInput.svelte';
 import { loadablesSources } from '$lib/stores/loadable-sources';

 export let data;

 async function submitSettings(avoidStructUpdate) {
     return await data.settings.save(avoidStructUpdate);
 }

 function addCustomSource() {
     data.settings.custom_sources.push({ edit_name: true, kv: { name: "" } });
     data.settings.custom_sources = data.settings.custom_sources;

     tick().then(() => {
         document.getElementById("src_name_" + (data.settings.custom_sources.length - 1)).focus();
     })
 }

 async function deleteCustomSource(isrc) {
     data.settings.custom_sources.splice(isrc);
     data.settings = await submitSettings();
 }
</script>

<div class="container">
    <div class="card my-3">
        <h4 class="card-header">
            <div class="d-flex justify-content-between">
                <div>
                    <i class="bi bi-gear"></i>
                    General Settings
                </div>
            </div>
        </h4>
        <div class="card-body">
            <span class="text-muted">Version:</span>
            {#await data.version then version}
                {version.version}
            {/await}
        </div>
    </div>

    <div class="card my-3">
        <h4 class="card-header">
            <div class="d-flex justify-content-between">
                <div>
                    <i class="bi bi-cassette"></i>
                    Custom Sources
                </div>
                <button
                    class="btn btn-sm btn-primary"
                    on:click={() => addCustomSource()}
                    on:keypress={() => addCustomSource()}
                >
                    <i class="bi bi-plus"></i>
                </button>
            </div>
        </h4>
        <div class="list-group list-group-flush">
            {#each data.settings.custom_sources as source, isrc}
                <div class="list-group-item">
                    <div class="d-flex justify-content-between align-items-start mb-1">
                        <div class="d-flex gap-2 align-items-center">
                            <h5
                                class="mb-0"
                                on:click={() => source.edit_name = true}
                            >
                                {#if source.kv && source.kv.name !== undefined}
                                    {#if !source.edit_name}
                                        {source.kv.name}
                                    {:else}
                                        <input
                                            type="text"
                                            id={"src_name_" + isrc}
                                            class="form-control"
                                            placeholder="Source's name"
                                            bind:value={data.settings.custom_sources[isrc].kv["name"]}
                                            on:input={submitSettings}
                                        >
                                    {/if}
                                {:else}
                                    {source.src} #{isrc}
                                {/if}
                            </h5>
                            {#if source.src}
                                <small class="badge bg-secondary">{source.src}</small>
                            {/if}
                        </div>
                        <button
                            type="button"
                            class="btn btn-sm btn-danger"
                            on:click={() => deleteCustomSource(isrc)}
                            on:keypress={() => deleteCustomSource(isrc)}
                        >
                            <i class="bi bi-trash"></i>
                        </button>
                    </div>
                    {#if source.src && source.kv}
                        {#if $loadablesSources[source.src].fields}
                            {#each $loadablesSources[source.src].fields.filter((e) => e.id !== "name") as field (field.id)}
                                <div class="d-flex gap-3 mb-2">
                                    <label for={"input_" + isrc + "_" + field.id} class="form-label">
                                        {field.label}
                                    </label>
                                    <Input
                                        id={"input_" + isrc + "_" + field.id}
                                        {field}
                                        bind:value={data.settings.custom_sources[isrc].kv[field.id]}
                                        on:input={submitSettings}
                                    />
                                </div>
                            {/each}
                        {:else}
                            {#each Object.keys(source.kv).filter((e) => e !== "name") as k}
                                <div class="d-flex gap-3 mb-2">
                                    <label for={"input_" + isrc + "_" + k} class="form-label">
                                        {k}
                                    </label>
                                    <Input
                                        id={"input_" + isrc + "_" + k}
                                        {field}
                                        bind:value={data.settings.custom_sources[isrc].kv[k]}
                                        on:input={submitSettings}
                                    />
                                </div>
                            {/each}
                        {/if}
                    {:else}
                        Choose a new custom source:
                        <div class="d-flex gap-3">
                            {#if $loadablesSources}
                                {#each Object.keys($loadablesSources) as kls}
                                    <button
                                        class="btn btn-sm btn-primary"
                                        title={$loadablesSources[kls].description}
                                        on:click={() => data.settings.custom_sources[isrc].src = kls}
                                        on:keypress={() => data.settings.custom_sources[isrc].src = kls}
                                    >
                                        {kls}
                                    </button>
                                {/each}
                            {/if}
                        </div>
                    {/if}
                </div>
            {/each}
        </div>
    </div>
</div>
